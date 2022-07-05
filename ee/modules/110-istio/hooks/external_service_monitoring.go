/*
Copyright 2022 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook/metrics"
	"github.com/flant/addon-operator/sdk"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/deckhouse/deckhouse/ee/modules/110-istio/hooks/internal"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: internal.Queue("revisions-discovery-monitoring"),
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "services",
			ApiVersion: "v1",
			Kind:       "Service",
			FilterFunc: applyServiceFilter,
			// ignore d8 services, we don't use misconfigured external services
			LabelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "heritage",
						Operator: metav1.LabelSelectorOpNotIn,
						Values:   []string{"deckhouse"},
					},
				}},
		},
	},
}, handleExternalNameService)

func handleExternalNameService(input *go_hook.HookInput) error {
	input.MetricsCollector.Expire("d8_istio_service")
	snapshot := input.Snapshots["services"]

	for _, snap := range snapshot {
		if snap == nil {
			continue
		}

		service := snap.(externalService)

		input.MetricsCollector.Set("d8_istio_irrelevant_service", 1, map[string]string{"namespace": service.Namespace, "name": service.Name}, metrics.WithGroup("d8_istio_service"))
	}

	return nil
}

func applyServiceFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var service v1.Service

	err := sdk.FromUnstructured(obj, &service)
	if err != nil {
		return nil, err
	}

	if service.Spec.Type != v1.ServiceTypeExternalName {
		return nil, nil
	}

	if len(service.Spec.Ports) == 0 {
		return nil, nil
	}

	return externalService{
		Namespace: service.Namespace,
		Name:      service.Name,
	}, nil
}

type externalService struct {
	Namespace string
	Name      string
}
