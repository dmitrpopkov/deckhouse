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

package destination

import "github.com/deckhouse/deckhouse/modules/460-log-shipper/apis/v1alpha1"

type Vector struct {
	CommonSettings

	Version string `json:"version,omitempty"`

	Address string `json:"address"`

	TLS VectorTLS `json:"tls,omitempty"`

	Keepalive VectorKeepalive `json:"keepalive,omitempty"`
}

type VectorTLS struct {
	CommonTLS         `json:",inline"`
	VerifyCertificate bool `json:"verify_certificate"`
	Enabled           bool `json:"enabled"`
}

type VectorKeepalive struct {
	TimeSecs int `json:"time_secs"`
}

func NewVector(name string, cspec v1alpha1.ClusterLogDestinationSpec) *Vector {
	spec := cspec.Vector

	// Disable buffer. It is buggy. Vector developers know about problems with buffer.
	// More info about buffer rewriting here - https://github.com/vectordotdev/vector/issues/9476
	// common.Buffer = buffer{
	//	Size: 100 * 1024 * 1024, // 100MiB in bytes for vector persistent queue
	//	Type: "disk",
	// }

	var enabledTLS bool
	if spec.TLS.KeyFile != "" || spec.TLS.CertFile != "" || spec.TLS.CAFile != "" {
		enabledTLS = true
	}

	return &Vector{
		CommonSettings: CommonSettings{
			Name: ComposeName(name),
			Type: "vector",
		},
		TLS: VectorTLS{
			CommonTLS: CommonTLS{
				CAFile:         decodeB64(spec.TLS.CAFile),
				CertFile:       decodeB64(spec.TLS.CertFile),
				KeyFile:        decodeB64(spec.TLS.KeyFile),
				KeyPass:        decodeB64(spec.TLS.KeyPass),
				VerifyHostname: spec.TLS.VerifyHostname,
			},
			VerifyCertificate: spec.TLS.VerifyCertificate,
			Enabled:           enabledTLS,
		},
		Version: "2",
		Address: spec.Endpoint,
		// TODO(nabokikhms): Only available for vector the first version sink, consider different load balancing solution
		//
		// Keepalive: VectorKeepalive{
		//	TimeSecs: 7200,
		// },
	}
}
