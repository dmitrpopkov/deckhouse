/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v1alpha1 "bashible-apiserver/pkg/apis/bashible/v1alpha1"
	bashiblev1alpha1 "bashible-apiserver/pkg/generated/applyconfiguration/bashible/v1alpha1"

	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=bashible.deckhouse.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("Bashible"):
		return &bashiblev1alpha1.BashibleApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Bootstrap"):
		return &bashiblev1alpha1.BootstrapApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NodeGroupBundle"):
		return &bashiblev1alpha1.NodeGroupBundleApplyConfiguration{}

	}
	return nil
}
