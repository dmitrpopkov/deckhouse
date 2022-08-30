/*
Copyright 2022 Flant JSC

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

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	mTLSSwitchPath = "monitoringCustom.internal.prometheusScraperIstioMTLSEnabled"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: "/modules/monitoring-custom/switch_istio_mtls",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "prometheus_secret_mtls",
			ApiVersion: "v1",
			Kind:       "Secret",
			FilterFunc: applySecertMTLSFilter,
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-monitoring"},
				},
			},
			NameSelector: &types.NameSelector{
				MatchNames: []string{"prometheus-scraper-istio-mtls"},
			},
		},
	},
}, istioMTLSswitchHook)

func applySecertMTLSFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	secret := &v1.Secret{}
	err := sdk.FromUnstructured(obj, secret)

	if err != nil {
		return nil, fmt.Errorf("can't convert ca secret to secret struct: %v", err)
	}
	return true, nil
}

func istioMTLSswitchHook(input *go_hook.HookInput) error {
	mTLSCertSnap := input.Snapshots["prometheus_secret_mtls"]
	if len(mTLSCertSnap) == 1 {
		input.Values.Set(mTLSSwitchPath, true)
	} else {
		input.Values.Set(mTLSSwitchPath, false)
	}
	return nil
}
