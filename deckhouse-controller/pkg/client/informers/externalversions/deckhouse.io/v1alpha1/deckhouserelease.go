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

// DeckhouseReleaseInformer provides access to a shared informer and lister for
// DeckhouseReleases.
type DeckhouseReleaseInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.DeckhouseReleaseLister
}

type deckhouseReleaseInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewDeckhouseReleaseInformer constructs a new informer for DeckhouseRelease type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDeckhouseReleaseInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDeckhouseReleaseInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredDeckhouseReleaseInformer constructs a new informer for DeckhouseRelease type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDeckhouseReleaseInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DeckhouseV1alpha1().DeckhouseReleases().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.DeckhouseV1alpha1().DeckhouseReleases().Watch(context.TODO(), options)
			},
		},
		&deckhouseiov1alpha1.DeckhouseRelease{},
		resyncPeriod,
		indexers,
	)
}

func (f *deckhouseReleaseInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDeckhouseReleaseInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *deckhouseReleaseInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&deckhouseiov1alpha1.DeckhouseRelease{}, f.defaultInformer)
}

func (f *deckhouseReleaseInformer) Lister() v1alpha1.DeckhouseReleaseLister {
	return v1alpha1.NewDeckhouseReleaseLister(f.Informer().GetIndexer())
}
