/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

// import (
// 	"encoding/base64"
// 	"fmt"
// 	"sort"

// 	resources_common "github.com/deckhouse/deckhouse/ee/modules/038-system-registry/hooks/internal/resources/common"
// 	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
// 	"github.com/flant/addon-operator/sdk"
// 	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
// 	v1_core "k8s.io/api/core/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// )

// var _ = sdk.RegisterFunc(&go_hook.HookConfig{
// 	Queue:        "/modules/system-registry/pki",
// 	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
// 	Kubernetes: []go_hook.KubernetesConfig{
// 		{
// 			Name:       "pki",
// 			ApiVersion: "v1",
// 			Kind:       "Secret",
// 			NamespaceSelector: &types.NamespaceSelector{
// 				NameSelector: &types.NameSelector{
// 					MatchNames: []string{"kube-system"},
// 				},
// 			},
// 			NameSelector: &types.NameSelector{
// 				MatchNames: []string{"d8-pki"},
// 			},
// 			FilterFunc: filterPkiSecret,
// 		},
// 	},
// }, handlePKIData)

// func filterPkiSecret(unstructured *unstructured.Unstructured) (go_hook.FilterResult, error) {
// 	var sec v1_core.Secret

// 	err := sdk.FromUnstructured(unstructured, &sec)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resources_common.SecretData(sec.Data), nil
// }

// func handlePKIData(input *go_hook.HookInput) error {
// 	snap := input.Snapshots["pki"]

// 	if len(snap) == 0 {
// 		return fmt.Errorf(`there is no Secret named "d8-pki" in NS "kube-system"`)
// 	}

// 	sData := snap[0].(resources_common.SecretData)

// 	keys := make([]string, 0, len(sData))
// 	kvSData := make([]resources_common.SecretDataKV, 0, len(sData))

// 	// sort map values by key
// 	for k := range sData {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)

// 	// create kv sData
// 	for _, key := range keys {
// 		kvSData = append(kvSData, resources_common.SecretDataKV{Key: key, Value: base64.StdEncoding.EncodeToString(sData[key])})
// 	}

// 	input.Values.Set("systemRegistry.internal.pki.data", kvSData)
// 	return nil
// }
