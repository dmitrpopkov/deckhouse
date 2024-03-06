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

package requirements

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/deckhouse/deckhouse/go_lib/dependency/requirements"
)

const (
	requirementsKey                  = "istioVer"
	k8sKey                           = "k8s"
	minVersionValuesKey              = "istio:minimalVersion"
	isK8sVersionAutomaticKey         = "istio:isK8sVersionAutomatic"
	compatibilityOperatorToK8sVerKey = "istio:compatibilityOperatorToK8sVer"
)

func init() {
	checkRequirementFunc := func(requirementValue string, getter requirements.ValueGetter) (bool, error) {
		comingK8sVersion, err := semver.NewVersion(requirementValue)
		if err != nil {
			return false, err
		}
		currentVersionRaw, exists := getter.Get(minVersionValuesKey)
		if !exists {
			return true, nil
		}
		currentVersionStr := currentVersionRaw.(string)
		currentVersion, err := semver.NewVersion(currentVersionStr)
		if err != nil {
			return false, err
		}

		if currentVersion.LessThan(comingK8sVersion) {
			return false, fmt.Errorf("installed Istio version '%s' is lower than required", currentVersionStr)
		}

		return true, nil
	}

	checkMaximalK8sVersionForOperator := func(requirementValue string, getter requirements.ValueGetter) (bool, error) {
		comingK8sVersion := requirementValue

		currentIstioVersionRaw, exists := getter.Get(minVersionValuesKey)
		if !exists {
			return true, nil
		}
		currentMinIstioVersionStr := currentIstioVersionRaw.(string)

		isAtomaticK8sVerRaw, exists := getter.Get(isK8sVersionAutomaticKey)
		if !exists {
			return true, nil
		}
		isAtomaticK8sVer := isAtomaticK8sVerRaw.(bool)
		// Only if k8s version is set to Automatic in cluster
		if !isAtomaticK8sVer {
			return true, nil
		}

		compatibilityMapRaw, exists := getter.Get(compatibilityOperatorToK8sVerKey)
		if !exists {
			return true, nil
		}
		compatibilityMap, ok := compatibilityMapRaw.(map[string][]string)
		if !ok {
			return true, nil
		}

		if k8sVersions, ok := compatibilityMap[currentMinIstioVersionStr]; ok {
			for _, k8sVersion := range k8sVersions {
				// If k8s version in compatibility list of operator version
				if comingK8sVersion == k8sVersion {
					return true, nil
				}
			}
			return false, fmt.Errorf("after update kubernetes version '%s' will be incompatible with operator version '%s'", comingK8sVersion, currentMinIstioVersionStr)
		}

		return true, nil
	}

	requirements.RegisterCheck(requirementsKey, checkRequirementFunc)
	requirements.RegisterCheck(k8sKey, checkMaximalK8sVersionForOperator)
}
