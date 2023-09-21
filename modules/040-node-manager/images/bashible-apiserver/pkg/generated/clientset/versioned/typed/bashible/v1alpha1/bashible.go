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

package v1alpha1

import (
	v1alpha1 "d8.io/bashible/apis/bashible/v1alpha1"
	bashiblev1alpha1 "d8.io/bashible/generated/applyconfiguration/bashible/v1alpha1"
	scheme "d8.io/bashible/generated/clientset/versioned/scheme"
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BashiblesGetter has a method to return a BashibleInterface.
// A group's client should implement this interface.
type BashiblesGetter interface {
	Bashibles() BashibleInterface
}

// BashibleInterface has methods to work with Bashible resources.
type BashibleInterface interface {
	Create(ctx context.Context, bashible *v1alpha1.Bashible, opts v1.CreateOptions) (*v1alpha1.Bashible, error)
	Update(ctx context.Context, bashible *v1alpha1.Bashible, opts v1.UpdateOptions) (*v1alpha1.Bashible, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Bashible, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.BashibleList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Bashible, err error)
	Apply(ctx context.Context, bashible *bashiblev1alpha1.BashibleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Bashible, err error)
	BashibleExpansion
}

// bashibles implements BashibleInterface
type bashibles struct {
	client rest.Interface
}

// newBashibles returns a Bashibles
func newBashibles(c *BashibleV1alpha1Client) *bashibles {
	return &bashibles{
		client: c.RESTClient(),
	}
}

// Get takes name of the bashible, and returns the corresponding bashible object, and an error if there is any.
func (c *bashibles) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Bashible, err error) {
	result = &v1alpha1.Bashible{}
	err = c.client.Get().
		Resource("bashibles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Bashibles that match those selectors.
func (c *bashibles) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BashibleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.BashibleList{}
	err = c.client.Get().
		Resource("bashibles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested bashibles.
func (c *bashibles) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("bashibles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a bashible and creates it.  Returns the server's representation of the bashible, and an error, if there is any.
func (c *bashibles) Create(ctx context.Context, bashible *v1alpha1.Bashible, opts v1.CreateOptions) (result *v1alpha1.Bashible, err error) {
	result = &v1alpha1.Bashible{}
	err = c.client.Post().
		Resource("bashibles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(bashible).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a bashible and updates it. Returns the server's representation of the bashible, and an error, if there is any.
func (c *bashibles) Update(ctx context.Context, bashible *v1alpha1.Bashible, opts v1.UpdateOptions) (result *v1alpha1.Bashible, err error) {
	result = &v1alpha1.Bashible{}
	err = c.client.Put().
		Resource("bashibles").
		Name(bashible.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(bashible).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the bashible and deletes it. Returns an error if one occurs.
func (c *bashibles) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("bashibles").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *bashibles) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("bashibles").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched bashible.
func (c *bashibles) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Bashible, err error) {
	result = &v1alpha1.Bashible{}
	err = c.client.Patch(pt).
		Resource("bashibles").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied bashible.
func (c *bashibles) Apply(ctx context.Context, bashible *bashiblev1alpha1.BashibleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Bashible, err error) {
	if bashible == nil {
		return nil, fmt.Errorf("bashible provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(bashible)
	if err != nil {
		return nil, err
	}
	name := bashible.Name
	if name == nil {
		return nil, fmt.Errorf("bashible.Name must be provided to Apply")
	}
	result = &v1alpha1.Bashible{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("bashibles").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
