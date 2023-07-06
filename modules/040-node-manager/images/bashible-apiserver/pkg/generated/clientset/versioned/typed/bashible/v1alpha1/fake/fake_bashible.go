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
	v1alpha1 "bashible-apiserver/pkg/apis/bashible/v1alpha1"
	bashiblev1alpha1 "bashible-apiserver/pkg/generated/applyconfiguration/bashible/v1alpha1"
	"context"
	json "encoding/json"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBashibles implements BashibleInterface
type FakeBashibles struct {
	Fake *FakeBashibleV1alpha1
}

var bashiblesResource = v1alpha1.SchemeGroupVersion.WithResource("bashibles")

var bashiblesKind = v1alpha1.SchemeGroupVersion.WithKind("Bashible")

// Get takes name of the bashible, and returns the corresponding bashible object, and an error if there is any.
func (c *FakeBashibles) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Bashible, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(bashiblesResource, name), &v1alpha1.Bashible{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bashible), err
}

// List takes label and field selectors, and returns the list of Bashibles that match those selectors.
func (c *FakeBashibles) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BashibleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(bashiblesResource, bashiblesKind, opts), &v1alpha1.BashibleList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.BashibleList{ListMeta: obj.(*v1alpha1.BashibleList).ListMeta}
	for _, item := range obj.(*v1alpha1.BashibleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested bashibles.
func (c *FakeBashibles) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(bashiblesResource, opts))
}

// Create takes the representation of a bashible and creates it.  Returns the server's representation of the bashible, and an error, if there is any.
func (c *FakeBashibles) Create(ctx context.Context, bashible *v1alpha1.Bashible, opts v1.CreateOptions) (result *v1alpha1.Bashible, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(bashiblesResource, bashible), &v1alpha1.Bashible{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bashible), err
}

// Update takes the representation of a bashible and updates it. Returns the server's representation of the bashible, and an error, if there is any.
func (c *FakeBashibles) Update(ctx context.Context, bashible *v1alpha1.Bashible, opts v1.UpdateOptions) (result *v1alpha1.Bashible, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(bashiblesResource, bashible), &v1alpha1.Bashible{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bashible), err
}

// Delete takes name of the bashible and deletes it. Returns an error if one occurs.
func (c *FakeBashibles) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(bashiblesResource, name, opts), &v1alpha1.Bashible{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBashibles) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(bashiblesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.BashibleList{})
	return err
}

// Patch applies the patch and returns the patched bashible.
func (c *FakeBashibles) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Bashible, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(bashiblesResource, name, pt, data, subresources...), &v1alpha1.Bashible{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bashible), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied bashible.
func (c *FakeBashibles) Apply(ctx context.Context, bashible *bashiblev1alpha1.BashibleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Bashible, err error) {
	if bashible == nil {
		return nil, fmt.Errorf("bashible provided to Apply must not be nil")
	}
	data, err := json.Marshal(bashible)
	if err != nil {
		return nil, err
	}
	name := bashible.Name
	if name == nil {
		return nil, fmt.Errorf("bashible.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(bashiblesResource, *name, types.ApplyPatchType, data), &v1alpha1.Bashible{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bashible), err
}
