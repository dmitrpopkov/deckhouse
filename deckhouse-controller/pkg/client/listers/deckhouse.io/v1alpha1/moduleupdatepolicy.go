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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ModuleUpdatePolicyLister helps list ModuleUpdatePolicies.
// All objects returned here must be treated as read-only.
type ModuleUpdatePolicyLister interface {
	// List lists all ModuleUpdatePolicies in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ModuleUpdatePolicy, err error)
	// Get retrieves the ModuleUpdatePolicy from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ModuleUpdatePolicy, error)
	ModuleUpdatePolicyListerExpansion
}

// moduleUpdatePolicyLister implements the ModuleUpdatePolicyLister interface.
type moduleUpdatePolicyLister struct {
	indexer cache.Indexer
}

// NewModuleUpdatePolicyLister returns a new ModuleUpdatePolicyLister.
func NewModuleUpdatePolicyLister(indexer cache.Indexer) ModuleUpdatePolicyLister {
	return &moduleUpdatePolicyLister{indexer: indexer}
}

// List lists all ModuleUpdatePolicies in the indexer.
func (s *moduleUpdatePolicyLister) List(selector labels.Selector) (ret []*v1alpha1.ModuleUpdatePolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ModuleUpdatePolicy))
	})
	return ret, err
}

// Get retrieves the ModuleUpdatePolicy from the index for a given name.
func (s *moduleUpdatePolicyLister) Get(name string) (*v1alpha1.ModuleUpdatePolicy, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("moduleupdatepolicy"), name)
	}
	return obj.(*v1alpha1.ModuleUpdatePolicy), nil
}
