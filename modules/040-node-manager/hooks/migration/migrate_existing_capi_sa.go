/*
Copyright 2024 Flant JSC

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
	"fmt"

	"github.com/deckhouse/deckhouse/modules/110-istio/hooks/lib"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/ptr"
)

const (
	clusterAPINamespace = "d8-cloud-instance-manager"
)

type ServiceAccountInfo struct {
	IsLabeledAndAnnotated bool
	Name                  string
}

const (
	labelHelmManagedBy         = "app.kubernetes.io/managed-by"
	annotationReleaseName      = "meta.helm.sh/release-name"
	annotationReleaseNamespace = "meta.helm.sh/release-namespace"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue:        "/modules/node-manager",
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 1},
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:                         "capi_sa",
			ExecuteHookOnSynchronization: ptr.To(false),
			ExecuteHookOnEvents:          ptr.To(false),
			ApiVersion:                   "v1",
			Kind:                         "ServiceAccount",
			NamespaceSelector:            lib.NsSelector(),
			FilterFunc:                   applyServiceAccountFilter,
		},
	},
}, migrateServiceAccounts)

func applyServiceAccountFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	serviceAccount := &v1.ServiceAccount{}
	err := sdk.FromUnstructured(obj, serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("cannot convert ServiceAccount to struct: %v", err)
	}

	_, isLabeledWithHelmManaged := serviceAccount.Labels[labelHelmManagedBy]
	_, isAnnotatedWithReleseName := serviceAccount.Annotations[annotationReleaseName]
	_, isAnnotatedWithReleseNamespace := serviceAccount.Annotations[annotationReleaseNamespace]

	return ServiceAccountInfo{
		IsLabeledAndAnnotated: isLabeledWithHelmManaged && isAnnotatedWithReleseName && isAnnotatedWithReleseNamespace,
		Name:                  serviceAccount.GetName(),
	}, nil
}

func migrateServiceAccounts(input *go_hook.HookInput) error {
	patch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": map[string]string{
				labelHelmManagedBy: "Helm",
			},
			"annotations": map[string]string{
				annotationReleaseName:      "node-manager",
				annotationReleaseNamespace: "d8-system",
			},
		},
	}

	snap := input.Snapshots
	for _, serviceAccountSnap := range snap["capi_sa"] {
		serviceAccount := serviceAccountSnap.(ServiceAccountInfo)
		if !serviceAccount.IsLabeledAndAnnotated {
			input.PatchCollector.MergePatch(patch, "v1", "ServiceAccount", "d8-cloud-instance-manager", serviceAccount.Name, object_patch.IgnoreMissingObject())
		}
	}
	return nil
}
