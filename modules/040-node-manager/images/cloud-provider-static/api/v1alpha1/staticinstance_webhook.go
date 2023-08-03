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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var staticinstancelog = logf.Log.WithName("staticinstance-resource")

func (r *StaticInstance) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-staticinstance,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=staticinstances,verbs=create;update,versions=v1alpha1,name=mstaticinstance.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &StaticInstance{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *StaticInstance) Default() {
	staticinstancelog.Info("default", "name", r.Name)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-staticinstance,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=staticinstances,verbs=create;update,versions=v1alpha1,name=vstaticinstance.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &StaticInstance{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *StaticInstance) ValidateCreate() (admission.Warnings, error) {
	staticinstancelog.Info("validate create", "name", r.Name)

	var errs field.ErrorList

	if r.Spec.CredentialsRef != nil && r.Spec.CredentialsRef.Kind != "StaticInstanceCredentials" {
		errs = append(errs, field.Forbidden(field.NewPath("spec", "credentialsRef", "kind"), "must be a StaticInstanceCredentials"))
	}

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, errs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *StaticInstance) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	staticinstancelog.Info("validate update", "name", r.Name)

	var errs field.ErrorList

	if r.Spec.CredentialsRef != nil && r.Spec.CredentialsRef.Kind != "StaticInstanceCredentials" {
		errs = append(errs, field.Forbidden(field.NewPath("spec", "credentialsRef", "kind"), "must be a StaticInstanceCredentials"))
	}

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, errs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *StaticInstance) ValidateDelete() (admission.Warnings, error) {
	staticinstancelog.Info("validate delete", "name", r.Name)

	return nil, nil
}
