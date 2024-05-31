/*
Copyright 2024 Flant JSC

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// network.deckhouse.io/v1alpha1

type Route struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
}

type Routes struct {
	Routes []Route `json:"routes"`
}

// CR RoutingTable

type RoutingTable struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RoutingTableSpec   `json:"spec"`
	Status            RoutingTableStatus `json:"status,omitempty"`
}

type RoutingTableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RoutingTable `json:"items"`
}

type RoutingTableSpec struct {
	IPRoutingTableID int               `json:"ipRoutingTableID"`
	Routes           []Route           `json:"routes"`
	NodeSelector     map[string]string `json:"nodeSelector"`
}

type RoutingTableStatus struct {
	ObservedGeneration        int64               `json:"observedGeneration,omitempty"`
	IPRoutingTableID          int                 `json:"ipRoutingTableID,omitempty"`
	ReadyNodeRoutingTables    int                 `json:"readyNodeRoutingTables,omitempty"`
	AffectedNodeRoutingTables int                 `json:"affectedNodeRoutingTables,omitempty"`
	Conditions                []ExtendedCondition `json:"conditions,omitempty"`
}

// CR SDNInternalNodeRoutingTable

type SDNInternalNodeRoutingTable struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SDNInternalNodeRoutingTableSpec   `json:"spec"`
	Status            SDNInternalNodeRoutingTableStatus `json:"status,omitempty"`
}

type SDNInternalNodeRoutingTableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SDNInternalNodeRoutingTable `json:"items"`
}

type SDNInternalNodeRoutingTableSpec struct {
	NodeName         string  `json:"nodeName"`
	IPRoutingTableID int     `json:"ipRoutingTableID"`
	Routes           []Route `json:"routes"`
}

type SDNInternalNodeRoutingTableStatus struct {
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	AppliedRoutes      []Route             `json:"appliedRoutes,omitempty"`
	Conditions         []ExtendedCondition `json:"conditions,omitempty"`
}
