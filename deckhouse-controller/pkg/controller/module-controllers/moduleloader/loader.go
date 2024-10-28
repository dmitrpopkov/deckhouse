// Copyright 2024 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package moduleloader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
	d8env "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/env"
	"github.com/deckhouse/deckhouse/go_lib/dependency/extenders"
	"github.com/flant/addon-operator/pkg/module_manager/loader"
	"github.com/flant/addon-operator/pkg/module_manager/models/modules"
	"github.com/flant/addon-operator/pkg/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	moduleOrderIdx = 2
	moduleNameIdx  = 3
)

var (
	// some ephemeral modules, which we even don't want to load
	excludeModules = map[string]struct{}{
		"000-common":           {},
		"007-registrypackages": {},
	}

	// validModuleNameRe defines a valid module name. It may have a number prefix: it is an order of the module.
	validModuleNameRe = regexp.MustCompile(`^(([0-9]+)-)?(.+)$`)
)

var _ loader.ModuleLoader = &Loader{}

type Loader struct {
	client      client.Client
	modulesDirs []string
	modules     map[string]*Module
}

func New(client client.Client, modulesDir string) *Loader {
	return &Loader{
		client:      client,
		modulesDirs: utils.SplitToPaths(modulesDir),
		modules:     make(map[string]*Module),
	}
}

// LoadModules implements the module loader interface from addon-operator, used for registering modules in addon-operator
func (l *Loader) LoadModules() ([]*modules.BasicModule, error) {
	result := make([]*modules.BasicModule, 0, len(l.modules))

	for _, module := range l.modules {
		result = append(result, module.GetBasicModule())
	}

	return result, nil
}

// LoadModule implements the module loader interface from addon-operator, it reads single directory and returns BasicModule
// modulePath is in the following format: /deckhouse-controller/downloaded/<module_name>/<module_version>
func (l *Loader) LoadModule(_, modulePath string) (*modules.BasicModule, error) {
	if _, err := readDir(modulePath); err != nil {
		return nil, err
	}

	// run moduleDefinitionByDir("<module_name>", "/deckhouse-controller/downloaded/<module_name>/<module_version>")
	def, err := l.moduleDefinitionByDir(filepath.Base(filepath.Dir(modulePath)), modulePath)
	if err != nil {
		return nil, err
	}

	module, err := l.processModuleDefinition(def)
	if err != nil {
		return nil, err
	}
	l.modules[def.Name] = module

	return module.GetBasicModule(), nil
}

func (l *Loader) processModuleDefinition(def *Definition) (*Module, error) {
	if err := validateModuleName(def.Name); err != nil {
		return nil, err
	}

	// load values for module
	valuesModuleName := utils.ModuleNameToValuesKey(def.Name)
	// 1. from static values.yaml inside the module
	moduleStaticValues, err := utils.LoadValuesFileFromDir(def.Path)
	if err != nil {
		return nil, err
	}

	if moduleStaticValues.HasKey(valuesModuleName) {
		moduleStaticValues = moduleStaticValues.GetKeySection(valuesModuleName)
	}

	// 2. from openapi defaults
	configBytes, vb, err := utils.ReadOpenAPIFiles(filepath.Join(def.Path, "openapi"))
	if err != nil {
		return nil, err
	}

	module, err := newModule(def, moduleStaticValues, configBytes, vb)
	if err != nil {
		return nil, fmt.Errorf("new deckhouse module: %w", err)
	}

	// load conversions
	if _, err = os.Stat(filepath.Join(def.Path, "openapi", "conversions")); err == nil {
		log.Debugf("conversions for the '%q' module found", valuesModuleName)
		if err = conversion.Store().Add(def.Name, filepath.Join(def.Path, "openapi", "conversions")); err != nil {
			log.Debugf("loading conversions for the '%q' module failed", valuesModuleName)
			return nil, err
		}
	} else {
		if !os.IsNotExist(err) {
			//log.Debugf("loading conversions for the '%q' module failed", valuesModuleName)
			return nil, err
		}
		log.Debugf("conversions for the '%q' module not found", valuesModuleName)
	}

	// load constrains
	if err = extenders.AddConstraints(def.Name, def.Requirements); err != nil {
		return nil, err
	}

	return module, nil
}

func validateModuleName(name string) error {
	// Check if name is consistent for conversions between kebab-case and camelCase.
	restoredName := utils.ModuleNameFromValuesKey(utils.ModuleNameToValuesKey(name))

	if name != restoredName {
		return fmt.Errorf("'%s' name should be in kebab-case and be restorable from camelCase: consider renaming to '%s'", name, restoredName)
	}

	return nil
}

func (l *Loader) GetModuleByName(name string) (*Module, error) {
	module, ok := l.modules[name]
	if !ok {
		return nil, errors.New("module not found")
	}

	return module, nil
}

// LoadModulesFromFS parses and ensures modules from FS
func (l *Loader) LoadModulesFromFS(ctx context.Context) error {
	for _, dir := range l.modulesDirs {
		definitions, err := l.parseModulesDir(dir)
		if err != nil {
			return err
		}
		for _, def := range definitions {
			module := new(Module)
			if module, err = l.processModuleDefinition(def); err != nil {
				return fmt.Errorf("faield to process the '%s' module: %w", def.Name, err)
			}

			if _, ok := l.modules[def.Name]; ok {
				log.Warnf("the '%q' module is already exists. Skipping module from %q", def.Name, def.Path)
				continue
			}

			if !strings.HasPrefix(def.Path, d8env.GetDownloadedModulesDir()) {
				if err = l.ensureEmbeddedModule(ctx, def); err != nil {
					return fmt.Errorf("failed to ensure the embedded '%s' module: %w", def.Name, err)
				}
			}

			l.modules[def.Name] = module
		}
	}

	return nil
}

func (l *Loader) ensureEmbeddedModule(ctx context.Context, def *Definition) error {
	module := new(v1alpha1.Module)
	if err := l.client.Get(ctx, client.ObjectKey{Name: def.Name}, module); err != nil {
		if apierrors.IsNotFound(err) {
			if module.Properties.Weight != def.Weight {
				module.Properties.Weight = def.Weight
				// TODO(ipaqsa): add retry on conflict
				return l.client.Update(ctx, module)
			}
			return nil
		}
		return err
	}
	module = &v1alpha1.Module{
		TypeMeta: metav1.TypeMeta{
			Kind:       v1alpha1.ModuleGVK.Kind,
			APIVersion: v1alpha1.ModuleGVK.GroupVersion().String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: def.Name,
		},
		Properties: v1alpha1.ModuleProperties{
			Weight:           def.Weight,
			Source:           v1alpha1.ModuleSourceEmbedded,
			AvailableSources: []string{},
		},
		Status: v1alpha1.ModuleStatus{
			Phase: string(v1alpha1.ModulePhaseReady),
			Conditions: []v1alpha1.ModuleCondition{
				{
					Type:               v1alpha1.ModuleConditionEnabled,
					Status:             corev1.ConditionFalse,
					LastTransitionTime: metav1.Now(),
					LastProbeTime:      metav1.Now(),
				},
			},
		},
	}
	return l.client.Create(ctx, module)
}

// parseModulesDir returns modules definitions from the target dir
func (l *Loader) parseModulesDir(modulesDir string) ([]*Definition, error) {
	entries, err := readDir(modulesDir)
	if err != nil {
		return nil, err
	}

	definitions := make([]*Definition, 0)
	for _, entry := range entries {
		name, absPath, err := resolveDirEntry(modulesDir, entry)
		if err != nil {
			return nil, err
		}
		// skip non-directories.
		if name == "" {
			continue
		}

		if _, ok := excludeModules[name]; ok {
			continue
		}

		definition, err := l.moduleDefinitionByDir(name, absPath)
		if err != nil {
			return nil, err
		}

		definitions = append(definitions, definition)
	}

	return definitions, nil
}

// readDir checks if dir exists and returns entries
func readDir(dir string) ([]os.DirEntry, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path '%s' does not exist", dir)
		}
		return nil, fmt.Errorf("failed to list modules directory '%s': %s", dir, err)
	}
	return dirEntries, nil
}

func resolveDirEntry(dirPath string, entry os.DirEntry) (string, string, error) {
	name := entry.Name()
	absPath := filepath.Join(dirPath, name)

	if entry.IsDir() {
		return name, absPath, nil
	}
	// Check if entry is a symlink to a directory.
	targetPath, err := resolveSymlinkToDir(dirPath, entry)
	if err != nil {
		// TODO: probably we can use os.IsNotExist here
		if e, ok := err.(*fs.PathError); ok {
			if e.Err.Error() == "no such file or directory" {
				//log.Warnf("Symlink target %q does not exist. Ignoring module", dirPath)
				return "", "", nil
			}
		}

		return "", "", fmt.Errorf("failed to resolve '%s' as a possible symlink: %v", absPath, err)
	}

	if targetPath != "" {
		return name, targetPath, nil
	}

	if name != utils.ValuesFileName {
		log.Warnf("ignore '%s' while searching for modules", absPath)
	}

	return "", "", nil
}

func resolveSymlinkToDir(dirPath string, entry os.DirEntry) (string, error) {
	info, err := entry.Info()
	if err != nil {
		return "", err
	}

	targetDirPath, isTargetDir, err := utils.SymlinkInfo(filepath.Join(dirPath, info.Name()), info)
	if err != nil {
		return "", err
	}

	if isTargetDir {
		return targetDirPath, nil
	}

	return "", nil
}

// moduleDefinitionByDir parses module's definition from the target dir
func (l *Loader) moduleDefinitionByDir(moduleName, moduleDir string) (*Definition, error) {
	definition, err := l.moduleDefinitionByFile(moduleDir)
	if err != nil {
		return nil, err
	}

	if definition == nil {
		//log.Debugf("module.yaml for module %q does not exist", moduleName)
		definition, err = l.moduleDefinitionByDirName(moduleName, moduleDir)
		if err != nil {
			return nil, err
		}
	}

	return definition, nil
}

// moduleDefinitionByFile returns Module instance parsed from the module.yaml file
func (l *Loader) moduleDefinitionByFile(absPath string) (*Definition, error) {
	path := filepath.Join(absPath, DefinitionFile)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	def := new(Definition)
	if err = yaml.NewDecoder(f).Decode(def); err != nil {
		return nil, err
	}

	if def.Name == "" || def.Weight == 0 {
		return nil, nil
	}
	def.Path = absPath

	return def, nil
}

// moduleDefinitionByDirName returns Module instance filled with name, order and its absolute path.
func (l *Loader) moduleDefinitionByDirName(dirName string, absPath string) (*Definition, error) {
	matchRes := validModuleNameRe.FindStringSubmatch(dirName)
	if matchRes == nil {
		return nil, fmt.Errorf("'%s' is invalid name for module: should match regex '%s'", dirName, validModuleNameRe.String())
	}

	return &Definition{
		Name:   matchRes[moduleNameIdx],
		Path:   absPath,
		Weight: parseUintOrDefault(matchRes[moduleOrderIdx], 100),
	}, nil
}

func parseUintOrDefault(num string, defaultValue uint32) uint32 {
	val, err := strconv.ParseUint(num, 10, 31)
	if err != nil {
		return defaultValue
	}
	return uint32(val)
}
