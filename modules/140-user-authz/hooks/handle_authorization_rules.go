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

	"github.com/deckhouse/deckhouse/modules/140-user-authz/hooks/internal"
)

const (
	authRuleSnapshot = "authorization_rules"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: internal.Queue(authRuleSnapshot),
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       authRuleSnapshot,
			ApiVersion: "deckhouse.io/v1alpha1",
			Kind:       "AuthorizationRule",
			FilterFunc: internal.ApplyAuthorizationRuleFilter,
		},
	},
}, internal.AuthorizationRulesHandler("userAuthz.internal.authRuleCrds", authRuleSnapshot))
