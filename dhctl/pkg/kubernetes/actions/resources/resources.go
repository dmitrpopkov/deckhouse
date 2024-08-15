// Copyright 2021 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/deckhouse/deckhouse/dhctl/pkg/config"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"

	"github.com/deckhouse/deckhouse/dhctl/pkg/app"
	"github.com/deckhouse/deckhouse/dhctl/pkg/kubernetes/actions"
	"github.com/deckhouse/deckhouse/dhctl/pkg/kubernetes/client"
	"github.com/deckhouse/deckhouse/dhctl/pkg/log"
	"github.com/deckhouse/deckhouse/dhctl/pkg/template"
	"github.com/deckhouse/deckhouse/dhctl/pkg/util/retry"
)

var ErrNotAllResourcesCreated = fmt.Errorf("Not all resources were creatated")

// apiResourceListGetter discovery and cache APIResources list for group version kind
type apiResourceListGetter struct {
	kubeCl             *client.KubernetesClient
	gvkToResourcesList map[string]*metav1.APIResourceList
}

func newAPIResourceListGetter(kubeCl *client.KubernetesClient) *apiResourceListGetter {
	return &apiResourceListGetter{
		kubeCl:             kubeCl,
		gvkToResourcesList: make(map[string]*metav1.APIResourceList),
	}
}

func (g *apiResourceListGetter) Get(gvk *schema.GroupVersionKind) (*metav1.APIResourceList, error) {
	key := gvk.GroupVersion().String()
	if resourcesList, ok := g.gvkToResourcesList[key]; ok {
		return resourcesList, nil
	}

	var resourcesList *metav1.APIResourceList
	var err error
	err = retry.NewSilentLoop("Get resources list", 5, 5*time.Second).Run(func() error {
		// ServerResourcesForGroupVersion does not return error if API returned NotFound (404) or Forbidden (403)
		// https://github.com/kubernetes/client-go/blob/51a4fd4aee686931f6a53148b3f4c9094f80d512/discovery/discovery_client.go#L204
		// and if CRD was not deployed method will return empty APIResources list
		resourcesList, err = g.kubeCl.Discovery().ServerResourcesForGroupVersion(gvk.GroupVersion().String())
		if err != nil {
			return fmt.Errorf("can't get preferred resources '%s': %w", key, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resourcesList, nil
}

type Creator struct {
	kubeCl    *client.KubernetesClient
	resources []*template.Resource
}

func NewCreator(kubeCl *client.KubernetesClient, resources template.Resources) *Creator {
	return &Creator{
		kubeCl:    kubeCl,
		resources: resources,
	}
}

func (c *Creator) createAll() error {
	apiResourceGetter := newAPIResourceListGetter(c.kubeCl)
	addedResourcesIndexes := make(map[int]struct{})

	defer func() {
		remainResources := make([]*template.Resource, 0)

		for i, resource := range c.resources {
			if _, ok := addedResourcesIndexes[i]; !ok {
				remainResources = append(remainResources, resource)
			}
		}

		c.resources = remainResources
	}()

	if err := c.ensureRequiredNamespacesExist(); err != nil {
		return err
	}

	for indx, resource := range c.resources {
		resourcesList, err := apiResourceGetter.Get(&resource.GVK)
		if err != nil {
			log.DebugF("apiResourceGetter returns error: %w", err)
			continue
		}

		for _, discoveredResource := range resourcesList.APIResources {
			if discoveredResource.Kind != resource.GVK.Kind {
				continue
			}
			if err := c.createSingleResource(resource); err != nil {
				return err
			}

			addedResourcesIndexes[indx] = struct{}{}
			break
		}
	}

	return nil
}

func (c *Creator) ensureRequiredNamespacesExist() error {
	knownNamespaces := make(map[string]struct{})

	for _, res := range c.resources {
		nsName := res.Object.GetNamespace()
		if _, nsWasSeenBefore := knownNamespaces[nsName]; nsName == "" || nsWasSeenBefore {
			continue // If this resource is not namespaces, or we saw this namespace already, there is no need to check
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if _, err := c.kubeCl.CoreV1().Namespaces().Get(ctx, nsName, metav1.GetOptions{}); err != nil {
			cancel()

			if !apierrors.IsNotFound(err) {
				return fmt.Errorf("can't get namespace %q: %w", nsName, err)
			}

			return fmt.Errorf(
				"%w: waiting for namespace %q is to create %q (%s)",
				ErrNotAllResourcesCreated,
				nsName,
				res.Object.GetName(),
				res.GVK.String(),
			)
		}
		cancel()
		knownNamespaces[nsName] = struct{}{}
	}
	return nil
}

func (c *Creator) TryToCreate() error {
	if err := c.createAll(); err != nil {
		return err
	}

	gvks := make(map[string]struct{})
	resourcesToCreate := make([]string, 0, len(c.resources))
	for _, resource := range c.resources {
		key := resource.GVK.String()
		if _, ok := gvks[key]; !ok {
			gvks[key] = struct{}{}
			resourcesToCreate = append(resourcesToCreate, key)
		}
	}

	if len(c.resources) > 0 {
		log.InfoF("\rResources to create: \n\t%s\n\n", strings.Join(resourcesToCreate, "\n\t"))
		return ErrNotAllResourcesCreated
	}

	return nil
}

func (c *Creator) isNamespaced(gvk schema.GroupVersionKind, name string) (bool, error) {
	return isNamespaced(c.kubeCl, gvk, name)
}

func (c *Creator) createSingleResource(resource *template.Resource) error {
	doc := resource.Object
	gvk := resource.GVK

	// Wait up to 10 minutes
	return retry.NewLoop(fmt.Sprintf("Create %s resources", gvk.String()), 60, 10*time.Second).Run(func() error {
		gvr, err := c.kubeCl.GroupVersionResource(gvk.ToAPIVersionAndKind())
		if err != nil {
			return fmt.Errorf("can't get resource by kind and apiVersion: %w", err)
		}

		namespaced, err := c.isNamespaced(gvk, gvr.Resource)
		if err != nil {
			return fmt.Errorf("can't determine whether a resource is namespaced or not: %v", err)
		}

		docCopy := doc.DeepCopy()
		namespace := docCopy.GetNamespace()
		if namespace == metav1.NamespaceNone && namespaced {
			namespace = metav1.NamespaceDefault
		}

		manifestTask := actions.ManifestTask{
			Name:     getUnstructuredName(docCopy),
			Manifest: func() interface{} { return nil },
			CreateFunc: func(manifest interface{}) error {
				_, err := c.kubeCl.Dynamic().Resource(gvr).
					Namespace(namespace).
					Create(context.TODO(), docCopy, metav1.CreateOptions{})
				return err
			},
			UpdateFunc: func(manifest interface{}) error {
				content, err := docCopy.MarshalJSON()
				if err != nil {
					return err
				}
				// using patch here because of https://github.com/kubernetes/kubernetes/issues/70674
				_, err = c.kubeCl.Dynamic().Resource(gvr).
					Namespace(namespace).
					Patch(context.TODO(), docCopy.GetName(), types.MergePatchType, content, metav1.PatchOptions{})
				return err
			},
		}

		return manifestTask.CreateOrUpdate()
	})
}

func CreateResourcesLoop(kubeCl *client.KubernetesClient, metaConfig *config.MetaConfig, resources template.Resources, checkers []Checker) error {
	ctx, cancel := context.WithTimeout(context.TODO(), app.ResourcesTimeout)
	defer cancel()

	resourceCreator := NewCreator(kubeCl, resources)

	clusterIsBootstrappedChecker, err := tryToGetClusterIsBootstrappedChecker(kubeCl, resources, metaConfig)
	if err != nil {
		return err
	}

	isBootstrapped := false

	err = retry.NewSilentLoop("Wait for resources creation", math.MaxInt, 10*time.Second).
		WithContext(ctx).
		Run(func() error {
			err := resourceCreator.TryToCreate()
			if err != nil && !errors.Is(err, ErrNotAllResourcesCreated) {
				return err
			}

			if clusterIsBootstrappedChecker != nil {
				var err error

				if !isBootstrapped {
					isBootstrapped, err = clusterIsBootstrappedChecker.IsBootstrapped()
				}

				if err != nil {
					return err
				}
			} else {
				isBootstrapped = true
			}

			if isBootstrapped && err == nil {
				return nil
			}

			if clusterIsBootstrappedChecker != nil {
				clusterIsBootstrappedChecker.outputProgressInfo()
			}

			return ErrNotAllResourcesCreated
		})
	if err != nil {
		return err
	}

	waiter := NewWaiter(checkers)

	err = retry.NewSilentLoop("Wait for resources readiness", math.MaxInt, 5*time.Second).
		WithContext(ctx).
		Run(func() error {
			ready, err := waiter.ReadyAll(ctx)
			if err != nil {
				return err
			}

			if !ready {
				return fmt.Errorf("not all resources are ready")
			}

			return nil
		})
	if err != nil {
		return err
	}

	return nil
}

func getUnstructuredName(obj *unstructured.Unstructured) string {
	namespace := obj.GetNamespace()
	if namespace == "" {
		return fmt.Sprintf("%s %s", obj.GetKind(), obj.GetName())
	}
	return fmt.Sprintf("%s %s/%s", obj.GetKind(), obj.GetNamespace(), obj.GetName())
}

func DeleteResourcesLoop(ctx context.Context, kubeCl *client.KubernetesClient, resources template.Resources) error {
	for _, res := range resources {
		name := res.Object.GetName()
		namespace := res.Object.GetNamespace()
		gvk := res.GVK

		gvr, err := kubeCl.GroupVersionResource(gvk.ToAPIVersionAndKind())
		if err != nil {
			return fmt.Errorf("bad group version resource %s: %w", res.GVK.String(), err)
		}

		namespaced, err := isNamespaced(kubeCl, gvk, gvr.Resource)
		if err != nil {
			return fmt.Errorf("can't determine whether a resource is namespaced or not: %v", err)
		}
		if namespace == metav1.NamespaceNone && namespaced {
			namespace = metav1.NamespaceDefault
		}

		var resourceClient dynamic.ResourceInterface
		if namespaced {
			resourceClient = kubeCl.Dynamic().Resource(gvr).Namespace(namespace)
		} else {
			resourceClient = kubeCl.Dynamic().Resource(gvr)
		}

		if err := resourceClient.Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				return fmt.Errorf("unable to delete %s %s: %w", gvr.String(), name, err)
			}
			log.DebugF("Unable to delete resource: %s %s: %s\n", gvr.String(), name, err)
		}
	}

	return nil
}

func isNamespaced(kubeCl *client.KubernetesClient, gvk schema.GroupVersionKind, name string) (bool, error) {
	lists, err := kubeCl.APIResourceList(gvk.GroupVersion().String())
	if err != nil && len(lists) == 0 {
		// apiVersion is defined and there is a ServerResourcesForGroupVersion error
		return false, err
	}

	for _, list := range lists {
		for _, resource := range list.APIResources {
			if len(resource.Verbs) == 0 {
				continue
			}
			if resource.Name == name {
				return resource.Namespaced, nil
			}
		}
	}
	return false, nil
}
