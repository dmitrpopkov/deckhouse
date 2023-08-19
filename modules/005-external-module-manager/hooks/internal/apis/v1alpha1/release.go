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
	"encoding/json"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

const (
	PhasePending    = "Pending"
	PhaseDeployed   = "Deployed"
	PhaseSuperseded = "Superseded"
	PhaseSuspended  = "Suspended"
)

// ModuleRelease is a Module release object.
type ModuleRelease struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ModuleReleaseSpec `json:"spec"`

	Status ModuleReleaseStatus `json:"status,omitempty"`
}

type ModuleReleaseSpec struct {
	ModuleName string          `json:"moduleName"`
	Version    *semver.Version `json:"version,omitempty"`
	Weight     int             `json:"weight,omitempty"`

	ApplyAfter   *time.Time        `json:"applyAfter,omitempty"`
	Requirements map[string]string `json:"requirements,omitempty"`
}

type ModuleReleaseStatus struct {
	Phase          string    `json:"phase,omitempty"`
	Approved       bool      `json:"approved"`
	TransitionTime time.Time `json:"transitionTime,omitempty"`
	Message        string    `json:"message"`
}

type ModuleReleaseKind struct{}

func (in *ModuleReleaseStatus) GetObjectKind() schema.ObjectKind {
	return &ModuleReleaseKind{}
}

func (f *ModuleReleaseKind) SetGroupVersionKind(_ schema.GroupVersionKind) {}
func (f *ModuleReleaseKind) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{Group: "deckhouse.io", Version: "v1alpha1", Kind: "ModuleRelease"}
}

// Duration custom type for appropriate json marshalling / unmarshalling (like "15m")
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
