/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
	Queue:        "/modules/metallb/discovery",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "mlbc",
			ApiVersion: "network.deckhouse.io/v1alpha1",
			Kind:       "MetalLoadBalancerClass",
			FilterFunc: applyMetalLoadBalancerClassFilter,
		},
		{
			Name:       "nodes",
			ApiVersion: "v1",
			Kind:       "Node",
			FilterFunc: applyNodeFilter,
		},
		{
			Name:       "services",
			ApiVersion: "v1",
			Kind:       "Service",
			FilterFunc: applyServiceFilter,
		},
	},
}, handleL2LoadBalancers)

func applyNodeFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var node v1.Node

	err := sdk.FromUnstructured(obj, &node)
	if err != nil {
		return nil, err
	}

	_, isLabeled := node.Labels[memberLabelKey]

	return NodeInfo{
		Name:      node.Name,
		Labels:    node.Labels,
		IsLabeled: isLabeled,
	}, nil
}

func applyServiceFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var service v1.Service

	err := sdk.FromUnstructured(obj, &service)
	if err != nil {
		return nil, err
	}

	if service.Spec.Type != v1.ServiceTypeLoadBalancer {
		// we only need service of LoadBalancer type
		return nil, nil
	}

	var externalIPsCount = 1
	if externalIPsCountStr, ok := service.Annotations[keyAnnotationExternalIPsCount]; ok {
		if externalIP, err := strconv.Atoi(externalIPsCountStr); err == nil {
			if externalIP > 1 {
				externalIPsCount = externalIP
			}
		}
	}

	var desiredIPs []string
	if DesiredIPsStr, ok := service.Annotations[l2LoadBalancerIPsAnnotate]; ok {
		desiredIPs = strings.Split(DesiredIPsStr, ",")
	}

	var lbAllowSharedIP string
	if lbAllowSharedIPStr, ok := service.Annotations[lbAllowSharedIPAnnotate]; ok {
		lbAllowSharedIP = lbAllowSharedIPStr
	}

	var loadBalancerClass string
	if service.Spec.LoadBalancerClass != nil {
		loadBalancerClass = *service.Spec.LoadBalancerClass
	}

	var assignedLoadBalancerClass string
	for _, condition := range service.Status.Conditions {
		if condition.Type == "network.deckhouse.io/load-balancer-class" {
			assignedLoadBalancerClass = condition.Message
			continue
		}
	}

	internalTrafficPolicy := v1.ServiceInternalTrafficPolicyCluster
	if service.Spec.InternalTrafficPolicy != nil {
		internalTrafficPolicy = *service.Spec.InternalTrafficPolicy
	}

	return ServiceInfo{
		Name:                      service.GetName(),
		Namespace:                 service.GetNamespace(),
		LoadBalancerClass:         loadBalancerClass,
		AssignedLoadBalancerClass: assignedLoadBalancerClass,
		ExternalIPsCount:          externalIPsCount,
		Ports:                     service.Spec.Ports,
		ExternalTrafficPolicy:     service.Spec.ExternalTrafficPolicy,
		InternalTrafficPolicy:     internalTrafficPolicy,
		Selector:                  service.Spec.Selector,
		ClusterIP:                 service.Spec.ClusterIP,
		PublishNotReadyAddresses:  service.Spec.PublishNotReadyAddresses,
		DesiredIPs:                desiredIPs,
		LBAllowSharedIP:           lbAllowSharedIP,
	}, nil
}

func applyMetalLoadBalancerClassFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var metalLoadBalancerClass MetalLoadBalancerClass

	err := sdk.FromUnstructured(obj, &metalLoadBalancerClass)
	if err != nil {
		return nil, err
	}

	interfaces := []string{}
	if len(metalLoadBalancerClass.Spec.L2.Interfaces) > 0 {
		interfaces = metalLoadBalancerClass.Spec.L2.Interfaces
	}

	return MetalLoadBalancerClassInfo{
		Name:         metalLoadBalancerClass.Name,
		AddressPool:  metalLoadBalancerClass.Spec.AddressPool,
		Interfaces:   interfaces,
		NodeSelector: metalLoadBalancerClass.Spec.NodeSelector,
		IsDefault:    metalLoadBalancerClass.Spec.IsDefault,
	}, nil
}

func handleL2LoadBalancers(input *go_hook.HookInput) error {
	l2LBServices := make([]L2LBServiceConfig, 0, 4)
	mlbcMap, mlbcDefaultName := makeMLBCMapFromSnapshot(input.Snapshots["mlbc"])

	for _, serviceSnap := range input.Snapshots["services"] {
		service, ok := serviceSnap.(ServiceInfo)
		if !ok {
			continue
		}

		patchStatusInformation := true
		var mlbcForUse MetalLoadBalancerClassInfo
		if mlbcTemp, ok := mlbcMap[mlbcDefaultName]; ok {
			// Use default MLBC (and add to status)
			mlbcForUse = mlbcTemp
		}
		if mlbcTemp, ok := mlbcMap[service.LoadBalancerClass]; ok {
			// Else use the MLBC that exists in the cluster (and add to status)
			mlbcForUse = mlbcTemp
		}
		if service.AssignedLoadBalancerClass != "" {
			if mlbcTemp, ok := mlbcMap[service.AssignedLoadBalancerClass]; ok {
				// Else use the MLBC associated earlier
				mlbcForUse = mlbcTemp
				patchStatusInformation = false
			} else {
				// MLBC is not among clustered MLBCs, but it is associated (in status)
				continue
			}
		}
		if mlbcForUse.Name == "" {
			continue
		}

		if patchStatusInformation {
			patch := map[string]interface{}{
				"status": map[string]interface{}{
					"conditions": []metav1.Condition{
						{
							Type:    "network.deckhouse.io/load-balancer-class",
							Message: mlbcForUse.Name,
							Status:  "True",
							Reason:  "LoadBalancerClassBound",
						},
					},
				},
			}

			input.PatchCollector.MergePatch(patch,
				"v1",
				"Service",
				service.Namespace,
				service.Name,
				object_patch.WithSubresource("/status"))
		}

		nodes := getNodesByMLBC(mlbcForUse, input.Snapshots["nodes"])
		if len(nodes) == 0 {
			// There is no node that matches the specified node selector.
			continue
		}

		desiredIPsCount := len(service.DesiredIPs)
		desiredIPsExist := desiredIPsCount > 0
		for i := 0; i < service.ExternalIPsCount; i++ {
			nodeIndex := i % len(nodes)
			config := L2LBServiceConfig{
				Name:                       fmt.Sprintf("%s-%s-%d", service.Name, mlbcForUse.Name, i),
				Namespace:                  service.Namespace,
				ServiceName:                service.Name,
				ServiceNamespace:           service.Namespace,
				PreferredNode:              nodes[nodeIndex].Name,
				ExternalTrafficPolicy:      service.ExternalTrafficPolicy,
				InternalTrafficPolicy:      service.InternalTrafficPolicy,
				PublishNotReadyAddresses:   service.PublishNotReadyAddresses,
				ClusterIP:                  service.ClusterIP,
				Ports:                      service.Ports,
				Selector:                   service.Selector,
				MetalLoadBalancerClassName: mlbcForUse.Name,
				LBAllowSharedIP:            service.LBAllowSharedIP,
			}
			if desiredIPsExist && i < desiredIPsCount {
				config.DesiredIP = service.DesiredIPs[i]
			}
			l2LBServices = append(l2LBServices, config)
		}
	}

	// L2 Load Balancers are sorted before saving
	l2LoadBalancersInternal := make([]MetalLoadBalancerClassInfo, 0, len(mlbcMap))
	for _, value := range mlbcMap {
		l2LoadBalancersInternal = append(l2LoadBalancersInternal, value)
	}
	sort.Slice(l2LoadBalancersInternal, func(i, j int) bool {
		return l2LoadBalancersInternal[i].Name < l2LoadBalancersInternal[j].Name
	})
	input.Values.Set("metallb.internal.l2loadbalancers", l2LoadBalancersInternal)

	// L2 Load Balancer Services are sorted by Namespace and then Name before saving
	sort.Slice(l2LBServices, func(i, j int) bool {
		if l2LBServices[i].Namespace == l2LBServices[j].Namespace {
			return l2LBServices[i].Name < l2LBServices[j].Name
		}
		return l2LBServices[i].Namespace < l2LBServices[j].Namespace
	})
	input.Values.Set("metallb.internal.l2lbservices", l2LBServices)
	return nil
}

func makeMLBCMapFromSnapshot(snapshot []go_hook.FilterResult) (map[string]MetalLoadBalancerClassInfo, string) {
	mlbcMap := make(map[string]MetalLoadBalancerClassInfo)
	var mlbcDefaultName string
	for _, mlbcSnap := range snapshot {
		if mlbc, ok := mlbcSnap.(MetalLoadBalancerClassInfo); ok {
			mlbcMap[mlbc.Name] = mlbc
			if mlbc.IsDefault {
				mlbcDefaultName = mlbc.Name
			}
		}
	}
	return mlbcMap, mlbcDefaultName
}
