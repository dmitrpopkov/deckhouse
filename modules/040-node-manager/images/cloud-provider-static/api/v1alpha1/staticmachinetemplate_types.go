/*
Copyright 2023.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StaticMachineTemplateSpec defines the desired state of StaticMachineTemplate
type StaticMachineTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Template StaticMachineTemplateSpecTemplate `json:"template"`
}

type StaticMachineTemplateSpecTemplate struct {
	Spec StaticMachineTemplateSpecTemplateSpec `json:"spec"`
}

type StaticMachineTemplateSpecTemplateSpec struct {
	LabelSelector metav1.LabelSelector `json:"labelSelector"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`
}

//+kubebuilder:object:root=true

// StaticMachineTemplate is the Schema for the staticmachinetemplates API
type StaticMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec StaticMachineTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// StaticMachineTemplateList contains a list of StaticMachineTemplate
type StaticMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StaticMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StaticMachineTemplate{}, &StaticMachineTemplateList{})
}
