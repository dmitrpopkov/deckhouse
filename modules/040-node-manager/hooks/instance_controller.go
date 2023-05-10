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
	"fmt"
	"time"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"

	mcmv1alpha1 "github.com/deckhouse/deckhouse/modules/040-node-manager/hooks/internal/mcm/v1alpha1"
	d8v1 "github.com/deckhouse/deckhouse/modules/040-node-manager/hooks/internal/v1"
	d8v1alpha1 "github.com/deckhouse/deckhouse/modules/040-node-manager/hooks/internal/v1alpha1"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Settings: &go_hook.HookConfigSettings{
		ExecutionMinInterval: 1 * time.Second,
		ExecutionBurst:       2,
	},
	Queue: "/modules/node-manager/instance_controller",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "instances",
			Kind:       "Instance",
			ApiVersion: "deckhouse.io/v1alpha1",
			FilterFunc: instanceFilter,
		},
		{
			Name:       "ngs",
			Kind:       "NodeGroup",
			ApiVersion: "deckhouse.io/v1",
			FilterFunc: instanceNodeGroupFilter,
		},
		{
			Name:       "machines",
			ApiVersion: "machine.sapcloud.io/v1alpha1",
			Kind:       "Machine",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-cloud-instance-manager"},
				},
			},
			FilterFunc: instanceMachineFilter,
		},
	},
}, instanceController)

func instanceMachineFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var machine mcmv1alpha1.Machine

	err := sdk.FromUnstructured(obj, &machine)
	if err != nil {
		return nil, err
	}

	return &machineForInstance{
		NodeGroup:         machine.Spec.NodeTemplateSpec.Labels["node.deckhouse.io/group"],
		NodeName:          machine.GetLabels()["node"],
		Name:              machine.GetName(),
		CurrentStatus:     machine.Status.CurrentStatus,
		LastOperation:     machine.Status.LastOperation,
		DeletionTimestamp: machine.GetDeletionTimestamp(),
	}, nil
}

var (
	deleteFinalizersPatch = map[string]interface{}{
		"metadata": map[string]interface{}{
			"finalizers": nil,
		},
	}
)

func instanceNodeGroupFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var ng d8v1.NodeGroup

	err := sdk.FromUnstructured(obj, &ng)
	if err != nil {
		return nil, err
	}

	return &nodeGroupForInstance{
		Name:           ng.GetName(),
		UID:            ng.GetUID(),
		ClassReference: ng.Spec.CloudInstances.ClassReference,
	}, nil
}

func instanceFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var ic d8v1alpha1.Instance

	err := sdk.FromUnstructured(obj, &ic)
	if err != nil {
		return nil, err
	}

	return &ic, nil
}

func instanceController(input *go_hook.HookInput) error {
	instances := make(map[string]*d8v1alpha1.Instance, len(input.Snapshots["instances"]))
	machines := make(map[string]*machineForInstance, len(input.Snapshots["machines"]))
	nodeGroups := make(map[string]*nodeGroupForInstance, len(input.Snapshots["ngs"]))

	for _, i := range input.Snapshots["instances"] {
		ins := i.(*d8v1alpha1.Instance)
		instances[ins.GetName()] = ins
	}

	for _, m := range input.Snapshots["machines"] {
		mc := m.(*machineForInstance)
		machines[mc.Name] = mc
	}

	for _, m := range input.Snapshots["ngs"] {
		ng := m.(*nodeGroupForInstance)
		nodeGroups[ng.Name] = ng
	}

	// first, check mapping from machines to instance
	// here we handle two cases:
	//   1. if instance exists, then we have 2 subcases:
	//  	a. first, we should compare machines fields and instance status and if we have changes path status
	//		b. next, we should check metadata.deletionTimestamp for instance
	//         if we have non-zero timestamp this means that client deletes instance, and we should delete machine.
	//         Before deletion, we should check metadata.deletionTimestamp in the machine for prevent multiple deletion
	//   2. if instance does not exist, create instance for machine
	for name, machine := range machines {
		ng, ok := nodeGroups[machine.NodeGroup]
		if !ok {
			return fmt.Errorf("NodeGroup %s not found", machine.NodeGroup)
		}

		if ic, ok := instances[name]; ok {
			statusPatch := getInstanceStatusPatch(ic, machine, ng)
			if len(statusPatch) > 0 {
				patch := map[string]interface{}{
					"status": statusPatch,
				}
				input.PatchCollector.MergePatch(patch, "deckhouse.io/v1alpha1", "Instance", "", ic.Name, object_patch.WithSubresource("/status"))
			}

			if ic.DeletionTimestamp != nil && !ic.DeletionTimestamp.IsZero() {
				if machine.DeletionTimestamp == nil || machine.DeletionTimestamp.IsZero() {
					// delete in background, because machine has finalizer
					input.PatchCollector.Delete("machine.sapcloud.io/v1alpha1", "Machine", "d8-cloud-instance-manager", machine.Name, object_patch.InBackground())
				}
			}
		} else {
			newIc := newInstance(machine, ng)
			input.PatchCollector.Create(newIc, object_patch.IgnoreIfExists())
		}
	}

	// next, check mapping from instance  to machines
	// here we handle next cases:
	//   1. if machine exists, then skip it, we handle it above
	//   2. if machine does not exist:
	//  	a. we should delete finalizers only if metadata.deletionTimestamp is non-zero
	//		b. we should delete finalizers and delete instance if metadata.deletionTimestamp is zero
	for _, ic := range instances {
		if _, ok := machines[ic.GetName()]; !ok {
			input.PatchCollector.MergePatch(deleteFinalizersPatch, "deckhouse.io/v1alpha1", "Instance", "", ic.Name)

			ds := ic.GetDeletionTimestamp()
			if ds == nil || ds.IsZero() {
				input.PatchCollector.Delete("deckhouse.io/v1alpha1", "Instance", "", ic.Name)
			}
		}
	}

	return nil
}

func newInstance(machine *machineForInstance, ng *nodeGroupForInstance) *d8v1alpha1.Instance {
	return &d8v1alpha1.Instance{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Instance",
			APIVersion: "deckhouse.io/v1alpha1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: machine.Name,
			Labels: map[string]string{
				"node.deckhouse.io/group": machine.NodeGroup,
			},

			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "deckhouse.io/v1",
					BlockOwnerDeletion: pointer.Bool(false),
					Controller:         pointer.Bool(false),
					Kind:               "NodeGroup",
					Name:               ng.Name,
					UID:                ng.UID,
				},
			},

			Finalizers: []string{
				"node-manager.hooks.deckhouse.io/instance-controller",
			},
		},

		Status: d8v1alpha1.InstanceStatus{
			ClassReference: d8v1alpha1.ClassReference{
				Name: ng.ClassReference.Name,
				Kind: ng.ClassReference.Kind,
			},
			MachineRef: d8v1alpha1.MachineRef{
				APIVersion: "machine.sapcloud.io/v1alpha1",
				Kind:       "Machine",
				Name:       machine.Name,
				Namespace:  "d8-cloud-instance-manager",
			},
			CurrentStatus: d8v1alpha1.CurrentStatus{
				Phase:          d8v1alpha1.InstancePhase(machine.CurrentStatus.Phase),
				LastUpdateTime: machine.CurrentStatus.LastUpdateTime,
			},
			LastOperation: d8v1alpha1.LastOperation{
				LastUpdateTime: machine.LastOperation.LastUpdateTime,
				Description:    machine.LastOperation.Description,
				State:          d8v1alpha1.State(machine.LastOperation.State),
				Type:           d8v1alpha1.OperationType(machine.LastOperation.Type),
			},
		},
	}
}

func instanceLastOpMap(s map[string]interface{}) map[string]interface{} {
	m, ok := s["lastOperation"]
	if !ok {
		m = make(map[string]interface{})
		s["lastOperation"] = m
	}

	return m.(map[string]interface{})
}

func getInstanceStatusPatch(ic *d8v1alpha1.Instance, machine *machineForInstance, ng *nodeGroupForInstance) map[string]interface{} {
	status := make(map[string]interface{})

	if ic.Status.NodeRef.Name != machine.NodeName {
		status["nodeRef"] = map[string]interface{}{
			"name": machine.NodeName,
		}
	}

	// if machine is Running we do not need bootstrap info
	if machine.CurrentStatus.Phase == mcmv1alpha1.MachineRunning && (ic.Status.BootstrapStatus.LogsEndpoint != "" || ic.Status.BootstrapStatus.Description != "") {
		status["bootstrapStatus"] = nil
	}

	if string(ic.Status.CurrentStatus.Phase) != string(machine.CurrentStatus.Phase) {
		status["currentStatus"] = map[string]interface{}{
			"phase":          string(machine.CurrentStatus.Phase),
			"lastUpdateTime": machine.CurrentStatus.LastUpdateTime.Format(time.RFC3339),
		}
	}

	if ic.Status.LastOperation.Description != machine.LastOperation.Description {
		m := instanceLastOpMap(status)
		m["description"] = machine.LastOperation.Description
	}

	if string(ic.Status.LastOperation.Type) != string(machine.LastOperation.Type) {
		m := instanceLastOpMap(status)
		m["type"] = string(machine.LastOperation.Type)
	}

	if string(ic.Status.LastOperation.State) != string(machine.LastOperation.State) {
		m := instanceLastOpMap(status)
		m["state"] = string(machine.LastOperation.State)
	}

	if !ic.Status.LastOperation.LastUpdateTime.Equal(&machine.LastOperation.LastUpdateTime) {
		m := instanceLastOpMap(status)
		m["lastUpdateTime"] = machine.LastOperation.LastUpdateTime.Format(time.RFC3339)
	}

	if ic.Status.ClassReference.Kind != ng.ClassReference.Kind || ic.Status.ClassReference.Name != ng.ClassReference.Name {
		status["classReference"] = map[string]interface{}{
			"kind": ng.ClassReference.Kind,
			"name": ng.ClassReference.Name,
		}
	}

	return status
}

type machineForInstance struct {
	NodeGroup         string
	NodeName          string
	Name              string
	CurrentStatus     mcmv1alpha1.CurrentStatus
	LastOperation     mcmv1alpha1.LastOperation
	DeletionTimestamp *metav1.Time
}

type nodeGroupForInstance struct {
	Name           string
	UID            k8stypes.UID
	ClassReference d8v1.ClassReference
}
