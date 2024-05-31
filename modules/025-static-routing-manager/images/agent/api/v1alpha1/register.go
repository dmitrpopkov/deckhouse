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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	APIGroup         = "network.deckhouse.io"
	InternalAPIGroup = "internal.network.deckhouse.io"
	APIVersion       = "v1alpha1"
	RTKind           = "RoutingTable"
	NRTKind          = "SDNInternalNodeRoutingTable"
	IRSKind          = "IPRuleSet"
	NIRSKind         = "SDNInternalNodeIPRuleSet"
)

// SchemeGroupVersion is group version used to register these objects
var (
	SchemeGroupVersion = schema.GroupVersion{
		Group:   APIGroup,
		Version: APIVersion,
	}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme

	InternalSchemeGroupVersion = schema.GroupVersion{
		Group:   InternalAPIGroup,
		Version: APIVersion,
	}
	InternalSchemeBuilder = runtime.NewSchemeBuilder(addKnownInternalTypes)
	AddInternalToScheme   = InternalSchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&RoutingTable{},
		&RoutingTableList{},
		&IPRuleSet{},
		&IPRuleSetList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// Adds the list of known Internal types to Scheme.
func addKnownInternalTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(InternalSchemeGroupVersion,
		&SDNInternalNodeRoutingTable{},
		&SDNInternalNodeRoutingTableList{},
		&SDNInternalNodeIPRuleSet{},
		&SDNInternalNodeIPRuleSetList{},
	)
	metav1.AddToGroupVersion(scheme, InternalSchemeGroupVersion)
	return nil
}
