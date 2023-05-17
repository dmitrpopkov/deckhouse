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

// this hook checks if there any clusterAuthorizationRules with limitNamespaces option set

package hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook/metrics"
	"github.com/flant/addon-operator/sdk"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "cluster_authorization_rules",
			ApiVersion: "deckhouse.io/v1",
			Kind:       "ClusterAuthorizationRule",
			FilterFunc: applyClusterAuthorizationRuleFilter,
		},
	},
}, handleClusterAuthorizationRulesWithLimitNamespaces)

type ObjectNameKind struct {
	Name string
	Kind string
}

func applyClusterAuthorizationRuleFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	namespaces, found, err := unstructured.NestedStringSlice(obj.Object, "spec", "limitNamespaces")
	if err != nil {
		return nil, err
	}
	if found {
		if len(namespaces) > 0 {
			return &ObjectNameKind{
				Name: obj.GetName(),
				Kind: obj.GetKind(),
			}, nil
		}
	}
	return nil, nil
}

func handleClusterAuthorizationRulesWithLimitNamespaces(input *go_hook.HookInput) error {
	input.MetricsCollector.Expire("d8_deprecated_car_spec_limitnamespaces")
	for _, obj := range input.Snapshots["cluster_authorization_rules"] {
		if obj == nil {
			continue
		}
		objMeta := obj.(*ObjectNameKind)
		input.MetricsCollector.Set("d8_deprecated_car_spec_limitnamespaces", 1, map[string]string{"kind": objMeta.Kind, "name": objMeta.Name}, metrics.WithGroup("d8_deprecated_car_spec_limitnamespaces"))
	}
	return nil
}
