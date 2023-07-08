/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package internal

import (
	"github.com/deckhouse/deckhouse/ee/modules/160-multitenancy-manager/hooks/apis/deckhouse.io/v1alpha1"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	ProjectTypeHookKubeConfig = go_hook.KubernetesConfig{
		Name:       ProjectTypesQueue,
		ApiVersion: APIVersion,
		Kind:       ProjectTypeKind,
		FilterFunc: filterProjectTypes,
	}
)

type ProjectTypeSnapshot struct {
	Name string
	Spec v1alpha1.ProjectTypeSpec
}

func filterProjectTypes(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	pt := &v1alpha1.ProjectType{}
	if err := sdk.FromUnstructured(obj, pt); err != nil {
		return nil, err
	}

	return ProjectTypeSnapshot{
		Name: pt.Name,
		Spec: pt.Spec,
	}, nil
}
