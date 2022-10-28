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

	"github.com/deckhouse/deckhouse/go_lib/certificate"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func etcdSecretFilter(unstructured *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var sec corev1.Secret

	err := sdk.FromUnstructured(unstructured, &sec)
	if err != nil {
		return nil, err
	}

	var cert certificate.Certificate
	if ca, ok := sec.Data["etcd-ca.crt"]; ok {
		cert.CA = string(ca)
		cert.Cert = string(ca)
	}
	if key, ok := sec.Data["etcd-ca.key"]; ok {
		cert.Key = string(key)
	}
	return &cert, nil
}

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "etcd-certificate",
			ApiVersion: "v1",
			Kind:       "Secret",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"kube-system"},
				},
			},
			NameSelector: &types.NameSelector{MatchNames: []string{"d8-pki"}},
			FilterFunc:   etcdSecretFilter,
		},
	},
}, storeSecret)

func storeSecret(input *go_hook.HookInput) error {
	snap := input.Snapshots["etcd-certificate"]

	if len(snap) == 0 {
		return fmt.Errorf("no etcd certificate found")
	}

	cert := snap[0].(*certificate.Certificate)
	input.Values.Set("monitoringKubernetesControlPlane.internal.etcdClientSecret", cert)

	return nil
}
