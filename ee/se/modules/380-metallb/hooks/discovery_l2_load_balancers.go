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
}, handleLoadBalancers)

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

	var DesiredIPs []string
	if DesiredIPsStr, ok := service.Annotations[l2LoadBalancerIPsAnnotate]; ok {
		DesiredIPs = strings.Split(DesiredIPsStr, ",")
	}

	var loadBalancerClass string
	if service.Spec.LoadBalancerClass != nil {
		loadBalancerClass = *service.Spec.LoadBalancerClass
	}

	var assignedLoadBalancerClass string
	for _, condition := range service.Status.Conditions {
		if condition.Type == "network.deckhouse.io/LoadBalancerClass" {
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
		DesiredIPs:                DesiredIPs,
	}, nil
}

func applyMetalLoadBalancerClassFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var metalloadbalancerclass MetalLoadBalancerClass

	err := sdk.FromUnstructured(obj, &metalloadbalancerclass)
	if err != nil {
		return nil, err
	}

	interfaces := []string{}
	if len(metalloadbalancerclass.Spec.L2.Interfaces) > 0 {
		interfaces = metalloadbalancerclass.Spec.L2.Interfaces
	}

	return MetalLoadBalancerClassInfo{
		Name:         metalloadbalancerclass.Name,
		AddressPool:  metalloadbalancerclass.Spec.AddressPool,
		Interfaces:   interfaces,
		NodeSelector: metalloadbalancerclass.Spec.NodeSelector,
		IsDefault:    metalloadbalancerclass.Spec.IsDefault,
	}, nil
}

func handleLoadBalancers(input *go_hook.HookInput) error {
	l2lbservices := make([]L2LBServiceConfig, 0, 4)
	metalLoadBalancerClasses, mlbcDefaultName := makeMLBCMapFromSnapshot(input.Snapshots["mlbc"])

	for _, serviceSnap := range input.Snapshots["services"] {
		service, ok := serviceSnap.(ServiceInfo)
		if !ok {
			continue
		}

		var mlbc MetalLoadBalancerClassInfo
		var needPatchService bool
		if mlbc, ok = metalLoadBalancerClasses[service.AssignedLoadBalancerClass]; ok {
			needPatchService = false
		} else if service.AssignedLoadBalancerClass != "" {
			// AssignedLoadBalancerClass is missing from existing MetalLoadBalancerClasses
			continue
		} else if mlbc, ok = metalLoadBalancerClasses[service.LoadBalancerClass]; ok {
			needPatchService = true
		} else if mlbc, ok = metalLoadBalancerClasses[mlbcDefaultName]; ok {
			needPatchService = true
		} else {
			continue
		}

		if needPatchService {
			patch := map[string]interface{}{
				"status": map[string]interface{}{
					"conditions": []metav1.Condition{
						{
							Type:    "network.deckhouse.io/LoadBalancerClass",
							Message: mlbc.Name,
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

		nodes := getNodesByMLBC(mlbc, input.Snapshots["nodes"])
		if len(nodes) == 0 {
			// There is no node that matches the specified node selector.
			continue
		}

		loopCount := len(service.DesiredIPs)
		isMigrateService := loopCount > 0
		if !isMigrateService {
			loopCount = service.ExternalIPsCount
		}
		for i := 0; i < loopCount; i++ {
			nodeIndex := i % len(nodes)
			config := L2LBServiceConfig{
				Name:                       fmt.Sprintf("%s-%s-%d", service.Name, mlbc.Name, i),
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
				MetalLoadBalancerClassName: mlbc.Name,
			}
			if isMigrateService {
				config.DesiredIP = service.DesiredIPs[i]
			}
			l2lbservices = append(l2lbservices, config)
		}
	}

	// L2 Load Balancers are sorted before saving
	l2loadbalancersInternal := make([]MetalLoadBalancerClassInfo, 0, len(metalLoadBalancerClasses))
	for _, value := range metalLoadBalancerClasses {
		l2loadbalancersInternal = append(l2loadbalancersInternal, value)
	}
	sort.Slice(l2loadbalancersInternal, func(i, j int) bool {
		return l2loadbalancersInternal[i].Name < l2loadbalancersInternal[j].Name
	})
	input.Values.Set("metallb.internal.l2loadbalancers", l2loadbalancersInternal)

	// L2 Load Balancer Services are sorted by Namespace and then Name before saving
	sort.Slice(l2lbservices, func(i, j int) bool {
		if l2lbservices[i].Namespace == l2lbservices[j].Namespace {
			return l2lbservices[i].Name < l2lbservices[j].Name
		}
		return l2lbservices[i].Namespace < l2lbservices[j].Namespace
	})
	input.Values.Set("metallb.internal.l2lbservices", l2lbservices)
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
