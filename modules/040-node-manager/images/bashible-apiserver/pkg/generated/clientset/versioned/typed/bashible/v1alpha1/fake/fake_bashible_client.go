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
	v1alpha1 "bashible-apiserver/pkg/generated/clientset/versioned/typed/bashible/v1alpha1"

	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeBashibleV1alpha1 struct {
	*testing.Fake
}

func (c *FakeBashibleV1alpha1) Bashibles() v1alpha1.BashibleInterface {
	return &FakeBashibles{c}
}

func (c *FakeBashibleV1alpha1) Bootstraps() v1alpha1.BootstrapInterface {
	return &FakeBootstraps{c}
}

func (c *FakeBashibleV1alpha1) NodeGroupBundles() v1alpha1.NodeGroupBundleInterface {
	return &FakeNodeGroupBundles{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeBashibleV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
