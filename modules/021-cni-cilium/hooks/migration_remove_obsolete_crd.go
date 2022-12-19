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
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnStartup: &go_hook.OrderedConfig{Order: 0},
}, removeObsoleteCiliumCRD)

func removeObsoleteCiliumCRD(input *go_hook.HookInput) error {
	input.PatchCollector.Delete("apiextensions.k8s.io/v1", "CustomResourceDefinition", "", "ciliumenvoyconfigs.cilium.io", object_patch.InForeground())
	input.PatchCollector.Delete("apiextensions.k8s.io/v1", "CustomResourceDefinition", "", "ciliumclusterwideenvoyconfigs.cilium.io", object_patch.InForeground())
	return nil
}
