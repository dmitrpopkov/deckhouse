// Copyright 2023 Flant JSC
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

package backend

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"

	logger "github.com/docker/distribution/context"
	"github.com/flant/addon-operator/pkg/kube_config_manager/config"
	"github.com/flant/addon-operator/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/client/clientset/versioned"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/client/informers/externalversions"
)

type ModuleConfig struct {
	mcKubeClient *versioned.Clientset
	logger       logger.Logger
}

// New returns native(Deckhouse) implementation for addon-operator's KubeConfigManager which works directly with
// deckhouse.io/ModuleConfig, avoiding moving configs to the ConfigMap
func New(config *rest.Config, logger logger.Logger) *ModuleConfig {
	mcClient, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &ModuleConfig{
		mcClient,
		logger,
	}
}

func (mc ModuleConfig) StartInformer(ctx context.Context, eventC chan config.Event) {
	// define resyncPeriod for informer
	resyncPeriod := time.Duration(15) * time.Minute

	informer := externalversions.NewSharedInformerFactory(mc.mcKubeClient, resyncPeriod)
	mcInformer := informer.Deckhouse().V1alpha1().ModuleConfigs().Informer()

	// we can ignore the error here because we have only 1 error case here:
	//   if mcInformer was stopped already. But we are controlling its behavior
	_, _ = mcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mconfig := obj.(*v1alpha1.ModuleConfig)
			mc.handleEvent(mconfig, eventC)
		},
		UpdateFunc: func(prevObj interface{}, obj interface{}) {
			mconfig := obj.(*v1alpha1.ModuleConfig)
			mc.handleEvent(mconfig, eventC)
		},
		DeleteFunc: func(obj interface{}) {
			mc.handleEvent(obj.(*v1alpha1.ModuleConfig), eventC)
		},
	})

	go func() {
		mcInformer.Run(ctx.Done())
	}()
}

func (mc ModuleConfig) handleEvent(obj *v1alpha1.ModuleConfig, eventC chan config.Event) {
	cfg := config.NewConfig()
	values := utils.Values{}

	// if ModuleConfig was deleted - values are empty
	if obj.DeletionTimestamp == nil {
		// convert values on the fly
		chain := conversion.Registry().Chain(obj.Name)
		fmt.Println("CHAIN VERSION", chain.LatestVersion(), obj.Spec.Version)
		if chain.LatestVersion() != obj.Spec.Version {
			newVersion, newSettings, err := chain.ConvertToLatest(obj.Spec.Version, obj.Spec.Settings)
			if err != nil {
				// TODO: handle panic
				panic(err)
			}
			obj.Spec.Version = newVersion
			obj.Spec.Settings = newSettings
		}

		values = utils.Values(obj.Spec.Settings)
	}

	switch obj.Name {
	case "global":
		cfg.Global = &config.GlobalKubeConfig{
			Values:   values,
			Checksum: values.Checksum(),
		}

	default:
		mcfg := utils.NewModuleConfig(obj.Name, values)
		mcfg.IsEnabled = obj.Spec.Enabled
		cfg.Modules[obj.Name] = &config.ModuleKubeConfig{
			ModuleConfig: *mcfg,
			Checksum:     mcfg.Checksum(),
		}
	}
	eventC <- config.Event{Key: obj.Name, Config: cfg}
}

func (mc ModuleConfig) LoadConfig(ctx context.Context) (*config.KubeConfig, error) {
	// List all ModuleConfig and get settings
	cfg := config.NewConfig()

	list, err := mc.mcKubeClient.DeckhouseV1alpha1().ModuleConfigs().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range list.Items {
		// convert values on the fly
		chain := conversion.Registry().Chain(item.Name)
		fmt.Println("CHAIN VERSION2", chain.LatestVersion(), item.Spec.Version)
		if chain.LatestVersion() != item.Spec.Version {
			newVersion, newSettings, err := chain.ConvertToLatest(item.Spec.Version, item.Spec.Settings)
			if err != nil {
				// TODO: handle panic
				panic(err)
			}
			item.Spec.Version = newVersion
			item.Spec.Settings = newSettings
		}

		values := utils.Values(item.Spec.Settings)

		if item.Name == "global" {
			cfg.Global = &config.GlobalKubeConfig{
				Values:   values,
				Checksum: values.Checksum(),
			}
		} else {
			mcfg := utils.NewModuleConfig(item.Name, values)
			mcfg.IsEnabled = item.Spec.Enabled
			cfg.Modules[item.Name] = &config.ModuleKubeConfig{
				ModuleConfig: *mcfg,
				Checksum:     mcfg.Checksum(),
			}
		}
	}

	return cfg, nil
}

//func (mc ModuleConfig)

// SaveConfigValues saving patches in ModuleConfig. Used for settings-conversions
func (mc ModuleConfig) SaveConfigValues(_ context.Context, moduleName string, values utils.Values) ( /*checksum*/ string, error) {
	mc.logger.Errorf("module %s tries to save values in ModuleConfig: %s", moduleName, values.DebugString())
	return "", errors.New("saving patch values in ModuleConfig is forbidden")
	//fmt.Println("TRY TO update", moduleName, values)
	//mc.logger.Debugf("Saving config values %v for %s", values, moduleName)
	//
	//obj, err := mc.mcKubeClient.DeckhouseV1alpha1().ModuleConfigs().Get(ctx, moduleName, metav1.GetOptions{})
	//if err != nil {
	//	return "", err
	//}
	//
	//// values are stored like: map[prometheus:map[longtermRetentionDays:0 retentionDays:7]]
	//// we have to extract top level key
	//fmt.Println("TRY TO 0", values.HasKey(moduleName))
	//fmt.Println("TRY TO 01", values[moduleName])
	//if values.HasKey(moduleName) {
	//	fmt.Println("TRY TO1 HAS KEY", moduleName)
	//	values = values.SectionByKey(moduleName)
	//	fmt.Println("TRY TO2 AFTER", values)
	//}
	//
	//obj.Spec.Settings = v1alpha1.SettingsValues(values)
	//
	//_, err = mc.mcKubeClient.DeckhouseV1alpha1().ModuleConfigs().Update(ctx, obj, metav1.UpdateOptions{})
	//
	//return values.Checksum(), err
}
