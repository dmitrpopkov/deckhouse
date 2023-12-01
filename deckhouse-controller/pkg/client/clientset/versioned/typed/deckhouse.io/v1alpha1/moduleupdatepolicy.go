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
	"context"
	"time"

	v1alpha1 "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	scheme "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ModuleUpdatePoliciesGetter has a method to return a ModuleUpdatePolicyInterface.
// A group's client should implement this interface.
type ModuleUpdatePoliciesGetter interface {
	ModuleUpdatePolicies() ModuleUpdatePolicyInterface
}

// ModuleUpdatePolicyInterface has methods to work with ModuleUpdatePolicy resources.
type ModuleUpdatePolicyInterface interface {
	Create(ctx context.Context, moduleUpdatePolicy *v1alpha1.ModuleUpdatePolicy, opts v1.CreateOptions) (*v1alpha1.ModuleUpdatePolicy, error)
	Update(ctx context.Context, moduleUpdatePolicy *v1alpha1.ModuleUpdatePolicy, opts v1.UpdateOptions) (*v1alpha1.ModuleUpdatePolicy, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ModuleUpdatePolicy, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ModuleUpdatePolicyList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ModuleUpdatePolicy, err error)
	ModuleUpdatePolicyExpansion
}

// moduleUpdatePolicies implements ModuleUpdatePolicyInterface
type moduleUpdatePolicies struct {
	client rest.Interface
}

// newModuleUpdatePolicies returns a ModuleUpdatePolicies
func newModuleUpdatePolicies(c *DeckhouseV1alpha1Client) *moduleUpdatePolicies {
	return &moduleUpdatePolicies{
		client: c.RESTClient(),
	}
}

// Get takes name of the moduleUpdatePolicy, and returns the corresponding moduleUpdatePolicy object, and an error if there is any.
func (c *moduleUpdatePolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ModuleUpdatePolicy, err error) {
	result = &v1alpha1.ModuleUpdatePolicy{}
	err = c.client.Get().
		Resource("moduleupdatepolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ModuleUpdatePolicies that match those selectors.
func (c *moduleUpdatePolicies) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ModuleUpdatePolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ModuleUpdatePolicyList{}
	err = c.client.Get().
		Resource("moduleupdatepolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested moduleUpdatePolicies.
func (c *moduleUpdatePolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("moduleupdatepolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a moduleUpdatePolicy and creates it.  Returns the server's representation of the moduleUpdatePolicy, and an error, if there is any.
func (c *moduleUpdatePolicies) Create(ctx context.Context, moduleUpdatePolicy *v1alpha1.ModuleUpdatePolicy, opts v1.CreateOptions) (result *v1alpha1.ModuleUpdatePolicy, err error) {
	result = &v1alpha1.ModuleUpdatePolicy{}
	err = c.client.Post().
		Resource("moduleupdatepolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(moduleUpdatePolicy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a moduleUpdatePolicy and updates it. Returns the server's representation of the moduleUpdatePolicy, and an error, if there is any.
func (c *moduleUpdatePolicies) Update(ctx context.Context, moduleUpdatePolicy *v1alpha1.ModuleUpdatePolicy, opts v1.UpdateOptions) (result *v1alpha1.ModuleUpdatePolicy, err error) {
	result = &v1alpha1.ModuleUpdatePolicy{}
	err = c.client.Put().
		Resource("moduleupdatepolicies").
		Name(moduleUpdatePolicy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(moduleUpdatePolicy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the moduleUpdatePolicy and deletes it. Returns an error if one occurs.
func (c *moduleUpdatePolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("moduleupdatepolicies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *moduleUpdatePolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("moduleupdatepolicies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched moduleUpdatePolicy.
func (c *moduleUpdatePolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ModuleUpdatePolicy, err error) {
	result = &v1alpha1.ModuleUpdatePolicy{}
	err = c.client.Patch(pt).
		Resource("moduleupdatepolicies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
