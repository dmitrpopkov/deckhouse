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
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/deckhouse/deckhouse/go_lib/dependency/requirements"
	"github.com/deckhouse/deckhouse/modules/110-istio/hooks/lib"
	"github.com/deckhouse/deckhouse/modules/110-istio/hooks/lib/crd"
	"github.com/deckhouse/deckhouse/modules/110-istio/hooks/lib/istio_versions"
)

const (
	minVersionValuesKey      = "istio:minimalVersion"
	operatorK8sMinVersionKey = "istio:minimalVersionK8sMinimal"
	operatorK8sMaxVersionKey = "istio:minimalVersionK8sMaximal"
)

type IstioOperatorCrdInfo struct {
	Name     string
	Revision string
}

func applyIstioOperatorFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var iop crd.IstioOperator

	err := sdk.FromUnstructured(obj, &iop)
	if err != nil {
		return nil, err
	}

	return IstioOperatorCrdInfo{
		Name:     iop.GetName(),
		Revision: iop.Spec.Revision,
	}, nil
}

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: lib.Queue("discovery"),
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:              "istiooperators",
			ApiVersion:        "install.istio.io/v1alpha1",
			Kind:              "IstioOperator",
			FilterFunc:        applyIstioOperatorFilter,
			NamespaceSelector: lib.NsSelector(),
		},
	},
}, operatorRevisionsToInstallDiscovery)

func operatorRevisionsToInstallDiscovery(input *go_hook.HookInput) error {
	var operatorVersionsToInstall = make([]string, 0)
	var unsupportedRevisions = make([]string, 0)

	versionMap := istio_versions.VersionMapJSONToVersionMap(input.Values.Get("istio.internal.versionMap").String())

	// Get array of compatibility k8s versions for every operator version
	k8sCompatibleVersions := make(map[string][]string)
	_ = json.Unmarshal([]byte(input.Values.Get("istio.internal.istioToK8sCompatibilityMap").String()), &k8sCompatibleVersions)

	var versionsToInstallResult = input.Values.Get("istio.internal.versionsToInstall").Array()
	for _, versionResult := range versionsToInstallResult {
		operatorVersionsToInstall = append(operatorVersionsToInstall, versionResult.String())
	}

	for _, iop := range input.Snapshots["istiooperators"] {
		iopInfo := iop.(IstioOperatorCrdInfo)
		iopVer := versionMap.GetVersionByRevision(iopInfo.Revision)
		if !versionMap.IsRevisionSupported(iopInfo.Revision) {
			unsupportedRevisions = append(unsupportedRevisions, iopInfo.Revision)
			continue
		}
		if !lib.Contains(operatorVersionsToInstall, iopVer) {
			operatorVersionsToInstall = append(operatorVersionsToInstall, iopVer)
		}
	}

	if len(unsupportedRevisions) > 0 {
		sort.Strings(unsupportedRevisions)
		return fmt.Errorf("unsupported revisions: [%s]", strings.Join(unsupportedRevisions, ","))
	}

	sort.Strings(operatorVersionsToInstall)
	input.Values.Set("istio.internal.operatorVersionsToInstall", operatorVersionsToInstall)

	// Getting minVersion
	var minVersion *semver.Version
	for _, version := range operatorVersionsToInstall {
		versionSemver, err := semver.NewVersion(version)
		if err != nil {
			return err
		}
		if minVersion == nil || versionSemver.LessThan(minVersion) {
			minVersion = versionSemver
		}
	}
	if minVersion == nil {
		requirements.RemoveValue(minVersionValuesKey)
		requirements.RemoveValue(operatorK8sMinVersionKey)
		requirements.RemoveValue(operatorK8sMaxVersionKey)
	} else {
		requirements.SaveValue(minVersionValuesKey, minVersion.String())
		minVer, maxVer := getMinAndMaxVersionsFromCompMap(k8sCompatibleVersions, minVersion.String())
		requirements.SaveValue(operatorK8sMinVersionKey, minVer)
		requirements.SaveValue(operatorK8sMaxVersionKey, maxVer)
	}

	return nil
}

// Use the fact compatability versions array is ["1.18", "1.19", "1.20", "1.21"] and key is like "1.19.3"
func getMinAndMaxVersionsFromCompMap(cmpMap map[string][]string, key string) (string, string) {
	v := strings.Split(key, ".")
	if len(v) < 2 {
		return "", ""
	}
	keySuffix := v[0] + "." + v[1]

	if values, ok := cmpMap[keySuffix]; ok {
		if len(values) > 1 {
			return values[0], values[len(values)-1]
		}
	}
	return "", ""
}
