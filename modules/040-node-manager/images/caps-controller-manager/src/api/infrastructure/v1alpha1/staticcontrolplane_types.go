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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StaticControlPlaneSpec defines the desired state of StaticControlPlane
type StaticControlPlaneSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// StaticControlPlaneStatus defines the observed state of StaticControlPlane
type StaticControlPlaneStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	Ready bool `json:"ready,omitempty"`

	// Initialized is true when the control plane is available for initial contact.
	// This may occur before the control plane is fully ready.
	// +optional
	Initialized bool `json:"initialized,omitempty"`

	// +optional
	ExternalManagedControlPlane bool `json:"externalManagedControlPlane,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StaticControlPlane is the Schema for the staticcontrolplanes API
type StaticControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StaticControlPlaneSpec   `json:"spec,omitempty"`
	Status StaticControlPlaneStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StaticControlPlaneList contains a list of StaticControlPlane
type StaticControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StaticControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StaticControlPlane{}, &StaticControlPlaneList{})
}
