/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/deckhouse/deckhouse/go_lib/dependency"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnStartup: &go_hook.OrderedConfig{Order: 5},
	Queue:     "/modules/metallb/discovery-services",
}, dependency.WithExternalDependencies(discoveryServicesForMigrate))

// TODO: Disable logging
func discoveryServicesForMigrate(input *go_hook.HookInput, dc dependency.Container) error {
	input.LogEntry.Infoln("MMMLB: Start a hook")
	k8sClient, err := dc.GetK8sClient()
	if err != nil {
		return err
	}
	serviceList, err := k8sClient.CoreV1().Services("").List(
		context.TODO(),
		metav1.ListOptions{},
	)
	if err != nil {
		return nil
	}
	input.LogEntry.Infoln("MMMLB: Start a loop")
	for _, service := range serviceList.Items {
		// Is it not a Loadbalancer?
		if service.Spec.Type != "LoadBalancer" {
			continue
		}
		input.LogEntry.Infof("MMMLB: service.Spec.Type=%v\n", service.Spec.Type)

		// Has MLBC status?
		statusExists := false
		for _, condition := range service.Status.Conditions {
			if condition.Type == "network.deckhouse.io/load-balancer-class" {
				statusExists = true
				input.LogEntry.Infof("MMMLB: condition.Type=%v\n", condition.Type)
				break
			}
		}
		if statusExists {
			continue
		}

		// Has the annotations?
		if _, ok := service.ObjectMeta.Annotations["network.deckhouse.io/l2-load-balancer-ips"]; ok {
			input.LogEntry.Infoln("MMMLB: First annotation")
			continue
		}
		if _, ok := service.ObjectMeta.Annotations["metallb.universe.tf/address-pool"]; !ok {
			input.LogEntry.Infoln("MMMLB: Second annotation")
			continue
		}

		input.LogEntry.Infoln("MMMLB: Patching!")
		var sliceIPs []string
		for _, ingress := range service.Status.LoadBalancer.Ingress {
			sliceIPs = append(sliceIPs, ingress.IP)
		}
		stringIPs := strings.Join(sliceIPs, ",")

		// Patch the service
		patch := map[string]interface{}{
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"network.deckhouse.io/l2-load-balancer-ips": stringIPs,
				},
			},
		}

		data, err := json.Marshal(patch)
		if err != nil {
			return fmt.Errorf("error to marshal patch-object: %w", err)
		}

		_, err = k8sClient.CoreV1().Services(service.ObjectMeta.Namespace).Patch(
			context.TODO(),
			service.ObjectMeta.Name,
			types.MergePatchType,
			data,
			metav1.PatchOptions{},
		)
		if err != nil {
			return fmt.Errorf("error to apply patch to Service %s: %w", service.ObjectMeta.Name, err)
		}
	}
	return nil
}
