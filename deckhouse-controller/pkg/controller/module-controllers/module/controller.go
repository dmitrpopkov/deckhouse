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

package module

import (
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strconv"
	"sync"
	"time"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/moduleloader"
	d8config "github.com/deckhouse/deckhouse/go_lib/deckhouse-config"
	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
	"github.com/flant/addon-operator/pkg/module_manager/models/modules/events"
	"github.com/flant/addon-operator/pkg/utils/logger"
	metricstorage "github.com/flant/shell-operator/pkg/metric_storage"
	log "github.com/sirupsen/logrus"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	controllerName = "d8-module-controller"

	deleteReleaseAfterHours = 72
)

func RegisterController(runtimeManager manager.Manager, loader *moduleloader.Loader, mm moduleManager, ms *metricstorage.MetricStorage, bundle string) error {
	r := &reconciler{
		init:          sync.WaitGroup{},
		client:        runtimeManager.GetClient(),
		log:           log.WithField("component", "ModuleController"),
		metricStorage: ms,
		moduleManager: mm,
		moduleLoader:  loader,
		bundle:        bundle,
	}

	moduleController, err := controller.New(controllerName, runtimeManager, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// sync modules and configs, and run modules event loop
	r.init.Add(1)
	if err = runtimeManager.Add(manager.RunnableFunc(r.syncModules)); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(runtimeManager).
		For(&v1alpha1.Module{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(moduleController)
}

type reconciler struct {
	init          sync.WaitGroup
	client        client.Client
	log           logger.Logger
	moduleManager moduleManager
	moduleLoader  *moduleloader.Loader
	metricStorage *metricstorage.MetricStorage
	bundle        string
}

type moduleManager interface {
	AreModulesInited() bool
	IsModuleEnabled(moduleName string) bool
	GetModuleNames() []string
	GetModuleEventsChannel() chan events.ModuleEvent
}

func (r *reconciler) syncModules(ctx context.Context) error {
	defer r.init.Done()

	// wait until module manager init
	if err := wait.PollUntilContextCancel(ctx, time.Second, true, func(_ context.Context) (bool, error) {
		return r.moduleManager.AreModulesInited(), nil
	}); err != nil {
		return err
	}

	if err := r.syncModuleConfigs(ctx); err != nil {
		return err
	}

	for _, moduleName := range r.moduleManager.GetModuleNames() {
		module := new(v1alpha1.Module)
		if err := r.client.Get(ctx, types.NamespacedName{Name: moduleName}, module); err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}
			return err
		}
		if err := r.setModuleEnabled(ctx, module, r.moduleManager.IsModuleEnabled(moduleName)); err != nil {
			return err
		}
	}

	return r.runModuleEventLoop(ctx)
}

// syncModuleConfigs syncs module configs at start up
func (r *reconciler) syncModuleConfigs(ctx context.Context) error {
	return retry.OnError(retry.DefaultRetry, apierrors.IsServiceUnavailable, func() error {
		configs := new(v1alpha1.ModuleConfigList)
		if err := r.client.List(ctx, configs); err != nil {
			return err
		}

		for _, config := range configs.Items {
			if err := r.refreshModuleConfigAndModule(ctx, config.Name); err != nil {
				r.log.Errorf("failed to refresh the '%s' module config on sync: %s", config.Name, err)
			}
		}
		return nil
	})
}

func (r *reconciler) refreshModuleConfigAndModule(ctx context.Context, configName string) error {
	return retry.OnError(retry.DefaultRetry, apierrors.IsServiceUnavailable, func() error {
		return retry.RetryOnConflict(retry.DefaultRetry, func() error {
			metricGroup := fmt.Sprintf("%s_%s", "obsoleteVersion", configName)
			r.metricStorage.GroupedVault.ExpireGroupMetrics(metricGroup)

			config := new(v1alpha1.ModuleConfig)
			if err := r.client.Get(ctx, client.ObjectKey{Name: configName}, config); err != nil {
				if !apierrors.IsNotFound(err) {
					return err
				}
				return nil
			}

			newConfigStatus := d8config.Service().StatusReporter().ForConfig(config)
			if (config.Status.Message != newConfigStatus.Message) || (config.Status.Version != newConfigStatus.Version) {
				config.Status.Message = newConfigStatus.Message
				config.Status.Version = newConfigStatus.Version

				if err := r.client.Status().Update(ctx, config); err != nil {
					return err
				}
			}

			// update metrics
			converter := conversion.Store().Get(config.Name)
			if config.Spec.Version > 0 && config.Spec.Version < converter.LatestVersion() {
				r.metricStorage.GroupedVault.GaugeSet(metricGroup, "module_config_obsolete_version", 1.0, map[string]string{
					"name":    config.Name,
					"version": strconv.Itoa(config.Spec.Version),
					"latest":  strconv.Itoa(converter.LatestVersion()),
				})
			}

			// refresh related module
			if err := r.refreshModuleByModuleConfig(ctx, configName); err != nil && !apierrors.IsNotFound(err) {
				return err
			}

			return nil
		})
	})
}

func (r *reconciler) refreshModuleByModuleConfig(ctx context.Context, moduleName string) error {
	module := new(v1alpha1.Module)
	if err := r.client.Get(ctx, client.ObjectKey{Name: moduleName}, module); err != nil {
		return err
	}

	config := new(v1alpha1.ModuleConfig)
	if err := r.client.Get(ctx, client.ObjectKey{Name: moduleName}, config); err != nil {
		if apierrors.IsNotFound(err) {
			config = nil
		} else {
			return err
		}
	}

	// TODO(ipaqsa): update module status reporter
	newModuleStatus := d8config.Service().StatusReporter().ForModule(module, config, r.bundle)
	if module.Status.Phase != newModuleStatus.Phase || module.Status.Message != newModuleStatus.Message || module.Status.HooksState != newModuleStatus.HooksState {
		module.Status.Phase = newModuleStatus.Phase
		module.Status.Message = newModuleStatus.Message
		module.Status.HooksState = newModuleStatus.HooksState

		//r.log.Debugf("update status for the '%s' module: status '%s' to '%s', message '%s' to '%s'", moduleName, module.Status.Status, newModuleStatus.Status, module.Status.Message, newModuleStatus.Message)

		if err := r.client.Status().Update(ctx, module); err != nil {
			return err
		}
	}
	return nil
}

func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (ctrl.Result, error) {
	// wait for init
	r.init.Wait()

	module := new(v1alpha1.Module)
	if err := r.client.Get(ctx, req.NamespacedName, module); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	if !module.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	return r.handleModule(ctx, module)
}

func (r *reconciler) handleModule(ctx context.Context, module *v1alpha1.Module) (ctrl.Result, error) {
	if module.IsEmbedded() {
		return ctrl.Result{}, nil
	}

	if module.IsEnabled() {
		return r.handleEnabledModule(ctx, module)
	}

	if module.DisabledMoreThan(deleteReleaseAfterHours) {
		moduleReleases := new(v1alpha1.ModuleReleaseList)
		if err := r.client.List(ctx, moduleReleases, &client.MatchingLabels{"module": module.Name}); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		for _, release := range moduleReleases.Items {
			if err := r.client.Delete(ctx, &release); err != nil {
				return ctrl.Result{Requeue: true}, err
			}
		}
	}

	if err := r.setModulePhase(ctx, module, v1alpha1.ModulePhaseNotInstalled); err != nil {
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

func (r *reconciler) handleEnabledModule(ctx context.Context, module *v1alpha1.Module) (ctrl.Result, error) {
	if module.Properties.Source == "" {
		if len(module.Properties.AvailableSources) == 1 {
			module.Properties.Source = module.Properties.AvailableSources[0]
		} else {
			if module.Properties.Source = module.Spec.SelectedSource; module.Properties.Source == "" {
				if err := r.setModulePhase(ctx, module, v1alpha1.ModulePhaseConflict); err != nil {
					return ctrl.Result{Requeue: true}, nil
				}
				return ctrl.Result{}, nil
			}
		}
	} else {
		if module.Spec.SelectedSource != "" && module.Spec.SelectedSource != module.Properties.Source {
			module.Properties.Source = module.Spec.SelectedSource
		} else {
			return ctrl.Result{}, nil
		}
	}

	if err := r.setModulePhase(ctx, module, v1alpha1.ModulePhaseDownloading); err != nil {
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

func (r *reconciler) setModulePhase(ctx context.Context, module *v1alpha1.Module, phase string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := r.client.Get(ctx, types.NamespacedName{Name: module.Name}, module); err != nil {
			return err
		}
		module.Status.Phase = string(phase)
		return r.client.Status().Update(ctx, module)
	})
}

func (r *reconciler) setModuleEnabled(ctx context.Context, module *v1alpha1.Module, enabled bool) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := r.client.Get(ctx, types.NamespacedName{Name: module.Name}, module); err != nil {
			return err
		}
		if (module.IsEnabled() && enabled) || (!module.IsEnabled() && !enabled) {
			return nil
		}
		module.SetEnabled(enabled)
		return r.client.Status().Update(ctx, module)
	})
}
