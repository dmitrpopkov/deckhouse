/*
Copyright 2023 Flant JSC

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

package controller

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	modulemanager "github.com/flant/addon-operator/pkg/module_manager"
	"github.com/flant/addon-operator/pkg/utils"
	metricstorage "github.com/flant/shell-operator/pkg/metric_storage"
	"github.com/go-logr/logr"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	coordv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/yaml"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	deckhouserelease "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/deckhouse-release"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/docbuilder"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/module"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/moduleloader"
	modulerelease "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/release"
	modulesource "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/source"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/helpers"
	d8env "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/env"
	"github.com/deckhouse/deckhouse/go_lib/dependency"
	"github.com/deckhouse/deckhouse/go_lib/dependency/extenders"
)

const (
	docsLeaseLabel = "deckhouse.io/documentation-builder-sync"

	deckhouseNamespace  = "d8-system"
	kubernetesNamespace = "kube-system"
)

type DeckhouseController struct {
	runtimeManager     manager.Manager
	preflightCountDown *sync.WaitGroup

	moduleLoader *moduleloader.Loader

	embeddedPolicy    *helpers.ModuleUpdatePolicySpecContainer
	deckhouseSettings *helpers.DeckhouseSettingsContainer
}

func NewDeckhouseController(ctx context.Context, config *rest.Config, mm *modulemanager.ModuleManager, ms *metricstorage.MetricStorage) (*DeckhouseController, error) {
	addToScheme := []func(s *runtime.Scheme) error{
		corev1.AddToScheme,
		coordv1.AddToScheme,
		v1alpha1.AddToScheme,
		appsv1.AddToScheme,
	}

	scheme := runtime.NewScheme()
	for _, add := range addToScheme {
		if err := add(scheme); err != nil {
			return nil, err
		}
	}

	// Setting the controller-runtime logger to a no-op logger by default,
	// unless debug mode is enabled. This is because the controller-runtime
	// logger is *very* verbose even at info level. This is not really needed,
	// but otherwise we get a warning from the controller-runtime.
	controllerruntime.SetLogger(logr.New(ctrllog.NullLogSink{}))

	runtimeManager, err := controllerruntime.NewManager(config, controllerruntime.Options{
		Scheme: scheme,
		BaseContext: func() context.Context {
			return ctx
		},
		// disable manager's metrics for a while
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
		GracefulShutdownTimeout: ptr.To(10 * time.Second),
		Cache: cache.Options{
			ByObject: map[client.Object]cache.ByObject{
				// for ModuleDocumentation controller
				&coordv1.Lease{}: {
					Namespaces: map[string]cache.Config{
						deckhouseNamespace: {
							LabelSelector: labels.SelectorFromSet(map[string]string{docsLeaseLabel: ""}),
						},
					},
				},
				// for ModuleRelease controller and DeckhouseRelease controller
				&corev1.Secret{}: {
					Namespaces: map[string]cache.Config{
						deckhouseNamespace: {
							LabelSelector: labels.SelectorFromSet(map[string]string{"heritage": "deckhouse", "module": "deckhouse"}),
						},
						kubernetesNamespace: {
							LabelSelector: labels.SelectorFromSet(map[string]string{"name": "d8-cluster-configuration"}),
						},
					},
				},
				// for DeckhouseRelease controller
				&corev1.Pod{}: {
					Namespaces: map[string]cache.Config{
						deckhouseNamespace: {
							LabelSelector: labels.SelectorFromSet(map[string]string{"app": "deckhouse"}),
						},
					},
				},
				// for DeckhouseRelease controller
				&corev1.ConfigMap{}: {
					Namespaces: map[string]cache.Config{
						deckhouseNamespace: {
							LabelSelector: labels.SelectorFromSet(map[string]string{"heritage": "deckhouse"}),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// register extenders
	for _, extender := range extenders.Extenders() {
		if err = mm.AddExtender(extender); err != nil {
			return nil, err
		}
	}

	// create a default policy, it'll be filled in with relevant settings from the deckhouse moduleConfig
	embeddedPolicy := helpers.NewModuleUpdatePolicySpecContainer(&v1alpha1.ModuleUpdatePolicySpec{
		Update: v1alpha1.ModuleUpdatePolicySpecUpdate{
			Mode: "Auto",
		},
		ReleaseChannel: "Stable",
	})

	dc := dependency.NewDependencyContainer()
	dsContainer := helpers.NewDeckhouseSettingsContainer(nil)

	preflightCountDown := &sync.WaitGroup{}

	bundle := os.Getenv("DECKHOUSE_BUNDLE")

	loader := moduleloader.New(runtimeManager.GetClient(), mm.ModulesDir)

	// register module loader
	mm.SetModuleLoader(loader)

	err = deckhouserelease.NewDeckhouseReleaseController(ctx, runtimeManager, dc, mm, dsContainer, ms, preflightCountDown)
	if err != nil {
		return nil, fmt.Errorf("new Deckhouse release controller: %w", err)
	}

	err = module.RegisterController(runtimeManager, loader, mm, ms, bundle)
	if err != nil {
		return nil, err
	}

	err = modulesource.RegisterController(runtimeManager, dc, embeddedPolicy, preflightCountDown)
	if err != nil {
		return nil, err
	}

	err = modulerelease.NewModuleReleaseController(runtimeManager, dc, embeddedPolicy, mm, ms, preflightCountDown)
	if err != nil {
		return nil, err
	}

	err = modulerelease.NewModulePullOverrideController(runtimeManager, dc, mm, preflightCountDown)
	if err != nil {
		return nil, err
	}

	err = docbuilder.NewModuleDocumentationController(runtimeManager, dc)
	if err != nil {
		return nil, err
	}

	return &DeckhouseController{
		runtimeManager:     runtimeManager,
		moduleLoader:       loader,
		preflightCountDown: preflightCountDown,

		embeddedPolicy:    embeddedPolicy,
		deckhouseSettings: dsContainer,
	}, nil
}

// Start loads modules from FS, starts pluggable controllers and runs deckhouse config event loop
func (c *DeckhouseController) Start(ctx context.Context, deckhouseConfigCh <-chan utils.Values) error {
	// run preflight checks first
	if d8env.GetDownloadedModulesDir() != "" {
		c.startPluggableModulesControllers(ctx)
	}

	// load and ensure modules from fs at start
	if err := c.moduleLoader.LoadModulesFromFS(ctx); err != nil {
		return err
	}

	go c.runDeckhouseConfigObserver(deckhouseConfigCh)

	return nil
}

// startPluggableModulesControllers starts all child controllers linked with Modules
func (c *DeckhouseController) startPluggableModulesControllers(ctx context.Context) {
	// syncs the fs with the cluster state, starts the manager and various controllers
	go func() {
		if err := c.runtimeManager.Start(ctx); err != nil {
			log.Fatalf("Start controller manager failed: %s", err)
		}
	}()

	log.Info("Waiting for the preflight checks to run")
	c.preflightCountDown.Wait()
	log.Info("The preflight checks are done")
}

// runDeckhouseConfigObserver updates embeddedPolicy and deckhouseSettings with the configuration from the deckhouse moduleConfig
func (c *DeckhouseController) runDeckhouseConfigObserver(configCh <-chan utils.Values) {
	for {
		cfg := <-configCh

		configBytes, _ := cfg.AsBytes("yaml")
		settings := &helpers.DeckhouseSettings{
			ReleaseChannel: "",
		}
		settings.Update.Mode = "Auto"
		settings.Update.DisruptionApprovalMode = "Auto"

		if err := yaml.Unmarshal(configBytes, settings); err != nil {
			log.Errorf("Error occurred during the Deckhouse settings unmarshalling: %s", err)
			continue
		}

		c.deckhouseSettings.Set(settings)

		// if deckhouse moduleConfig has releaseChannel unset, apply default releaseChannel Stable to the embedded Deckhouse policy
		if len(settings.ReleaseChannel) == 0 {
			settings.ReleaseChannel = "Stable"
			log.Debugf("Embedded deckhouse policy release channel set to %s", settings.ReleaseChannel)
		}
		c.embeddedPolicy.Set(settings)
	}
}

// GetModuleByName implements moduleStorage interface for validation webhook
func (c *DeckhouseController) GetModuleByName(name string) (*moduleloader.Module, error) {
	return c.moduleLoader.GetModuleByName(name)
}
