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

package v1alpha1

import (
	"context"
	time "time"

	deckhouseiov1alpha1 "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	versioned "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/client/clientset/versioned"
	internalinterfaces "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/client/listers/deckhouse.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ModuleConfigInformer provides access to a shared informer and lister for
// ModuleConfigs.
type ModuleConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ModuleConfigLister
}

type moduleConfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewModuleConfigInformer constructs a new informer for ModuleConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewModuleConfigInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredModuleConfigInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredModuleConfigInformer constructs a new informer for ModuleConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredModuleConfigInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DeckhouseV1alpha1().ModuleConfigs().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DeckhouseV1alpha1().ModuleConfigs().Watch(context.TODO(), options)
			},
		},
		&deckhouseiov1alpha1.ModuleConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *moduleConfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredModuleConfigInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *moduleConfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&deckhouseiov1alpha1.ModuleConfig{}, f.defaultInformer)
}

func (f *moduleConfigInformer) Lister() v1alpha1.ModuleConfigLister {
	return v1alpha1.NewModuleConfigLister(f.Informer().GetIndexer())
}
