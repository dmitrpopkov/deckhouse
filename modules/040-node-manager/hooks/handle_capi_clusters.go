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
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	namespace  = "d8-cloud-instance-manager"
	apiVersion = "cluster.x-k8s.io/v1beta1"
)

// filterDynamicProbeNodeGroups returns the name of a nodegroup to consider or emptystring if it should be skipped
func filterClusters(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	return obj.GetName(), nil
}

// This hook discovers nodegroup names for dynamic probes in upmeter
var _ = sdk.RegisterFunc(
	&go_hook.HookConfig{
		Queue:        "/modules/node-manager",
		OnBeforeHelm: &go_hook.OrderedConfig{Order: 20},
		Kubernetes: []go_hook.KubernetesConfig{
			{
				Name:       "clusters",
				ApiVersion: apiVersion,
				Kind:       "Cluster",
				NamespaceSelector: &types.NamespaceSelector{
					NameSelector: &types.NameSelector{
						MatchNames: []string{namespace},
					},
				},
				FilterFunc: filterClusters,
			},
		},
	},
	updateClusterStatus,
)

// collectDynamicProbeConfig sets names of objects to internal values
func updateClusterStatus(input *go_hook.HookInput) error {
	statusPatch := map[string]interface{}{
		"status": map[string]interface{}{
			"infrastructureReady": true,
		},
	}

	snap := input.Snapshots["clusters"]
	for _, sn := range snap {
		input.PatchCollector.MergePatch(statusPatch, apiVersion, "Cluster", namespace, sn.(string), object_patch.IgnoreMissingObject())
	}
	return nil
}
