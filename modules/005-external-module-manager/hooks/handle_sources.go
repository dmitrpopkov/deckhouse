/*
Copyright 2023 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hooks

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/flant/addon-operator/pkg/module_manager"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	labels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/pointer"

	deckhouse_config "github.com/deckhouse/deckhouse/go_lib/deckhouse-config"
	"github.com/deckhouse/deckhouse/go_lib/dependency"
	"github.com/deckhouse/deckhouse/go_lib/dependency/cr"
	"github.com/deckhouse/deckhouse/modules/005-external-module-manager/hooks/internal/apis/v1alpha1"
)

const (
	defaultModuleWeight = 900
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: "/modules/external-module-source/sources",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:                "sources",
			ApiVersion:          "deckhouse.io/v1alpha1",
			Kind:                "ModuleSource",
			ExecuteHookOnEvents: pointer.Bool(true),
			FilterFunc:          filterSource,
		},
		{
			Name:                "policies",
			ApiVersion:          "deckhouse.io/v1alpha1",
			Kind:                "ModuleUpdatePolicy",
			ExecuteHookOnEvents: pointer.Bool(true),
			FilterFunc:          filterPolicy,
		},
	},
	Schedule: []go_hook.ScheduleConfig{
		{
			Name:    "check_module_releases",
			Crontab: "*/3 * * * *",
		},
	},
	Settings: &go_hook.HookConfigSettings{
		EnableSchedulesOnStartup: true,
	},
}, dependency.WithExternalDependencies(handleSource))

func filterPolicy(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var mus v1alpha1.ModuleUpdatePolicy

	err := sdk.FromUnstructured(obj, &mus)
	if err != nil {
		return nil, err
	}

	newmus := v1alpha1.ModuleUpdatePolicy{
		TypeMeta: mus.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name: mus.Name,
			UID:  mus.UID,
		},
		Spec: mus.Spec,
	}

	return newmus, nil
}

func filterSource(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var ms v1alpha1.ModuleSource

	err := sdk.FromUnstructured(obj, &ms)
	if err != nil {
		return nil, err
	}

	if ms.Spec.Registry.Scheme == "" {
		// fallback to default https protocol
		ms.Spec.Registry.Scheme = "HTTPS"
	}

	// remove unused fields
	newms := v1alpha1.ModuleSource{
		TypeMeta: ms.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name: ms.Name,
		},
		Spec: ms.Spec,
		Status: v1alpha1.ModuleSourceStatus{
			ModuleErrors: ms.Status.ModuleErrors,
		},
	}

	return newms, nil
}

func handleSource(input *go_hook.HookInput, dc dependency.Container) error {
	externalModulesDir := os.Getenv("EXTERNAL_MODULES_DIR")
	checksumFilePath := path.Join(externalModulesDir, "checksum.json")
	ts := time.Now().UTC()

	snapPolicies := input.Snapshots["policies"]
	if len(snapPolicies) == 0 {
		input.LogEntry.Info("not a single module update policy was found - skip external module releases")
		return nil
	}

	snapSources := input.Snapshots["sources"]
	if len(snapSources) == 0 {
		return nil
	}

	sourcesChecksum, err := getSourceChecksums(checksumFilePath)
	if err != nil {
		return err
	}

	policies := make(map[string]*modulePolicy)

	for _, pol := range snapPolicies {
		policy := pol.(v1alpha1.ModuleUpdatePolicy)
		policies[policy.Name] = &modulePolicy{
			spec:            policy.Spec,
			affectedModules: make([]string, 0),
			errors:          make([]string, 0),
		}
	}

	for _, source := range snapSources {
		ex := source.(v1alpha1.ModuleSource)
		sc := v1alpha1.ModuleSourceStatus{
			SyncTime: ts,
		}

		opts := make([]cr.Option, 0)
		if ex.Spec.Registry.DockerCFG != "" {
			opts = append(opts, cr.WithAuth(ex.Spec.Registry.DockerCFG))
		} else {
			opts = append(opts, cr.WithDisabledAuth())
		}

		if ex.Spec.Registry.CA != "" {
			opts = append(opts, cr.WithCA(ex.Spec.Registry.CA))
		}

		if ex.Spec.Registry.Scheme == "HTTP" {
			opts = append(opts, cr.WithInsecureSchema(true))
		}

		regCli, err := dc.GetRegistryClient(ex.Spec.Registry.Repo, opts...)
		if err != nil {
			sc.Msg = err.Error()
			updateSourceStatus(input, ex.Name, sc)
			continue
		}

		tags, err := regCli.ListTags()
		if err != nil {
			sc.Msg = err.Error()
			updateSourceStatus(input, ex.Name, sc)
			continue
		}

		sort.Strings(tags)

		sc.Msg = ""
		sc.AvailableModules = tags
		sc.ModulesCount = len(tags)
		moduleErrors := make([]v1alpha1.ModuleError, 0)

		mChecksum := make(moduleChecksum)

		if data, ok := sourcesChecksum[ex.Name]; ok {
			mChecksum = data
		}

		for _, moduleName := range tags {
			if moduleName == "modules" {
				input.LogEntry.Warn("'modules' name for module is forbidden. Skip module.")
				continue
			}

			foundPolicy, policyConflict, policyName, releaseChannel, err := getReleasePolicy(ex.Name, moduleName, policies)
			if err != nil {
				return fmt.Errorf("couldn't check module %s release: %v", moduleName, err)
			}

			if !foundPolicy {
				input.LogEntry.Infof("module %s release from module source %s was skipped as no matching policies were found", moduleName, ex.Name)
				continue
			}

			if policyConflict {
				input.LogEntry.Warnf("module %s release from module source %s was skipped as it is matched by multiple module policies", moduleName, ex.Name)
				continue
			}

			moduleVersion, err := fetchModuleVersion(input.LogEntry, dc, ex, moduleName, releaseChannel, mChecksum, opts)
			if err != nil {
				moduleErrors = append(moduleErrors, v1alpha1.ModuleError{
					Name:  moduleName,
					Error: err.Error(),
				})
				continue
			}

			if moduleVersion == "" {
				if len(ex.Status.ModuleErrors) > 0 {
					// inherit errors to keep failed modules status
					for _, mer := range ex.Status.ModuleErrors {
						if mer.Name == moduleName {
							moduleErrors = append(moduleErrors, mer)
							break
						}
					}
				}
				// checksum has not been changed
				continue
			}

			moduleVersionPath := path.Join(externalModulesDir, moduleName, moduleVersion)

			err = fetchAndCopyModuleByVersion(dc, moduleVersionPath, ex, moduleName, moduleVersion, opts)
			if err != nil {
				moduleErrors = append(moduleErrors, v1alpha1.ModuleError{
					Name:  moduleName,
					Error: err.Error(),
				})
				continue
			}

			weight := fetchModuleWeight(moduleVersionPath)

			err = validateModule(moduleName, moduleVersionPath, weight)
			if err != nil {
				moduleErrors = append(moduleErrors, v1alpha1.ModuleError{
					Name:  moduleName,
					Error: err.Error(),
				})
				continue
			}

			createRelease(input, ex.Name, moduleName, moduleVersion, policyName, weight)
		}

		sc.ModuleErrors = moduleErrors
		if len(sc.ModuleErrors) > 0 {
			sc.Msg = "Some errors occurred. Inspect status for details"
		} else {
			sourcesChecksum[ex.Name] = mChecksum
		}
		updateSourceStatus(input, ex.Name, sc)
	}

	updatePolicyStatuses(input, policies)

	// save checksums
	err = saveSourceChecksums(checksumFilePath, sourcesChecksum)
	if err != nil {
		return err
	}

	return nil
}

func validateModule(moduleName, absPath string, weight int) error {
	module, err := module_manager.NewModuleWithNameValidation(moduleName, absPath, weight)
	if err != nil {
		return err
	}

	err = deckhouse_config.Service().ValidateModule(module)
	if err != nil {
		return err
	}

	if weight < 900 || weight > 999 {
		return fmt.Errorf("external module weight must be between 900 and 999")
	}

	return nil
}

func getSourceChecksums(checksumFilePath string) (sourceChecksum, error) {
	var sourcesChecksum sourceChecksum

	if _, err := os.Stat(checksumFilePath); err == nil {
		checksumFile, err := os.Open(checksumFilePath)
		if err != nil {
			return nil, err
		}
		defer checksumFile.Close()

		err = json.NewDecoder(checksumFile).Decode(&sourcesChecksum)
		if err != nil {
			if err == io.EOF {
				return make(sourceChecksum), nil
			}
			return nil, err
		}

		return sourcesChecksum, nil
	}

	return make(sourceChecksum), nil
}

func saveSourceChecksums(checksumFilePath string, checksums sourceChecksum) error {
	data, _ := json.Marshal(checksums)

	return os.WriteFile(checksumFilePath, data, 0666)
}

func fetchModuleVersion(logger *logrus.Entry, dc dependency.Container, moduleSource v1alpha1.ModuleSource, moduleName, moduleReleaseChannel string, modulesChecksum map[string]string, registryOptions []cr.Option) ( /* moduleVersion */ string, error) {
	regCli, err := dc.GetRegistryClient(path.Join(moduleSource.Spec.Registry.Repo, moduleName, "release"), registryOptions...)
	if err != nil {
		return "", fmt.Errorf("fetch release image error: %v", err)
	}

	img, err := regCli.Image(strcase.ToKebab(moduleReleaseChannel))
	if err != nil {
		return "", fmt.Errorf("fetch image error: %v", err)
	}

	digest, err := img.Digest()
	if err != nil {
		return "", fmt.Errorf("fetch digest error: %v", err)
	}

	if prev, ok := modulesChecksum[moduleName]; ok {
		if prev == digest.String() {
			logger.Infof("Module %s checksum has not been changed. Ignoring.", moduleName)
			return "", nil
		}
	}

	modulesChecksum[moduleName] = digest.String()

	moduleMetadata, err := fetchModuleReleaseMetadata(img)
	if err != nil {
		return "", fmt.Errorf("fetch release metadata error: %v", err)
	}

	return "v" + moduleMetadata.Version.String(), nil
}

func fetchModuleWeight(moduleVersionPath string) int {
	moduleDefFile := path.Join(moduleVersionPath, module_manager.ModuleDefinitionFileName)

	if _, err := os.Stat(moduleDefFile); err != nil {
		return defaultModuleWeight
	}

	var def module_manager.ModuleDefinition

	f, err := os.Open(moduleDefFile)
	if err != nil {
		return defaultModuleWeight
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&def)
	if err != nil {
		return defaultModuleWeight
	}

	return def.Weight
}

func fetchAndCopyModuleByVersion(dc dependency.Container, moduleVersionPath string, moduleSource v1alpha1.ModuleSource, moduleName, moduleVersion string, registryOptions []cr.Option) error {
	regCli, err := dc.GetRegistryClient(path.Join(moduleSource.Spec.Registry.Repo, moduleName), registryOptions...)
	if err != nil {
		return fmt.Errorf("fetch module error: %v", err)
	}

	img, err := regCli.Image(moduleVersion)
	if err != nil {
		return fmt.Errorf("fetch module version error: %v", err)
	}

	_ = os.RemoveAll(moduleVersionPath)

	err = copyModuleToFS(moduleVersionPath, img)
	if err != nil {
		return fmt.Errorf("copy module error: %v", err)
	}

	// inject registry to values
	err = injectRegistryToModuleValues(moduleVersionPath, moduleSource)
	if err != nil {
		return fmt.Errorf("inject registry error: %v", err)
	}

	return nil
}

func injectRegistryToModuleValues(moduleVersionPath string, moduleSource v1alpha1.ModuleSource) error {
	valuesFile := path.Join(moduleVersionPath, "openapi", "values.yaml")

	valuesData, err := os.ReadFile(valuesFile)
	if err != nil {
		return err
	}

	valuesData, err = mutateOpenapiSchema(valuesData, moduleSource)
	if err != nil {
		return err
	}

	return os.WriteFile(valuesFile, valuesData, 0666)
}

func mutateOpenapiSchema(sourceValuesData []byte, moduleSource v1alpha1.ModuleSource) ([]byte, error) {
	reg := new(registrySchemaForValues)
	reg.SetBase(moduleSource.Spec.Registry.Repo)
	reg.SetDockercfg(moduleSource.Spec.Registry.DockerCFG)

	var yamlData injectedValues

	err := yaml.Unmarshal(sourceValuesData, &yamlData)
	if err != nil {
		return nil, err
	}

	yamlData.Properties.Registry = reg

	buf := bytes.NewBuffer(nil)

	yamlEncoder := yaml.NewEncoder(buf)
	yamlEncoder.SetIndent(2)
	err = yamlEncoder.Encode(yamlData)

	return buf.Bytes(), err
}

func copyModuleToFS(rootPath string, img v1.Image) error {
	rc := mutate.Extract(img)
	defer rc.Close()

	err := copyLayersToFS(rootPath, rc)
	if err != nil {
		return fmt.Errorf("copy tar to fs: %w", err)
	}

	return nil
}

func copyLayersToFS(rootPath string, rc io.ReadCloser) error {
	if err := os.MkdirAll(rootPath, 0700); err != nil {
		return fmt.Errorf("mkdir root path: %w", err)
	}

	tr := tar.NewReader(rc)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of archive
			return nil
		}
		if err != nil {
			return fmt.Errorf("tar reader next: %w", err)
		}

		if strings.Contains(hdr.Name, "..") {
			// CWE-22 check, prevents path traversal
			return fmt.Errorf("path traversal detected in the module archive: malicious path %v", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path.Join(rootPath, hdr.Name), 0700); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(path.Join(rootPath, hdr.Name))
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("copy: %w", err)
			}
			outFile.Close()

			err = os.Chmod(outFile.Name(), os.FileMode(hdr.Mode)&0700) // remove only 'user' permission bit, E.x.: 644 => 600, 755 => 700
			if err != nil {
				return fmt.Errorf("chmod: %w", err)
			}
		case tar.TypeSymlink:
			link := path.Join(rootPath, hdr.Name)
			if err := os.Symlink(hdr.Linkname, link); err != nil {
				return fmt.Errorf("create symlink: %w", err)
			}
		case tar.TypeLink:
			err := os.Link(path.Join(rootPath, hdr.Linkname), path.Join(rootPath, hdr.Name))
			if err != nil {
				return fmt.Errorf("create hardlink: %w", err)
			}

		default:
			return errors.New("unknown tar type")
		}
	}
}

func untarMetadata(rc io.ReadCloser, rw io.Writer) error {
	tr := tar.NewReader(rc)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of archive
			return nil
		}
		if err != nil {
			return err
		}
		if strings.HasPrefix(hdr.Name, ".werf") {
			continue
		}

		switch hdr.Name {
		case "version.json":
			_, err = io.Copy(rw, tr)
			if err != nil {
				return err
			}
			return nil

		default:
			continue
		}
	}
}

func fetchModuleReleaseMetadata(img v1.Image) (moduleReleaseMetadata, error) {
	buf := bytes.NewBuffer(nil)
	var meta moduleReleaseMetadata

	layers, err := img.Layers()
	if err != nil {
		return meta, err
	}

	for _, layer := range layers {
		size, err := layer.Size()
		if err != nil {
			// dcr.logger.Warnf("couldn't calculate layer size")
			return meta, err
		}
		if size == 0 {
			// skip some empty werf layers
			continue
		}
		rc, err := layer.Uncompressed()
		if err != nil {
			return meta, err
		}

		err = untarMetadata(rc, buf)
		if err != nil {
			return meta, err
		}

		rc.Close()
	}

	err = json.Unmarshal(buf.Bytes(), &meta)

	return meta, err
}

func getReleaseLabels(source, module string) map[string]string {
	return map[string]string{"module": module, "source": source}
}

func getReleasePolicy(sourceName, moduleName string, policies map[string]*modulePolicy) (bool, bool, string, string, error) {
	var labelsSet labels.Set = getReleaseLabels(sourceName, moduleName)
	var matchedPolicy string
	var found, conflict bool
	for name, policy := range policies {
		if policy.spec.ModuleReleaseSelector.LabelSelector != nil {
			selector, err := metav1.LabelSelectorAsSelector(policy.spec.ModuleReleaseSelector.LabelSelector)
			if err != nil {
				return false, false, "", "", err
			}

			if selector.Matches(labelsSet) {
				policies[name].affectedModules = append(policies[name].affectedModules, moduleName)
				if found {
					policies[name].errors = append(policies[name].errors, fmt.Sprintf("module %s is also affected by %s policy", moduleName, matchedPolicy))
					policies[matchedPolicy].errors = append(policies[matchedPolicy].errors, fmt.Sprintf("module %s is also affected by %s policy", moduleName, name))
					conflict = true
				} else {
					found = true
					matchedPolicy = name
				}
			}
		}
	}

	if !found {
		return found, conflict, "", "", nil
	}

	return found, conflict, matchedPolicy, policies[matchedPolicy].spec.ReleaseChannel, nil
}

func createRelease(input *go_hook.HookInput, sourceName, moduleName, moduleVersion, policyName string, moduleWeight int) {
	moduleReleaseLabels := getReleaseLabels(sourceName, moduleName)
	moduleReleaseLabels["module-update-policy"] = policyName

	rl := &v1alpha1.ModuleRelease{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ModuleRelease",
			APIVersion: "deckhouse.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-%s", moduleName, moduleVersion),
			Annotations: make(map[string]string),
			Labels:      moduleReleaseLabels,
		},
		Spec: v1alpha1.ModuleReleaseSpec{
			ModuleName: moduleName,
			Version:    semver.MustParse(moduleVersion),
			Weight:     moduleWeight,
		},
	}

	input.PatchCollector.Create(rl, object_patch.UpdateIfExists())
}

func updatePolicyStatuses(input *go_hook.HookInput, policies map[string]*modulePolicy) {
	for name, policy := range policies {
		status := map[string]v1alpha1.ModuleUpdatePolicyStatus{
			"status": {
				MatchedModules: strings.Join(policy.affectedModules, ", "),
				Errors:         strings.Join(policy.errors, ", "),
			},
		}
		input.PatchCollector.MergePatch(status, "deckhouse.io/v1alpha1", "ModuleUpdatePolicy", "", name, object_patch.WithSubresource("/status"))
	}
}

func updateSourceStatus(input *go_hook.HookInput, name string, sc v1alpha1.ModuleSourceStatus) {
	st := map[string]v1alpha1.ModuleSourceStatus{
		"status": sc,
	}

	input.PatchCollector.MergePatch(st, "deckhouse.io/v1alpha1", "ModuleSource", "", name, object_patch.WithSubresource("/status"))
}

type moduleReleaseMetadata struct {
	Version *semver.Version `json:"version"`
}

type moduleChecksum map[string]string

type sourceChecksum map[string]moduleChecksum

type modulePolicy struct {
	errors          []string
	affectedModules []string
	spec            v1alpha1.ModuleUpdatePolicySpec
}

// part of openapi schema for injecting registry values
type registrySchemaForValues struct {
	Type       string   `yaml:"type"`
	Default    struct{} `yaml:"default"`
	Properties struct {
		Base struct {
			Type    string `yaml:"type"`
			Default string `yaml:"default"`
		} `yaml:"base"`
		Dockercfg struct {
			Type    string `yaml:"type"`
			Default string `yaml:"default,omitempty"`
		} `yaml:"dockercfg"`
	} `yaml:"properties"`
}

func (rsv *registrySchemaForValues) fillTypes() {
	rsv.Properties.Base.Type = "string"
	rsv.Properties.Dockercfg.Type = "string"
	rsv.Type = "object"
}

func (rsv *registrySchemaForValues) SetBase(registryBase string) {
	rsv.fillTypes()

	rsv.Properties.Base.Default = registryBase
}

func (rsv *registrySchemaForValues) SetDockercfg(dockercfg string) {
	if len(dockercfg) == 0 {
		return
	}

	rsv.Properties.Dockercfg.Default = dockercfg
}

type injectedValues struct {
	Type       string                 `json:"type" yaml:"type"`
	Xextend    map[string]interface{} `json:"x-extend" yaml:"x-extend"`
	Properties struct {
		YYY      map[string]interface{}   `json:",inline" yaml:",inline"`
		Registry *registrySchemaForValues `json:"registry" yaml:"registry"`
	} `json:"properties" yaml:"properties"`
	XXX map[string]interface{} `json:",inline" yaml:",inline"`
}
