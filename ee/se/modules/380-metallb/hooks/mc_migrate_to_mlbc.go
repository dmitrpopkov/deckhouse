/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"fmt"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 5},
	Queue:        "/modules/metallb/node-labels-update",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "metallb_module_config",
			ApiVersion: "deckhouse.io/v1alpha1",
			Kind:       "ModuleConfig",
			NameSelector: &types.NameSelector{
				MatchNames: []string{"metallb"},
			},
			FilterFunc: applyModuleConfigFilter,
		},
		{
			Name:       "l2advertisements",
			ApiVersion: "metallb.io/v1beta1",
			Kind:       "L2Advertisement",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-metallb"},
				},
			},
			FilterFunc: applyL2AdvertisementFilter,
		},
		{
			Name:       "ipaddresspools",
			ApiVersion: "metallb.io/v1beta1",
			Kind:       "IPAddressPool",
			FilterFunc: applyIPAddressPoolFilter,
		},
		// {
		// 	Name:       "mlbc",
		// 	ApiVersion: "network.deckhouse.io/v1alpha1",
		// 	Kind:       "MetalLoadBalancerClass",
		// 	FilterFunc: applyMetalLoadBalancerClassFilter,
		// },
	},
}, migrateMCtoMLBC)

func applyModuleConfigFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	mc := &ModuleConfig{}
	err := sdk.FromUnstructured(obj, mc)
	if err != nil {
		return nil, fmt.Errorf("cannot convert Metallb ModuleConfig: %v", err)
	}

	if mc.Spec.Version == 1 {
		return mc, nil
	}
	return nil, nil
}

func applyL2AdvertisementFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	l2Advertisement := &L2Advertisement{}
	err := sdk.FromUnstructured(obj, l2Advertisement)
	if err != nil {
		return nil, err
	}

	if len(l2Advertisement.Spec.IPAddressPools) == 0 {
		return nil, nil
	}

	return L2AdvertisementInfo{
		Name:           l2Advertisement.Name,
		IPAddressPools: l2Advertisement.Spec.IPAddressPools,
		NodeSelectors:  l2Advertisement.Spec.NodeSelectors,
		Namespace:      l2Advertisement.Namespace,
	}, nil
}

func applyIPAddressPoolFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	ipAddressPool := &IPAddressPool{}
	err := sdk.FromUnstructured(obj, ipAddressPool)
	if err != nil {
		return nil, err
	}

	return IPAddressPoolInfo{
		Name:      ipAddressPool.Name,
		Addresses: ipAddressPool.Spec.Addresses,
	}, err
}

func createMetalLoadBalancerClass(input *go_hook.HookInput, mlbcInfo *MetalLoadBalancerClassInfo) {
	mlbc := map[string]interface{}{
		"apiVersion": "network.deckhouse.io/v1alpha1",
		"kind":       "MetalLoadBalancerClass",
		"metadata": map[string]interface{}{
			"name": mlbcInfo.Name,
		},
		"spec": map[string]interface{}{
			"isDefault":    mlbcInfo.IsDefault,
			"type":         "L2",
			"addressPool":  mlbcInfo.AddressPool,
			"nodeSelector": mlbcInfo.NodeSelector,
			"tolerations":  mlbcInfo.Tolerations,
		},
	}
	mlbcUnstructured, err := sdk.ToUnstructured(&mlbc)
	if err != nil {
		return
	}
	input.PatchCollector.Create(mlbcUnstructured, object_patch.IgnoreIfExists())
}

func migrateMCtoMLBC(input *go_hook.HookInput) error {
	snapsMC := input.Snapshots["metallb_module_config"]
	if len(snapsMC) != 1 || snapsMC[0] == nil {
		return nil
	}
	mc, ok := snapsMC[0].(*ModuleConfig)
	if !ok {
		return nil
	}

	var mlbcDefault MetalLoadBalancerClassInfo
	mlbcDefault.Name = "l2-default"
	mlbcDefault.IsDefault = true

	// Getting addressPools, nodeSelector and tolerations from MC
	if len(mc.Spec.Settings.AddressPools) > 0 {
		for _, addressPool := range mc.Spec.Settings.AddressPools {
			if addressPool.Protocol != "layer2" {
				continue
			}
			mlbcDefault.AddressPool = append(mlbcDefault.AddressPool, addressPool.Addresses...)
		}
	}
	if mc.Spec.Settings.Speaker.NodeSelector != nil {
		mlbcDefault.NodeSelector = mc.Spec.Settings.Speaker.NodeSelector
	}
	mlbcDefault.Tolerations = mc.Spec.Settings.Speaker.Tolerations

	// Getting nodeSelector from L2Advertisement, create Default MetalLoadBalancerClass
	snapsL2A := input.Snapshots["l2advertisements"]
	for _, snapL2A := range snapsL2A {
		l2Advertisement := snapL2A.(L2AdvertisementInfo)
		if len(l2Advertisement.NodeSelectors) > 0 {
			for k, v := range l2Advertisement.NodeSelectors[0].MatchLabels {
				mlbcDefault.NodeSelector[k] = v
			}
		}
	}
	createMetalLoadBalancerClass(input, &mlbcDefault)

	// Collect ipAddressPoolAddresses from IPAddressPools
	ipAddressPollsMap := make(map[string][]string)
	snapsIAP := input.Snapshots["ipaddresspools"]
	for _, snapIAP := range snapsIAP {
		ipAddressPool := snapIAP.(IPAddressPoolInfo)
		ipAddressPollsMap[ipAddressPool.Name] = ipAddressPool.Addresses
	}

	var mlbc MetalLoadBalancerClassInfo
	mlbc.Labels = make(map[string]string)
	mlbc.Labels["heritage"] = "deckhouse"
	mlbc.Labels["auto-generated-by"] = "migration-hook"
	mlbc.IsDefault = false

	// Create other MetalLoadBalancerClasses
	l2AdvertisementNamesMap := make(map[string]bool)
	snapsL2A = input.Snapshots["l2advertisements"]
	for _, snapL2A := range snapsL2A {
		l2Advertisement := snapL2A.(L2AdvertisementInfo)
		l2AdvertisementNamesMap[l2Advertisement.Name] = true // Will be used to remove the MLBC
		for _, ipAddressPool := range l2Advertisement.IPAddressPools {
			if addresses, ok := ipAddressPollsMap[ipAddressPool]; ok {
				mlbc.AddressPool = append(mlbc.AddressPool, addresses...)
			}
		}
		mlbc.Name = l2Advertisement.Name
		createMetalLoadBalancerClass(input, &mlbc)
	}
	return nil
}
