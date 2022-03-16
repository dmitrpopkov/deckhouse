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
	"time"

	"github.com/cloudflare/cfssl/csr"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/pointer"

	"github.com/deckhouse/deckhouse/go_lib/certificate"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 5},
	Queue:        "/modules/cni-cilium/gen-cert",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "hubble-server-cert-secret",
			ApiVersion: "v1",
			Kind:       "Secret",
			NameSelector: &types.NameSelector{
				MatchNames: []string{"hubble-server-certs"},
			},
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-cni-cilium"},
				},
			},
			ExecuteHookOnSynchronization: pointer.BoolPtr(false),
			ExecuteHookOnEvents:          pointer.BoolPtr(false),
			FilterFunc:                   filterAdmissionSecret,
		},
	},
}, generateHubbleServerCert)

func filterAdmissionSecret(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var sec v1.Secret

	err := sdk.FromUnstructured(obj, &sec)
	if err != nil {
		return nil, err
	}

	return certificate.Certificate{
		CA:   string(sec.Data["ca.crt"]),
		Cert: string(sec.Data["tls.crt"]),
		Key:  string(sec.Data["tls.key"]),
	}, nil
}

func generateHubbleServerCert(input *go_hook.HookInput) error {
	snap := input.Snapshots["hubble-server-cert-secret"]

	if len(snap) > 0 {
		adm := snap[0].(certificate.Certificate)
		input.Values.Set("cniCilium.internal.hubbleServer.cert", adm.Cert)
		input.Values.Set("cniCilium.internal.hubbleServer.key", adm.Key)
		input.Values.Set("cniCilium.internal.hubbleServer.ca", adm.CA)

		return nil
	}

	const cn = "d8.hubble-ca.cilium.io"
	ca, err := certificate.GenerateCA(input.LogEntry, cn,
		certificate.WithKeyRequest(&csr.KeyRequest{
			A: "rsa",
			S: 2048,
		}),
		certificate.WithSANs("*.d8-cni-cilium.hubble-grpc.cilium.io"),
	)
	if err != nil {
		return errors.Wrap(err, "generate CA failed")
	}

	tls, err := certificate.GenerateSelfSignedCert(input.LogEntry,
		cn,
		ca,
		certificate.WithKeyRequest(&csr.KeyRequest{
			A: "rsa",
			S: 2048,
		}),
		certificate.WithSANs("*.d8-cni-cilium.hubble-grpc.cilium.io"),
		certificate.WithSigningDefaultExpiry(87600*time.Hour),
	)
	if err != nil {
		return errors.Wrap(err, "generate Cert failed")
	}

	input.Values.Set("cniCilium.internal.hubbleServer.cert", tls.Cert)
	input.Values.Set("cniCilium.internal.hubbleServer.key", tls.Key)
	input.Values.Set("cniCilium.internal.hubbleServer.ca", tls.CA)

	return nil
}
