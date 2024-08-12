/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ProjectStateValidationError = "ValidationError"
	ProjectStateUpgradeError    = "UpgradeError"

	ProjectStateTemplateNotFound     = "TemplateNotFound"
	ProjectStateGettingTemplateError = "GettingTemplateError"

	ProjectStateDeploying = "Deploying"
	ProjectStateDeployed  = "Deployed"
)

const (
	ProjectKind     = "Project"
	ProjectResource = "projects"
)

type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}

func (p *ProjectList) DeepCopyObject() runtime.Object {
	return p.DeepCopy()
}
func (p *ProjectList) DeepCopy() *ProjectList {
	if p == nil {
		return nil
	}
	newObj := new(ProjectList)
	p.DeepCopyInto(newObj)
	return newObj
}
func (p *ProjectList) DeepCopyInto(newObj *ProjectList) {
	*newObj = *p
	newObj.TypeMeta = p.TypeMeta
	p.ListMeta.DeepCopyInto(&newObj.ListMeta)
	if p.Items != nil {
		in, out := &p.Items, &newObj.Items
		*out = make([]Project, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

func (p *Project) DeepCopyObject() runtime.Object {
	return p.DeepCopy()
}
func (p *Project) DeepCopy() *Project {
	if p == nil {
		return nil
	}
	newObj := Project{}
	p.DeepCopyInto(&newObj)
	return &newObj
}
func (p *Project) DeepCopyInto(newObj *Project) {
	*newObj = *p
	newObj.TypeMeta = p.TypeMeta
	p.ObjectMeta.DeepCopyInto(&newObj.ObjectMeta)
	p.Spec.DeepCopyInto(&newObj.Spec)
	p.Status.DeepCopyInto(&newObj.Status)
}

type ProjectSpec struct {
	// Description of the Project
	Description string `json:"description,omitempty"`

	// Name of ProjectTemplate to use to create Project
	ProjectTemplateName string `json:"projectTemplateName,omitempty"`

	// Values for resource templates from ProjectTemplate
	// in helm values format that map to the open-api specification
	// from the ValuesSchema ProjectTemplate field
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

func (p *ProjectSpec) DeepCopy() *ProjectSpec {
	if p == nil {
		return nil
	}
	newObj := new(ProjectSpec)
	p.DeepCopyInto(newObj)
	return newObj
}
func (p *ProjectSpec) DeepCopyInto(newObj *ProjectSpec) {
	*newObj = *p
	newObj.Description = p.Description
	newObj.ProjectTemplateName = p.ProjectTemplateName
	newObj.Parameters = make(map[string]interface{})
	for key, value := range p.Parameters {
		newObj.Parameters[key] = value
	}
}

type ProjectStatus struct {
	// Status message.
	Message string `json:"message,omitempty"`

	// Current state.
	State string `json:"state,omitempty"`

	// Involved namespaces
	Namespaces []string `json:"namespaces,omitempty"`

	Sync bool `json:"sync,omitempty"`
}

func (p *ProjectStatus) DeepCopy() *ProjectStatus {
	if p == nil {
		return nil
	}
	newObj := new(ProjectStatus)
	p.DeepCopyInto(newObj)
	return newObj
}
func (p *ProjectStatus) DeepCopyInto(newObj *ProjectStatus) {
	*newObj = *p
	newObj.Message = p.Message
	newObj.State = p.State
	newObj.Sync = p.Sync
	if p.Namespaces != nil {
		in, out := &p.Namespaces, &newObj.Namespaces
		*in = make([]string, len(*out))
		copy(*in, *out)
	}
}
