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
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	ModuleResource = "modules"
	ModuleKind     = "Module"

	ModuleSourceEmbedded = "Embedded"

	ModuleConditionEnabled = "Enabled"

	ModulePhaseConflict     = "Conflict"
	ModulePhaseNotInstalled = "NotInstalled"
	ModulePhaseDownloading  = "Downloading"
	ModulePhaseEnqueued     = "Enqueued"
	ModulePhaseReady        = "Ready"
)

var (
	// ModuleGVR GroupVersionResource
	ModuleGVR = schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: ModuleResource,
	}
	// ModuleGVK GroupVersionKind
	ModuleGVK = schema.GroupVersionKind{
		Group:   SchemeGroupVersion.Group,
		Version: SchemeGroupVersion.Version,
		Kind:    ModuleKind,
	}
)

var _ runtime.Object = (*Module)(nil)

// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ModuleList is a list of Module resources
type ModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Module `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Module is a deckhouse module representation.
type Module struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ModuleSpec `json:"spec,omitempty"`

	Properties ModuleProperties `json:"properties,omitempty"`

	Status ModuleStatus `json:"status,omitempty"`
}

type ModuleSpec struct {
	SelectedSource string `json:"selectedSource"`
}

type ModuleProperties struct {
	Weight           uint32   `json:"weight,omitempty"`
	Source           string   `json:"source,omitempty"`
	ReleaseChannel   string   `json:"releaseChannel,omitempty"`
	Stage            string   `json:"stage,omitempty"`
	Description      string   `json:"description,omitempty"`
	AvailableSources []string `json:"availableSources,omitempty"`
}

type ModuleStatus struct {
	Phase      string            `json:"phase"`
	Message    string            `json:"message"`
	HooksState string            `json:"hooksState"`
	Conditions []ModuleCondition `json:"conditions,omitempty"`
}

type ModuleCondition struct {
	Type               string                 `json:"type,omitempty"`
	Reason             string                 `json:"reason,omitempty"`
	Message            string                 `json:"message,omitempty"`
	Status             corev1.ConditionStatus `json:"status,omitempty"`
	LastProbeTime      metav1.Time            `json:"lastProbeTime,omitempty"`
	LastTransitionTime metav1.Time            `json:"lastTransitionTime,omitempty"`
}

func (m *Module) IsEmbedded() bool {
	return m.Spec.SelectedSource == ModuleSourceEmbedded
}

func (m *Module) IsEnabled() bool {
	for _, cond := range m.Status.Conditions {
		if cond.Type == ModuleConditionEnabled {
			if cond.Status == corev1.ConditionTrue {
				return true
			}
			return false
		}
	}
	return false
}

func (m *Module) SetEnabled(enabled bool) {
	for _, cond := range m.Status.Conditions {
		if cond.Type == ModuleConditionEnabled {
			cond.LastProbeTime = metav1.Now()
			if enabled && cond.Status != corev1.ConditionTrue {
				cond.LastTransitionTime = metav1.Now()
				cond.Status = corev1.ConditionTrue
			}
			if !enabled && cond.Status != corev1.ConditionFalse {
				cond.LastTransitionTime = metav1.Now()
				cond.Status = corev1.ConditionFalse
			}
			break
		}
	}
}

func (m *Module) DisabledMoreThan(hours int) bool {
	for _, cond := range m.Status.Conditions {
		if cond.Type == ModuleConditionEnabled && cond.Status == corev1.ConditionFalse {
			return time.Now().Sub(cond.LastTransitionTime.Time).Hours() >= float64(hours)
		}
	}
	return false
}
