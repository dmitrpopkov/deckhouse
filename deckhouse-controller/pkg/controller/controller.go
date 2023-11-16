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
	"math/rand"

	"github.com/flant/addon-operator/pkg/module_manager/models/modules/events"
	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/pkg/values/validation"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/client/clientset/versioned"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/models"
)

const (
	epochLabelKey = "deckhouse.io/epoch"
)

var (
	epochLabelValue = fmt.Sprintf("%d", rand.Uint32())
)

type DeckhouseController struct {
	ctx context.Context

	dirs            []string
	valuesValidator *validation.ValuesValidator
	kubeClient      *versioned.Clientset

	deckhouseModules map[string]*models.DeckhouseModule
	// <module-name>: <module-source>
	sourceModules map[string]string
}

func NewDeckhouseController(ctx context.Context, config *rest.Config, moduleDirs string, vv *validation.ValuesValidator) (*DeckhouseController, error) {
	mcClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &DeckhouseController{
		ctx:             ctx,
		kubeClient:      mcClient,
		dirs:            utils.SplitToPaths(moduleDirs),
		valuesValidator: vv,

		deckhouseModules: make(map[string]*models.DeckhouseModule),
		sourceModules:    make(map[string]string),
	}, nil
}

func (dml *DeckhouseController) Start(ec chan events.ModuleEvent) error {
	err := dml.RestoreAbsentSourceModules()
	if err != nil {
		return err
	}

	err = dml.searchAndLoadDeckhouseModules()
	if err != nil {
		return err
	}

	go dml.runEventLoop(ec)

	return nil
}

func (dml *DeckhouseController) runEventLoop(ec chan events.ModuleEvent) {
	for event := range ec {
		// event without module name
		if event.EventType == events.FirstConvergeDone {
			err := dml.handleConvergeDone()
			if err != nil {
				log.Errorf("Error occurred during the converge done: %s", err)
			}
			continue
		}

		fmt.Println("GOT EVENT WITH MODULE", event.ModuleName, event.EventType)

		mod, ok := dml.deckhouseModules[event.ModuleName]
		if !ok {
			log.Errorf("Module %q registered but not found in Deckhouse. Possible bug?", event.ModuleName)
			continue
		}
		switch event.EventType {
		case events.ModuleRegistered:
			err := dml.handleModuleRegistration(mod)
			if err != nil {
				log.Errorf("Error occurred during the module %q registration: %s", mod.GetBasicModule().GetName(), err)
				continue
			}

		case events.ModuleEnabled:
			err := dml.handleEnabledModule(mod, true)
			if err != nil {
				log.Errorf("Error occurred during the module %q turning on: %s", mod.GetBasicModule().GetName(), err)
				continue
			}

		case events.ModuleDisabled:
			err := dml.handleEnabledModule(mod, false)
			if err != nil {
				log.Errorf("Error occurred during the module %q turning off: %s", mod.GetBasicModule().GetName(), err)
				continue
			}
		}
	}
}

// handleConvergeDone after converge we delete all absent Modules CR, which were not filled during this operator startup
func (dml *DeckhouseController) handleConvergeDone() error {
	epochLabelStr := fmt.Sprintf("%s!=%s", epochLabelKey, epochLabelValue)
	return retry.OnError(retry.DefaultRetry, errors.IsServiceUnavailable, func() error {
		return dml.kubeClient.DeckhouseV1alpha1().Modules().DeleteCollection(dml.ctx, v1.DeleteOptions{}, v1.ListOptions{LabelSelector: epochLabelStr})
	})
}

func (dml *DeckhouseController) handleModulePurge(m *models.DeckhouseModule) error {
	return retry.OnError(retry.DefaultRetry, errors.IsServiceUnavailable, func() error {
		return dml.kubeClient.DeckhouseV1alpha1().Modules().Delete(dml.ctx, m.GetBasicModule().GetName(), v1.DeleteOptions{})
	})
}

func (dml *DeckhouseController) handleModuleRegistration(m *models.DeckhouseModule) error {
	return retry.OnError(retry.DefaultRetry, errors.IsServiceUnavailable, func() error {
		source := dml.sourceModules[m.GetBasicModule().GetName()]
		newModule := m.AsKubeObject(source)
		newModule.SetLabels(map[string]string{epochLabelKey: epochLabelValue})

		existModule, err := dml.kubeClient.DeckhouseV1alpha1().Modules().Get(dml.ctx, newModule.GetName(), v1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				_, err = dml.kubeClient.DeckhouseV1alpha1().Modules().Create(dml.ctx, newModule, v1.CreateOptions{})
				return err
			}

			return err
		}

		existModule.Properties = newModule.Properties
		if len(existModule.Labels) == 0 {
			newModule.SetLabels(map[string]string{epochLabelKey: epochLabelValue})
		} else {
			existModule.Labels[epochLabelKey] = epochLabelValue
		}

		_, err = dml.kubeClient.DeckhouseV1alpha1().Modules().Update(dml.ctx, existModule, v1.UpdateOptions{})

		return err
	})
}

func (dml *DeckhouseController) handleEnabledModule(m *models.DeckhouseModule, enable bool) error {
	return retry.OnError(retry.DefaultRetry, errors.IsServiceUnavailable, func() error {
		obj, err := dml.kubeClient.DeckhouseV1alpha1().Modules().Get(dml.ctx, m.GetBasicModule().GetName(), v1.GetOptions{})
		if err != nil {
			return err
		}

		obj.Properties.State = "Disabled"
		if enable {
			obj.Properties.State = "Enabled"
		}

		_, err = dml.kubeClient.DeckhouseV1alpha1().Modules().Update(dml.ctx, obj, v1.UpdateOptions{})
		if err != nil {
			return err
		}

		// Update ModuleConfig if exists
		mc, err := dml.kubeClient.DeckhouseV1alpha1().ModuleConfigs().Get(dml.ctx, m.GetBasicModule().GetName(), v1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}

			return err
		}

		mc.Status.Status = "Disabled"
		if enable {
			mc.Status.Status = "Enabled"
		}

		_, err = dml.kubeClient.DeckhouseV1alpha1().ModuleConfigs().Update(dml.ctx, mc, v1.UpdateOptions{})

		return err
	})
}
