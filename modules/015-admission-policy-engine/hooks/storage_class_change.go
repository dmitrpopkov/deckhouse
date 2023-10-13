/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import "github.com/deckhouse/deckhouse/go_lib/hooks/storage_class_change"

var _ = storage_class_change.RegisterHook(storage_class_change.Args{
        ModuleName:         "admission-policy-engine",
        Namespace:          "d8-admission-policy-engine",
        LabelSelectorKey:   "app",
        LabelSelectorValue: "trivy-provider",
        ObjectKind:         "StatefulSet",
        ObjectName:         "trivy-provider",
})

