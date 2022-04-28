/*
Copyright 2021 Flant JSC

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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook/metrics"
	"github.com/flant/addon-operator/sdk"
	"github.com/golang/protobuf/proto"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/deckhouse/deckhouse/go_lib/dependency"
)

// this hook checks helm releases (v2 and v3) and find deprecated apis
// hook returns only metrics:
//   `resource_versions_compatibility` for apis:
//      1 - is deprecated
//      2 - in unsupported
// Also hook returns count on deployed releases `helm_releases_count`
// Hook checks only releases with status: deployed

// **Attention**
// Releases are checked via kubeclient not by snapshots to avoid huge memory consumption
// on some installations snapshots can take gigabytes of memory

const unsupportedVersionsYAML = `
"1.22":
  "admissionregistration.k8s.io/v1beta1": ["ValidatingWebhookConfiguration", "MutatingWebhookConfiguration"]
  "apiextensions.k8s.io/v1beta1": ["CustomResourceDefinition"]
  "apiregistration.k8s.io/v1beta1": ["APIService"]
  "authentication.k8s.io/v1beta1": ["TokenReview"]
  "authorization.k8s.io/v1beta1": ["SubjectAccessReview", "LocalSubjectAccessReview", "SelfSubjectAccessReview"]
  "certificates.k8s.io/v1beta1": ["CertificateSigningRequest"]
  "coordination.k8s.io/v1beta1": ["Lease"]
  "networking.k8s.io/v1beta1": ["Ingress"]
  "extensions/v1beta1": ["Ingress"]

"1.25":
  "batch/v1beta1": ["CronJob"]
  "discovery.k8s.io/v1beta1": ["EndpointSlice"]
  "events.k8s.io/v1beta1": ["Event"]
  "autoscaling/v2beta1": ["HorizontalPodAutoscaler"]
  "policy/v1beta1": ["PodDisruptionBudget", "PodSecurityPolicy"]
  "node.k8s.io/v1beta1": ["RuntimeClass"]
`

// delta for k8s versions which are checked for deprecated apis
// with delta == 2 for k8s 1.21 will also check apis for 1.22 and 1.23
const delta = 2

var (
	storage unsupportedVersionsStore
)

func init() {
	err := yaml.Unmarshal([]byte(unsupportedVersionsYAML), &storage)
	if err != nil {
		log.Fatal(err)
	}
}

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: "/modules/helm/helm_releases",
	Schedule: []go_hook.ScheduleConfig{
		{
			Name:    "helm_releases",
			Crontab: "*/1 * * * *",
		},
	},
}, dependency.WithExternalDependencies(handleHelmReleases))

func handleHelmReleases(input *go_hook.HookInput, dc dependency.Container) error {
	input.MetricsCollector.Expire("helm_deprecated_apiversions")

	k8sCurrentVersionRaw, ok := input.Values.GetOk("global.discovery.kubernetesVersion")
	if !ok {
		input.LogEntry.Warn("kubernetes version not found")
		return nil
	}
	k8sCurrentVersion := semver.MustParse(k8sCurrentVersionRaw.String())

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		helm3Releases, err := getHelm3Releases(ctx, dc)
		if err != nil {
			input.LogEntry.Error(err)
			return
		}

		input.MetricsCollector.Set("helm_releases_count", float64(len(helm3Releases)), map[string]string{"helm_version": "3"})
		err = processHelmReleases(k8sCurrentVersion, input, helm3Releases)
		if err != nil {
			input.LogEntry.Error(err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		helm2Releases, err := getHelm2Releases(ctx, dc)
		if err != nil {
			input.LogEntry.Error(err)
			return
		}

		input.MetricsCollector.Set("helm_releases_count", float64(len(helm2Releases)), map[string]string{"helm_version": "2"})
		err = processHelmReleases(k8sCurrentVersion, input, helm2Releases)
		if err != nil {
			input.LogEntry.Error(err)
			return
		}
	}()

	wg.Wait()

	return nil
}

func getHelm3Releases(ctx context.Context, dc dependency.Container) ([]*release, error) {
	var result []*release
	client, err := dc.GetK8sClient()
	if err != nil {
		return nil, err
	}
	secretsList, err := client.CoreV1().Secrets("").List(ctx, metav1.ListOptions{LabelSelector: "owner=helm,status=deployed"})
	if err != nil {
		return nil, err
	}

	for _, secret := range secretsList.Items {
		releaseData := secret.Data["release"]
		if len(releaseData) == 0 {
			continue
		}

		release, err := decodeRelease(string(releaseData))
		if err != nil {
			return nil, err
		}

		result = append(result, release)
	}

	return result, nil
}

func getHelm2Releases(ctx context.Context, dc dependency.Container) ([]*release, error) {
	var result []*release
	client, err := dc.GetK8sClient()
	if err != nil {
		return nil, err
	}
	cmList, err := client.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{LabelSelector: "OWNER=TILLER,STATUS=DEPLOYED"})
	if err != nil {
		return nil, err
	}

	for _, cm := range cmList.Items {
		releaseData := cm.Data["release"]
		if len(releaseData) == 0 {
			continue
		}

		release, err := helm2DecodeRelease(releaseData)
		if err != nil {
			return nil, err
		}

		result = append(result, release)
	}

	return result, nil
}

func processHelmReleases(k8sCurrentVersion *semver.Version, input *go_hook.HookInput, releases []*release) error {
	for _, rel := range releases {

		d := yaml.NewDecoder(strings.NewReader(rel.Manifest))

		for {
			resource := new(manifest)
			err := d.Decode(&resource)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				return err
			}

			if resource == nil {
				continue
			}

			incompatibility, k8sCompatibilityVersion := storage.CalculateCompatibility(k8sCurrentVersion, resource.APIVersion, resource.Kind)
			if incompatibility > 0 {
				input.MetricsCollector.Set("resource_versions_compatibility", float64(incompatibility), map[string]string{
					"helm_release_name":      rel.Name,
					"helm_release_namespace": rel.Namespace,
					"k8s_version":            k8sCompatibilityVersion,
					"resource_name":          resource.Metadata.Name,
					"resource_namespace":     resource.Metadata.Namespace,
					"kind":                   resource.Kind,
					"api_version":            resource.APIVersion,
				}, metrics.WithGroup("helm_deprecated_apiversions"))
			}

		}
	}

	return nil
}

// protobuf for handling helm2 releases - https://github.com/helm/helm/blob/47f0b88409e71fd9ca272abc7cd762a56a1c613e/pkg/proto/hapi/release/release.pb.go#L24
type release struct {
	Name      string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name,proto3"`
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,8,opt,name=namespace,proto3"`
	Manifest  string `json:"manifest,omitempty" protobuf:"bytes,5,opt,name=manifest,proto3" `
}

type manifest struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   struct {
		Name      string `json:"name" yaml:"name"`
		Namespace string `json:"namespace" yaml:"namespace"`
	} `json:"metadata" yaml:"metadata"`
}

// k8s version: APIVersion: Kind
type unsupportedVersionsStore map[string]unsupportedAPIVersions

func (uvs unsupportedVersionsStore) getByK8sVersion(version *semver.Version) (unsupportedAPIVersions, bool) {
	majorMinor := fmt.Sprintf("%d.%d", version.Major(), version.Minor())
	apis, ok := uvs[majorMinor]
	return apis, ok
}

// CalculateCompatibility check compatibility. Returns
//   0 - if resource is compatible
//   1 - if resource in deprecated and will be removed in the future
//   2 - if resource is unsupported for current k8s version
//  and k8s version in which deprecation would be
func (uvs unsupportedVersionsStore) CalculateCompatibility(currentVersion *semver.Version, resourceAPIVersion, resourceKind string) (uint, string) {
	// check unsupported api for current k8s version
	currentK8SAPIsStorage, exists := uvs.getByK8sVersion(currentVersion)
	if exists {
		isUnsupported := currentK8SAPIsStorage.isUnsupportedByAPIAndKind(resourceAPIVersion, resourceKind)
		if isUnsupported {
			return 2, fmt.Sprintf("%d.%d", currentVersion.Major(), currentVersion.Minor())
		}
	}

	// if api is supported - check deprecation in the next 2 minor k8s versions
	for i := 0; i < delta; i++ {
		newMinor := currentVersion.Minor() + uint64(i)
		nextVersion := semver.MustParse(fmt.Sprintf("%d.%d.0", currentVersion.Major(), newMinor))
		storage, exists := uvs.getByK8sVersion(nextVersion)
		if exists {
			isDeprecated := storage.isUnsupportedByAPIAndKind(resourceAPIVersion, resourceKind)
			if isDeprecated {
				return 1, fmt.Sprintf("%d.%d", nextVersion.Major(), nextVersion.Minor())
			}
		}
	}

	return 0, ""
}

// APIVersion: [Kind]
type unsupportedAPIVersions map[string][]string

func (ua unsupportedAPIVersions) isUnsupportedByAPIAndKind(api, ikind string) bool {
	kinds, ok := ua[api]
	if !ok {
		return false
	}
	for _, kind := range kinds {
		if kind == ikind {
			return true
		}
	}

	return false
}

// helm3 decoding

var magicGzip = []byte{0x1f, 0x8b, 0x08}

// Import this from helm3 lib - https://github.com/helm/helm/blob/49819b4ef782e80b0c7f78c30bd76b51ebb56dc8/pkg/storage/driver/util.go#L56
// decodeRelease decodes the bytes of data into a release
// type. Data must contain a base64 encoded gzipped string of a
// valid release, otherwise an error is returned.
func decodeRelease(data string) (*release, error) {
	// base64 decode string
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		defer r.Close()
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var rls release
	// unmarshal release object bytes
	if err := json.Unmarshal(b, &rls); err != nil {
		return nil, err
	}
	return &rls, nil
}

// https://github.com/helm/helm/blob/47f0b88409e71fd9ca272abc7cd762a56a1c613e/pkg/storage/driver/util.go#L57
func helm2DecodeRelease(data string) (*release, error) {
	// base64 decode string
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var rls release
	// unmarshal protobuf bytes
	if err := proto.Unmarshal(b, &rls); err != nil {
		return nil, err
	}
	return &rls, nil
}

// protobuf methods for helm2
func (m *release) Reset()         { *m = release{} }
func (m *release) String() string { return proto.CompactTextString(m) }
func (*release) ProtoMessage()    {}
