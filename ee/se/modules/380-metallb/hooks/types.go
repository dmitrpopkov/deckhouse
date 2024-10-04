/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	keyAnnotationL2BalancerName   = "network.deckhouse.io/l2-load-balancer-name"
	keyAnnotationExternalIPsCount = "network.deckhouse.io/l2-load-balancer-external-ips-count"
	memberLabelKey                = "l2-load-balancer.network.deckhouse.io/member"
	metallbAllocatedPool          = "metallb.universe.tf/ip-allocated-from-pool"
	l2LoadBalancerIPsAnnotate     = "network.deckhouse.io/l2-load-balancer-ips"
)

type NodeInfo struct {
	Name      string
	Labels    map[string]string
	IsLabeled bool
}

type ServiceInfo struct {
	PublishNotReadyAddresses  bool
	Name                      string
	Namespace                 string
	LoadBalancerClass         string
	AssignedLoadBalancerClass string
	ClusterIP                 string
	ExternalIPsCount          int
	Ports                     []v1.ServicePort
	ExternalTrafficPolicy     v1.ServiceExternalTrafficPolicy
	InternalTrafficPolicy     v1.ServiceInternalTrafficPolicy
	Selector                  map[string]string
	DesiredIPs                []string
}

type L2LBServiceStatusInfo struct {
	Name              string
	Namespace         string
	LoadBalancerClass string
	IP                string
}

type L2LBServiceConfig struct {
	PublishNotReadyAddresses   bool                            `json:"publishNotReadyAddresses"`
	Name                       string                          `json:"name"`
	Namespace                  string                          `json:"namespace"`
	ServiceName                string                          `json:"serviceName"`
	ServiceNamespace           string                          `json:"serviceNamespace"`
	PreferredNode              string                          `json:"preferredNode,omitempty"`
	ClusterIP                  string                          `json:"clusterIP"`
	Ports                      []v1.ServicePort                `json:"ports"`
	ExternalTrafficPolicy      v1.ServiceExternalTrafficPolicy `json:"externalTrafficPolicy"`
	InternalTrafficPolicy      v1.ServiceInternalTrafficPolicy `json:"internalTrafficPolicy"`
	Selector                   map[string]string               `json:"selector"`
	MetalLoadBalancerClassName string                          `json:"loadBalancerName"`
	DesiredIP                  string                          `json:"desiredIP"`
}

type MetalLoadBalancerClassInfo struct {
	Name         string            `json:"name"`
	AddressPool  []string          `json:"addressPool"`
	Interfaces   []string          `json:"interfaces"`
	NodeSelector map[string]string `json:"nodeSelector"`
	IsDefault    bool              `json:"isDefault"`
}

type MetalLoadBalancerClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetalLoadBalancerClassSpec   `json:"spec,omitempty"`
	Status MetalLoadBalancerClassStatus `json:"status,omitempty"`
}

type MetalLoadBalancerClassSpec struct {
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	L2           L2Type            `json:"l2,omitempty"`
	AddressPool  []string          `json:"addressPool,omitempty"`
	IsDefault    bool              `json:"isDefault,omitempty"`
}

type L2Type struct {
	Interfaces []string `json:"interfaces,omitempty"`
}

type MetalLoadBalancerClassStatus struct {
}

type SDNInternalL2LBServiceSpec struct {
	v1.ServiceSpec `json:",inline"`
	ServiceRef     SDNInternalL2LBServiceReference `json:"serviceRef"`
}

type SDNInternalL2LBServiceReference struct {
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

type SDNInternalL2LBService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   SDNInternalL2LBServiceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status v1.ServiceStatus           `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type SDNIpsLBName struct {
	Ips         []string
	LBClassName string
}
