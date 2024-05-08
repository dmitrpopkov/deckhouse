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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1alpha1 "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=deckhouse.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("deckhousereleases"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Deckhouse().V1alpha1().DeckhouseReleases().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("modules"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Deckhouse().V1alpha1().Modules().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("moduleconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Deckhouse().V1alpha1().ModuleConfigs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("moduledocumentations"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Deckhouse().V1alpha1().ModuleDocumentations().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("modulepulloverrides"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Deckhouse().V1alpha1().ModulePullOverrides().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("modulereleases"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Deckhouse().V1alpha1().ModuleReleases().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("modulesources"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Deckhouse().V1alpha1().ModuleSources().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("moduleupdatepolicies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Deckhouse().V1alpha1().ModuleUpdatePolicies().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
