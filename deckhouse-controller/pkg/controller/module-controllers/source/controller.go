// Copyright 2024 Flant JSC
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

package source

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/flant/addon-operator/pkg/utils/logger"
	"github.com/gofrs/uuid/v5"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/downloader"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/release"
	controllerUtils "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/utils"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/helpers"
	d8env "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/env"
	"github.com/deckhouse/deckhouse/go_lib/dependency"
)

const (
	controllerName = "d8-source-controller"

	deckhouseNamespace = "d8-system"

	deckhouseDiscoverySecret = "deckhouse-discovery"

	moduleSourceFinalizer = "modules.deckhouse.io/release-exists"

	moduleSourceAnnotationForceDelete      = "modules.deckhouse.io/force-delete"
	moduleSourceAnnotationRegistryChecksum = "modules.deckhouse.io/registry-spec-checksum"

	defaultScanInterval = 3 * time.Minute

	maxConcurrentReconciles = 3
	cacheSyncTimeout        = 3 * time.Minute
)

func RegisterController(runtimeManager manager.Manager, dc dependency.Container, embeddedPolicy *helpers.ModuleUpdatePolicySpecContainer, preflightCountDown *sync.WaitGroup) error {
	r := &reconciler{
		client:               runtimeManager.GetClient(),
		log:                  log.WithField("component", "ModuleSourceController"),
		downloadedModulesDir: d8env.GetDownloadedModulesDir(),
		embeddedPolicy:       embeddedPolicy,
		dependencyContainer:  dc,
	}

	preflightCountDown.Add(1)

	// add preflight to set the cluster UUID
	if err := runtimeManager.Add(manager.RunnableFunc(func(ctx context.Context) error {
		return r.setClusterUUID(ctx, preflightCountDown)
	})); err != nil {
		return err
	}

	sourceController, err := controller.New(controllerName, runtimeManager, controller.Options{
		MaxConcurrentReconciles: maxConcurrentReconciles,
		CacheSyncTimeout:        cacheSyncTimeout,
		NeedLeaderElection:      ptr.To(false),
		Reconciler:              r,
	})
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(runtimeManager).
		For(&v1alpha1.ModuleSource{}).
		Watches(&v1alpha1.Module{}, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
			return []reconcile.Request{{NamespacedName: types.NamespacedName{Name: obj.(*v1alpha1.Module).Properties.Source}}}
		}), builder.WithPredicates(predicate.Funcs{
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				oldModule := updateEvent.ObjectOld.(*v1alpha1.Module)
				module := updateEvent.ObjectNew.(*v1alpha1.Module)
				if !oldModule.IsEnabled() && module.IsEnabled() {
					return true
				}
				return false
			},
		})).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(sourceController)
}

type reconciler struct {
	client              client.Client
	log                 logger.Logger
	dependencyContainer dependency.Container
	embeddedPolicy      *helpers.ModuleUpdatePolicySpecContainer

	downloadedModulesDir string
	clusterUUID          string
}

func (r *reconciler) setClusterUUID(ctx context.Context, preflight *sync.WaitGroup) error {
	defer preflight.Done()

	// attempt to read the cluster UUID from a secret
	secret := new(corev1.Secret)
	if err := r.client.Get(ctx, types.NamespacedName{Namespace: deckhouseNamespace, Name: deckhouseDiscoverySecret}, secret); err != nil {
		r.log.Warnf("failed to read clusterUUID from the 'deckhouse-discovery' secret: %v. Generating random uuid", err)
		r.clusterUUID = uuid.Must(uuid.NewV4()).String()
		return nil
	}

	if clusterUUID, ok := secret.Data["clusterUUID"]; ok {
		r.clusterUUID = string(clusterUUID)
		return nil
	}

	// generate a random UUID if the key is missing
	r.clusterUUID = uuid.Must(uuid.NewV4()).String()
	return nil
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	moduleSource := new(v1alpha1.ModuleSource)
	if err := r.client.Get(ctx, req.NamespacedName, moduleSource); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, err
	}

	// handle delete event
	if !moduleSource.DeletionTimestamp.IsZero() {
		return r.deleteModuleSource(ctx, moduleSource)
	}

	// handle create/update events
	return r.handleModuleSource(ctx, moduleSource)
}

func (r *reconciler) handleModuleSource(ctx context.Context, source *v1alpha1.ModuleSource) (ctrl.Result, error) {
	// reset status fields
	source.Status.Msg = ""
	source.Status.ModuleErrors = make([]v1alpha1.ModuleError, 0)

	// generate options for connecting to the registry
	opts := controllerUtils.GenerateRegistryOptionsFromModuleSource(source, r.clusterUUID)

	// create a registry client
	registryClient, err := r.dependencyContainer.GetRegistryClient(source.Spec.Registry.Repo, opts...)
	if err != nil {
		source.Status.Msg = err.Error()
		if err = r.updateModuleSourceStatus(ctx, source); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		// error can occur on wrong auth only, we don't want to requeue the source until auth is fixed
		return ctrl.Result{Requeue: false}, nil
	}

	// sync registry settings if they have changed check
	shouldRequeue, err := r.syncRegistrySettings(ctx, source)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if shouldRequeue {
		// new registry settings checksum should be applied to module source
		if err = r.client.Update(ctx, source); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		// requeue moduleSource after modifying annotation
		return ctrl.Result{Requeue: true}, nil
	}

	// list available modules(tags) from the registry
	moduleNames, err := registryClient.ListTags(ctx)
	if err != nil {
		source.Status.Msg = err.Error()
		if err = r.updateModuleSourceStatus(ctx, source); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true}, err
	}

	sort.Strings(moduleNames)

	availableModules := make([]v1alpha1.AvailableModule, 0, len(moduleNames))
	source.Status.ModulesCount = len(moduleNames)

	// get all policies regardless of their labels
	policies := new(v1alpha1.ModuleUpdatePolicyList)
	if err = r.client.List(ctx, policies); err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	md := downloader.NewModuleDownloader(r.dependencyContainer, r.downloadedModulesDir, source, opts)

	for _, moduleName := range moduleNames {
		if moduleName == "modules" {
			r.log.Warn("the 'modules' name is a forbidden name. Skip the module.")
			continue
		}

		available, err := r.processModule(ctx, md, source, moduleName, policies.Items)
		if err != nil {
			// collect modules errors for reporting
			source.Status.ModuleErrors = append(source.Status.ModuleErrors, v1alpha1.ModuleError{
				Name:  moduleName,
				Error: err.Error(),
			})
		}
		availableModules = append(availableModules, available)
	}

	source.Status.AvailableModules = availableModules

	if len(source.Status.ModuleErrors) > 0 {
		source.Status.Msg = "Some errors occurred. Inspect status for details"
	}

	if err = r.updateModuleSourceStatus(ctx, source); err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	// everything is ok, check source on the other iteration
	return ctrl.Result{RequeueAfter: defaultScanInterval}, nil
}

func (r *reconciler) processModule(ctx context.Context, md *downloader.ModuleDownloader, source *v1alpha1.ModuleSource, moduleName string, policies []v1alpha1.ModuleUpdatePolicy) (v1alpha1.AvailableModule, error) {
	availableModule := v1alpha1.AvailableModule{Name: moduleName}
	for _, available := range source.Status.AvailableModules {
		if available.Name == moduleName {
			availableModule = available
		}
	}

	// get an update policy for the module or, if there is no matching policy, use the embedded on
	policy, err := r.releasePolicy(source.Name, moduleName, policies)
	if err != nil {
		return availableModule, err
	}
	availableModule.Policy = policy.Name

	// skip processing if the policy mode is "Ignore"
	if policy.Spec.Update.Mode == v1alpha1.ModuleUpdatePolicyModeIgnore {
		return availableModule, nil
	}

	// add the source to the module or create a new module with this source
	ensureRelease, err := r.ensureModule(ctx, source.Name, policy.Spec.ReleaseChannel, moduleName)
	if err != nil {
		return availableModule, err
	}

	// download module metadata from the specified release channel
	meta, err := md.DownloadMetadataFromReleaseChannel(moduleName, policy.Spec.ReleaseChannel, availableModule.Checksum)
	if err != nil {
		return availableModule, err
	}

	// if release is changed ensure module release
	if ensureRelease && availableModule.Checksum != meta.Checksum {
		availableModule.Checksum = meta.Checksum
		if err = r.ensureModuleRelease(ctx, source.Name, source.GetUID(), moduleName, policy.Name, meta); err != nil {
			return availableModule, err
		}
	}

	return availableModule, nil
}

func (r *reconciler) deleteModuleSource(ctx context.Context, source *v1alpha1.ModuleSource) (ctrl.Result, error) {
	if controllerutil.ContainsFinalizer(source, moduleSourceFinalizer) {
		if source.GetAnnotations()[moduleSourceAnnotationForceDelete] != "true" {
			// list deployed ModuleReleases associated with the ModuleSource
			releases := new(v1alpha1.ModuleReleaseList)
			if err := r.client.List(ctx, releases, client.MatchingLabels{"source": source.Name, "status": "deployed"}); err != nil {
				return ctrl.Result{Requeue: true}, err
			}

			// prevent deletion if there are deployed releases
			if len(releases.Items) > 0 {
				source.Status.Msg = "The ModuleSource contains at least 1 Deployed release and cannot be deleted. Please delete target ModuleReleases manually to continue"
				if err := r.updateModuleSourceStatus(ctx, source); err != nil {
					return ctrl.Result{Requeue: true}, nil
				}
				return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
			}

			for _, module := range source.Status.AvailableModules {
				if err := r.cleanSourceInModule(ctx, source.Name, module.Name); err != nil {
					return ctrl.Result{Requeue: true}, err
				}
			}
		}

		controllerutil.RemoveFinalizer(source, moduleSourceFinalizer)

		if err := r.client.Update(ctx, source); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *reconciler) cleanSourceInModule(ctx context.Context, sourceName, moduleName string) error {
	module := new(v1alpha1.Module)
	if err := r.client.Get(ctx, client.ObjectKey{Name: moduleName}, module); err != nil {
		return err
	}

	for idx, source := range module.Properties.AvailableSources {
		if source == sourceName {
			module.Properties.AvailableSources = append(module.Properties.AvailableSources[:idx], module.Properties.AvailableSources[idx+1:]...)
		}
	}
	if len(module.Properties.AvailableSources) == 0 {
		if err := r.client.Delete(ctx, module); err != nil {
			return err
		}
	}
	if err := r.client.Update(ctx, module); err != nil {
		return err
	}
	return nil
}

// releasePolicy checks if any update policy matches the module source and if it's so - returns the policy.
// if several policies match the module source labels, return error
// if no policy matches the module source, embeddedPolicy is returned
func (r *reconciler) releasePolicy(sourceName, moduleName string, policies []v1alpha1.ModuleUpdatePolicy) (*v1alpha1.ModuleUpdatePolicy, error) {
	var releaseLabelsSet labels.Set = map[string]string{"module": moduleName, "source": sourceName}
	var matchedPolicy *v1alpha1.ModuleUpdatePolicy
	for _, policy := range policies {
		if policy.Spec.ModuleReleaseSelector.LabelSelector == nil {
			continue
		}

		selector, err := metav1.LabelSelectorAsSelector(policy.Spec.ModuleReleaseSelector.LabelSelector)
		if err != nil {
			return nil, err
		}

		if selectorSourceName, exist := selector.RequiresExactMatch("source"); exist && selectorSourceName != sourceName {
			// the 'source' label is set, but does not match the given ModuleSource
			continue
		}

		if selector.Matches(releaseLabelsSet) {
			// ModuleUpdatePolicy matches ModuleSource and specified Module
			if matchedPolicy != nil {
				return nil, fmt.Errorf("more than one update policy matches the module: %s and %s", matchedPolicy.Name, policy.Name)
			}
			matchedPolicy = &policy
		}
	}

	if matchedPolicy == nil {
		r.log.Infof("no module update policy for the '%q' module source, and the '%q' module, deckhouse policy will be used: %+v", sourceName, moduleName, *r.embeddedPolicy.Get())
		return &v1alpha1.ModuleUpdatePolicy{
			TypeMeta: metav1.TypeMeta{
				Kind:       v1alpha1.ModuleUpdatePolicyGVK.Kind,
				APIVersion: v1alpha1.ModuleUpdatePolicyGVK.GroupVersion().String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "", // special empty default policy, inherits Deckhouse settings for update mode
			},
			Spec: *r.embeddedPolicy.Get(),
		}, nil
	}

	return matchedPolicy, nil
}

// syncRegistrySettings checks if modules source registry settings were updated (comparing moduleSourceAnnotationRegistryChecksum annotation and the current registry spec)
// and update relevant module releases' openapi values files if it is the case
func (r *reconciler) syncRegistrySettings(ctx context.Context, source *v1alpha1.ModuleSource) (bool, error) {
	marshaled, err := json.Marshal(source.Spec.Registry)
	if err != nil {
		return false, fmt.Errorf("failed to marshal the '%s' module source registry spec: %w", source.Name, err)
	}

	currentChecksum := fmt.Sprintf("%x", md5.Sum(marshaled))

	// if there is no annotations - only set the current checksum value
	if source.ObjectMeta.Annotations == nil {
		source.ObjectMeta.Annotations = make(map[string]string)
		source.ObjectMeta.Annotations[moduleSourceAnnotationRegistryChecksum] = currentChecksum
		return true, nil
	}

	// if the annotation matches current checksum - there is nothing to do here
	if source.ObjectMeta.Annotations[moduleSourceAnnotationRegistryChecksum] == currentChecksum {
		return false, nil
	}

	// get related releases
	moduleReleases := new(v1alpha1.ModuleReleaseList)
	if err = r.client.List(ctx, moduleReleases, client.MatchingLabels{"source": source.Name}); err != nil {
		return false, fmt.Errorf("failed to list module releases to update registry settings: %w", err)
	}

	for _, moduleRelease := range moduleReleases.Items {
		if moduleRelease.Status.Phase == v1alpha1.PhaseDeployed {
			for _, ref := range moduleRelease.GetOwnerReferences() {
				if ref.UID == source.UID && ref.Name == source.Name && ref.Kind == v1alpha1.ModuleSourceGVK.Kind {
					// update the values.yaml file in external-modules/<module_name>/v<module_version/openapi path
					modulePath := filepath.Join(r.downloadedModulesDir, moduleRelease.Spec.ModuleName, fmt.Sprintf("v%s", moduleRelease.Spec.Version))
					if err = downloader.InjectRegistryToModuleValues(modulePath, source); err != nil {
						return false, fmt.Errorf("failed to update the '%s' module release registry settings: %w", moduleRelease.Name, err)
					}

					if moduleRelease.ObjectMeta.Annotations == nil {
						moduleRelease.ObjectMeta.Annotations = make(map[string]string)
					}

					moduleRelease.ObjectMeta.Annotations[release.RegistrySpecChangedAnnotation] = r.dependencyContainer.GetClock().Now().UTC().Format(time.RFC3339)
					if err = r.client.Update(ctx, &moduleRelease); err != nil {
						return false, fmt.Errorf("failed to set RegistrySpecChangedAnnotation to the '%s' module release: %w", moduleRelease.Name, err)
					}
					break
				}
			}
		}
	}

	source.ObjectMeta.Annotations[moduleSourceAnnotationRegistryChecksum] = currentChecksum

	return true, nil
}

func (r *reconciler) updateModuleSourceStatus(ctx context.Context, sourceCopy *v1alpha1.ModuleSource) error {
	sourceCopy.Status.SyncTime = metav1.NewTime(r.dependencyContainer.GetClock().Now().UTC())
	return r.client.Status().Update(ctx, sourceCopy)
}

func (r *reconciler) ensureModule(ctx context.Context, sourceName, releaseChannel string, moduleName string) (bool, error) {
	module := new(v1alpha1.Module)
	err := r.client.Get(ctx, client.ObjectKey{Name: moduleName}, module)
	if err == nil {
		if !slices.Contains(module.Properties.AvailableSources, sourceName) {
			module.Properties.AvailableSources = append(module.Properties.AvailableSources, moduleName)
			// TODO(ipaqsa): add retry on conflict
			if err = r.client.Update(ctx, module); err != nil {
				return false, err
			}
		}
		if module.Properties.Source == sourceName {
			module.Properties.ReleaseChannel = releaseChannel
			if err = r.client.Update(ctx, module); err != nil {
				return false, err
			}
		}
		return module.IsEnabled() && module.Properties.Source == sourceName, nil
	}

	if apierrors.IsNotFound(err) {
		module = &v1alpha1.Module{
			TypeMeta: metav1.TypeMeta{
				Kind:       v1alpha1.ModuleGVK.Kind,
				APIVersion: v1alpha1.ModuleGVK.GroupVersion().String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: moduleName,
			},
			Properties: v1alpha1.ModuleProperties{
				AvailableSources: []string{sourceName},
			},
			Status: v1alpha1.ModuleStatus{
				Phase: v1alpha1.ModulePhaseNotInstalled,
				Conditions: []v1alpha1.ModuleCondition{
					{
						Type:               v1alpha1.ModuleConditionEnabled,
						Status:             corev1.ConditionFalse,
						LastTransitionTime: metav1.Now(),
						LastProbeTime:      metav1.Now(),
					},
				},
			},
		}
		if err = r.client.Create(ctx, module); err != nil {
			return false, err
		}
		return false, nil
	}

	return false, err
}

func (r *reconciler) ensureModuleRelease(ctx context.Context, sourceName string, sourceUID types.UID, moduleName, policy string, meta downloader.ModuleDownloadResult) error {
	// image digest has 64 symbols, while label can have maximum 63 symbols, so make md5 sum here
	checksum := fmt.Sprintf("%x", md5.Sum([]byte(meta.Checksum)))

	moduleRelease := &v1alpha1.ModuleRelease{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ModuleRelease",
			APIVersion: "deckhouse.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", moduleName, meta.ModuleVersion),
			Labels: map[string]string{
				"module":                  moduleName,
				"source":                  sourceName,
				"release-checksum":        checksum,
				release.UpdatePolicyLabel: policy,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: v1alpha1.ModuleSourceGVK.GroupVersion().String(),
					Kind:       v1alpha1.ModuleSourceGVK.Kind,
					Name:       sourceName,
					UID:        sourceUID,
					Controller: ptr.To(true),
				},
			},
		},
		Spec: v1alpha1.ModuleReleaseSpec{
			ModuleName: moduleName,
			Version:    semver.MustParse(meta.ModuleVersion),
			Weight:     meta.ModuleWeight,
			Changelog:  v1alpha1.Changelog(meta.Changelog),
		},
	}
	if meta.ModuleDefinition != nil {
		moduleRelease.Spec.Requirements = meta.ModuleDefinition.Requirements
	}

	if err := r.client.Create(ctx, moduleRelease); err != nil {
		if apierrors.IsAlreadyExists(err) {
			prevModuleRelease := new(v1alpha1.ModuleRelease)
			if err = r.client.Get(ctx, client.ObjectKey{Name: moduleRelease.Name}, prevModuleRelease); err != nil {
				return err
			}

			// seems weird to update already deployed/suspended release
			if prevModuleRelease.Status.Phase != v1alpha1.PhasePending {
				return nil
			}

			prevModuleRelease.Spec = moduleRelease.Spec
			return r.client.Update(ctx, prevModuleRelease)
		}

		return err
	}
	return nil
}
