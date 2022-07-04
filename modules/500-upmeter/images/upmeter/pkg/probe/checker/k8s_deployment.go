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

package checker

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"d8.io/upmeter/pkg/check"
	k8s "d8.io/upmeter/pkg/kubernetes"
	"d8.io/upmeter/pkg/probe/run"
)

const agentLabelKey = "upmeter-agent"

// DeploymentLifecycle is a checker constructor and configurator
type DeploymentLifecycle struct {
	Access                    k8s.Access
	Namespace                 string
	DeploymentCreationTimeout time.Duration
	DeploymentDeletionTimeout time.Duration
	PodAppearTimeout          time.Duration
	PodDisappearTimeout       time.Duration
	GarbageCollectionTimeout  time.Duration
	ControlPlaneAccessTimeout time.Duration
}

func (c DeploymentLifecycle) Checker() check.Checker {
	return &deploymentLifecycleChecker{
		agentId:                   run.ID(),
		access:                    c.Access,
		namespace:                 c.Namespace,
		deploymentCreationTimeout: c.DeploymentCreationTimeout,
		deploymentDeletionTimeout: c.DeploymentDeletionTimeout,
		podAppearTimeout:          c.PodAppearTimeout,
		podDisappearTimeout:       c.PodDisappearTimeout,
		garbageCollectionTimeout:  c.GarbageCollectionTimeout,
		controlPlaneAccessTimeout: c.ControlPlaneAccessTimeout,
	}
}

type deploymentLifecycleChecker struct {
	agentId                   string
	access                    k8s.Access
	namespace                 string
	deploymentCreationTimeout time.Duration
	deploymentDeletionTimeout time.Duration
	podAppearTimeout          time.Duration
	podDisappearTimeout       time.Duration

	garbageCollectionTimeout  time.Duration
	controlPlaneAccessTimeout time.Duration
}

func (c *deploymentLifecycleChecker) Check() check.Error {
	deployment := createDeploymentObject(c.agentId)
	checker := c.new(deployment)
	return checker.Check()
}

/*
 1. check control plane availability
 2. collect the garbage of the deployment/rs/pods from previous runs
 3. create the deployment in api
 4. wait for pending pod to appear
 5. delete the deployment in api
 6. wait for the pod to disappear
*/
func (c *deploymentLifecycleChecker) new(deployment *appsv1.Deployment) check.Checker {
	name := deployment.GetName()

	pingControlPlane := newControlPlaneChecker(c.access, c.controlPlaneAccessTimeout)

	// Clean all prior garbage that could be left by agent restarts. We rely on agent ID in
	// assumption that master nodes are not a subject for renaming.
	labels := map[string]string{agentLabelKey: c.agentId}
	collectGarbage := newGarbageCollectorCheckerByLabels(c.access, deployment.Kind, c.namespace, labels, c.garbageCollectionTimeout)

	createDeploymentOrUnknown := doOrUnknown(
		c.deploymentCreationTimeout,
		&deploymentCreator{
			access:     c.access,
			namespace:  c.namespace,
			deployment: deployment,
		},
	)

	deleteDeployment := &deploymentDeleter{
		access:    c.access,
		namespace: c.namespace,
		name:      name,
	}

	// Track pods only created by current deployment since the deployment name consists of agent
	// ID and random tail.
	podListOptions := listOptsByLabels(map[string]string{"deployment": name})

	verifyPodExists := withRetryEachSeconds(
		&pendingPodChecker{
			access:    c.access,
			namespace: c.namespace,
			listOpts:  podListOptions,
		},
		c.podAppearTimeout)

	verifyNoPod := withRetryEachSeconds(
		&objectIsNotListedChecker{
			access:    c.access,
			namespace: c.namespace,
			kind:      "Pod",
			listOpts:  podListOptions,
		},
		c.podDisappearTimeout)

	return sequence(
		pingControlPlane,
		collectGarbage,
		createDeploymentOrUnknown,
		withFinalizer(
			verifyPodExists,
			deleteDeployment,
		),
		verifyNoPod,
	)
}

type deploymentCreator struct {
	access     k8s.Access
	namespace  string
	deployment *appsv1.Deployment
}

func (c *deploymentCreator) Do(_ context.Context) error {
	client := c.access.Kubernetes()
	_, err := client.AppsV1().Deployments(c.namespace).Create(c.deployment)
	return err
}

type deploymentDeleter struct {
	access    k8s.Access
	namespace string
	name      string
}

func (c *deploymentDeleter) Do(_ context.Context) error {
	client := c.access.Kubernetes()
	err := client.AppsV1().Deployments(c.namespace).Delete(c.name, &metav1.DeleteOptions{})
	return err
}

func createDeploymentObject(agentId string) *appsv1.Deployment {
	name := run.StaticIdentifier("upmeter-controller-manager")
	replicas := int32(1)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"heritage":      "upmeter",
				agentLabelKey:   agentId,
				"upmeter-group": "control-plane",
				"upmeter-probe": "controller-manager",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					agentLabelKey:   agentId,
					"upmeter-group": "control-plane",
					"upmeter-probe": "controller-manager",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":           "upmeter",
						"heritage":      "upmeter",
						agentLabelKey:   agentId,
						"upmeter-group": "control-plane",
						"upmeter-probe": "controller-manager",
						"deployment":    name,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "pause",
							Image: "registry.k8s.io/upmeter-nonexistent:3.1415",
							Command: []string{
								"/pause",
							},
						},
					},
					NodeSelector: map[string]string{
						"label-to-avoid":          "scheduling-this-pod-on-any-node",
						"upmeter-only-tests-that": "controller-manager-creates-pods",
					},
					Tolerations: []v1.Toleration{
						{Operator: v1.TolerationOpExists},
					},
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
		},
	}
}
