/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package k8shandler

// import (
// 	"context"
// 	"fmt"
// 	pkg_cfg "system-registry-manager/pkg/cfg"
// 	"time"

// 	"github.com/sirupsen/logrus"
// 	appsv1 "k8s.io/api/apps/v1"
// 	corev1 "k8s.io/api/core/v1"
// 	pkg_logs "system-registry-manager/pkg/logs"
// )

// type CommonHandler struct {
// 	ctx                     context.Context
// 	log                     *logrus.Entry
// 	MasterNodesResource     *MasterNodesResource
// 	SeaweedfsPodsResource   *SeaweedfsPodsResource
// 	WorkerDaemonsetResource *WorkerDaemonsetResource
// 	WorkerEndpointResource  *WorkerEndpointResource

// 	KubernetesResourcesHandler *KubernetesResourcesHandler
// }

// func NewCommonHandler(ctx context.Context) (*CommonHandler, error) {
// 	log := pkg_logs.GetLoggerFromContext(ctx)

// 	cfg := pkg_cfg.GetConfig()

// 	commonHandler := CommonHandler{ctx: ctx, log: log}

// 	commonHandler.KubernetesResourcesHandler = NewKubernetesResourcesHandler(
// 		ctx,
// 		cfg.K8sClient,
// 		time.Duration(0)*time.Minute,
// 		cfg.Manager.Namespace,
// 	)

// 	var err error
// 	commonHandler.MasterNodesResource, err = NewMasterNodesResource(log)
// 	if err != nil {
// 		return nil, err
// 	}
// 	commonHandler.SeaweedfsPodsResource, err = NewSeaweedfsPodsResource(log, pkg_cfg.SeaweedfsStaticPodLabelsSelector)
// 	if err != nil {
// 		return nil, err
// 	}
// 	commonHandler.WorkerDaemonsetResource = NewWorkerDaemonsetResource(log, cfg.Manager.DaemonsetName)
// 	commonHandler.WorkerEndpointResource = NewWorkerEndpointResource(log, cfg.Manager.ServiceName)

// 	commonHandler.KubernetesResourcesHandler.Subscribe(commonHandler.MasterNodesResource)
// 	commonHandler.KubernetesResourcesHandler.Subscribe(commonHandler.SeaweedfsPodsResource)
// 	commonHandler.KubernetesResourcesHandler.Subscribe(commonHandler.WorkerDaemonsetResource)
// 	commonHandler.KubernetesResourcesHandler.Subscribe(commonHandler.WorkerEndpointResource)

// 	return &commonHandler, nil
// }

// func (c *CommonHandler) Start() {
// 	c.KubernetesResourcesHandler.Start()
// }

// func (c *CommonHandler) Stop() {
// 	c.KubernetesResourcesHandler.Stop()
// }

// func (c *CommonHandler) GetMasterNodeNameList() []string {
// 	data := c.MasterNodesResource.GetData()
// 	if data == nil {
// 		return nil
// 	}

// 	nodes := make([]string, 0, len(data))
// 	for nodeName := range data {
// 		nodes = append(nodes, nodeName)
// 	}
// 	return nodes
// }

// func (c *CommonHandler) GetWorkerEndpointByNodeName(nodeNmae string) *corev1.EndpointAddress {
// 	endpointsData := c.WorkerEndpointResource.GetData()

// 	if endpointsData == nil {
// 		return nil
// 	}
// 	for _, subsets := range endpointsData.Subsets {
// 		for _, address := range subsets.Addresses {
// 			if address.NodeName != nil && *address.NodeName == nodeNmae {
// 				return &address
// 			}
// 		}
// 	}
// 	return nil
// }

// func (c *CommonHandler) GetSeaweedfsPodByNodeName(nodeNmae string) *corev1.Pod {
// 	seaweedfsPods := c.SeaweedfsPodsResource.GetData()
// 	if seaweedfsPods == nil {
// 		return nil
// 	}
// 	for _, pod := range seaweedfsPods {
// 		if pod.Spec.NodeName == nodeNmae {
// 			return &pod
// 		}
// 	}
// 	return nil
// }

// func (c *CommonHandler) GetMasterNodeByNodeName(nodeNmae string) *corev1.Node {
// 	masterNodedata := c.MasterNodesResource.GetData()
// 	if masterNodedata == nil {
// 		return nil
// 	}
// 	if masterNode, ok := masterNodedata[nodeNmae]; ok {
// 		return &masterNode
// 	}
// 	return nil
// }

// func (c *CommonHandler) WaitWorkerDaemonset() (*appsv1.DaemonSet, error) {
// 	for i := 0; i < pkg_cfg.MaxRetries; i++ {
// 		defer time.Sleep(1 * time.Second)

// 		dsInfo := c.WorkerDaemonsetResource.GetData()
// 		if dsInfo == nil {
// 			continue
// 		}
// 		if dsInfo.Status.DesiredNumberScheduled == dsInfo.Status.NumberAvailable {
// 			return dsInfo, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("error WaitDaemonsetPods")
// }

// func (c *CommonHandler) WaitWorkerEndpoints() (*corev1.Endpoints, error) {
// 	for i := 0; i < pkg_cfg.MaxRetries; i++ {
// 		defer time.Sleep(1 * time.Second)

// 		ep := c.WorkerEndpointResource.GetData()

// 		if ep == nil {
// 			continue
// 		}

// 		for _, subset := range ep.Subsets {
// 			if len(subset.NotReadyAddresses) == 0 {
// 				return ep, nil
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("error WaitWorkerEndpoints")
// }

// func (c *CommonHandler) WaitAllWorkers() (int, error) {
// 	dsInfo, err := c.WaitWorkerDaemonset()
// 	if err != nil {
// 		return 0, err
// 	}

// 	ep, err := c.WaitWorkerEndpoints()
// 	if err != nil {
// 		return 0, err
// 	}

// 	numberOfNode := dsInfo.Status.DesiredNumberScheduled

// 	if len(ep.Subsets) == 0 {
// 		return 0, fmt.Errorf("error len(ep.Subsets) == 0")
// 	}
// 	if len(ep.Subsets[0].Addresses) != int(numberOfNode) {
// 		return 0, fmt.Errorf("error len(ep.Subsets[0].Addresses) != numberOfNode")
// 	}

// 	return int(numberOfNode), nil
// }

// func (c *CommonHandler) GetAllDataByNodeName(nodeNmae string) *MergeInfo {
// 	masterNode := c.GetMasterNodeByNodeName(nodeNmae)
// 	if masterNode == nil {
// 		return nil
// 	}
// 	workerEndpoint := c.GetWorkerEndpointByNodeName(nodeNmae)
// 	if workerEndpoint == nil {
// 		return nil
// 	}
// 	seaweedfsPod := c.GetSeaweedfsPodByNodeName(nodeNmae)

// 	return &MergeInfo{
// 		MasterNode:   *masterNode,
// 		Executor:       *workerEndpoint,
// 		SeaweedfsPod: seaweedfsPod,
// 	}
// }

// type MergeInfo struct {
// 	MasterNode   corev1.Node
// 	Executor       corev1.EndpointAddress
// 	SeaweedfsPod *corev1.Pod
// }
