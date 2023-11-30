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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeModuleReleases implements ModuleReleaseInterface
type FakeModuleReleases struct {
	Fake *FakeDeckhouseV1alpha1
}

var modulereleasesResource = schema.GroupVersionResource{Group: "deckhouse.io", Version: "v1alpha1", Resource: "modulereleases"}

var modulereleasesKind = schema.GroupVersionKind{Group: "deckhouse.io", Version: "v1alpha1", Kind: "ModuleRelease"}

// Get takes name of the moduleRelease, and returns the corresponding moduleRelease object, and an error if there is any.
func (c *FakeModuleReleases) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ModuleRelease, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(modulereleasesResource, name), &v1alpha1.ModuleRelease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ModuleRelease), err
}

// List takes label and field selectors, and returns the list of ModuleReleases that match those selectors.
func (c *FakeModuleReleases) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ModuleReleaseList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(modulereleasesResource, modulereleasesKind, opts), &v1alpha1.ModuleReleaseList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ModuleReleaseList{ListMeta: obj.(*v1alpha1.ModuleReleaseList).ListMeta}
	for _, item := range obj.(*v1alpha1.ModuleReleaseList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested moduleReleases.
func (c *FakeModuleReleases) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(modulereleasesResource, opts))
}

// Create takes the representation of a moduleRelease and creates it.  Returns the server's representation of the moduleRelease, and an error, if there is any.
func (c *FakeModuleReleases) Create(ctx context.Context, moduleRelease *v1alpha1.ModuleRelease, opts v1.CreateOptions) (result *v1alpha1.ModuleRelease, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(modulereleasesResource, moduleRelease), &v1alpha1.ModuleRelease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ModuleRelease), err
}

// Update takes the representation of a moduleRelease and updates it. Returns the server's representation of the moduleRelease, and an error, if there is any.
func (c *FakeModuleReleases) Update(ctx context.Context, moduleRelease *v1alpha1.ModuleRelease, opts v1.UpdateOptions) (result *v1alpha1.ModuleRelease, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(modulereleasesResource, moduleRelease), &v1alpha1.ModuleRelease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ModuleRelease), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeModuleReleases) UpdateStatus(ctx context.Context, moduleRelease *v1alpha1.ModuleRelease, opts v1.UpdateOptions) (*v1alpha1.ModuleRelease, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(modulereleasesResource, "status", moduleRelease), &v1alpha1.ModuleRelease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ModuleRelease), err
}

// Delete takes name of the moduleRelease and deletes it. Returns an error if one occurs.
func (c *FakeModuleReleases) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(modulereleasesResource, name, opts), &v1alpha1.ModuleRelease{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeModuleReleases) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(modulereleasesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ModuleReleaseList{})
	return err
}

// Patch applies the patch and returns the patched moduleRelease.
func (c *FakeModuleReleases) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ModuleRelease, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(modulereleasesResource, name, pt, data, subresources...), &v1alpha1.ModuleRelease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ModuleRelease), err
}
