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
	"fmt"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"

	hooks "github.com/deckhouse/deckhouse/go_lib/hooks/internal_tls"
)

const (
	webhookServiceHost = "webhook-handler.d8-system.svc"
)

var _ = hooks.RegisterInternalTLSHook(hooks.GenSelfSignedTLSHookConf{
	SANsGenerator: func(input *go_hook.HookInput) []string {
		webhookServiceFQDN := fmt.Sprintf(
			"%s.%s",
			webhookServiceHost,
			input.Values.Get("global.discovery.clusterDomain").String(),
		)

		return []string{
			webhookServiceHost,
			webhookServiceFQDN,
			"validating-" + webhookServiceHost,
			"conversion-" + webhookServiceHost,
			"validating-" + webhookServiceFQDN,
			"conversion-" + webhookServiceFQDN,
		}
	},

	CN: "webhook-handler.d8-system.svc",

	Namespace:     "d8-system",
	TLSSecretName: "webhook-handler-certs",

	CertValuesPath: "deckhouse.internal.webhookHandlerCert.crt",
	CAValuesPath:   "deckhouse.internal.webhookHandlerCert.ca",
	KeyValuesPath:  "deckhouse.internal.webhookHandlerCert.key",
})
