package controller

import (
	"context"

	"github.com/flant/addon-operator/pkg/module_manager"

	log "github.com/sirupsen/logrus"

	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/pkg/values/validation"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/client/clientset/versioned"
)

type DeckhouseController struct {
	ctx context.Context

	dirs            []string
	valuesValidator *validation.ValuesValidator
	kubeClient      *versioned.Clientset

	deckhouseModules map[string]*DeckhouseModule
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

		deckhouseModules: make(map[string]*DeckhouseModule),
	}, nil
}

func (dml *DeckhouseController) Start(ec chan module_manager.ModuleEvent) error {
	err := dml.searchAndLoadDeckhouseModules()
	if err != nil {
		return err
	}

	go dml.runEventLoop(ec)

	return nil
}

func (dml *DeckhouseController) runEventLoop(ec chan module_manager.ModuleEvent) {
	for event := range ec {
		mod, ok := dml.deckhouseModules[event.ModuleName]
		if !ok {
			log.Errorf("Module %q registered but not found in Deckhouse. Possible bug?", mod.basic.GetName())
			continue
		}
		switch event.EventType {
		case module_manager.ModuleRegistered:

			err := dml.handleModuleRegistration(mod)
			if err != nil {
				log.Errorf("Error occurred during the module %q registration: %s", mod.basic.GetName(), err)
				continue
			}

		case module_manager.ModuleEnabled:
			err := dml.handleEnabledModule(mod, true)
			if err != nil {
				log.Errorf("Error occurred during the module %q turning on: %s", mod.basic.GetName(), err)
				continue
			}

		case module_manager.ModuleDisabled:
			err := dml.handleEnabledModule(mod, false)
			if err != nil {
				log.Errorf("Error occurred during the module %q turning off: %s", mod.basic.GetName(), err)
				continue
			}

		}

	}
}

func (dml *DeckhouseController) handleModuleRegistration(m *DeckhouseModule) error {
	existModule, err := dml.kubeClient.DeckhouseV1alpha1().Modules().Get(dml.ctx, m.basic.GetName(), v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = dml.kubeClient.DeckhouseV1alpha1().Modules().Create(dml.ctx, m.AsKubeObject(), v1.CreateOptions{})
			return err
		}

		return err
	}

	existModule.Properties = m.AsKubeObject().Properties

	_, err = dml.kubeClient.DeckhouseV1alpha1().Modules().Update(dml.ctx, existModule, v1.UpdateOptions{})

	return err
}

func (dml *DeckhouseController) handleEnabledModule(m *DeckhouseModule, enable bool) error {
	obj, err := dml.kubeClient.DeckhouseV1alpha1().Modules().Get(dml.ctx, m.basic.GetName(), v1.GetOptions{})
	if err != nil {
		return err
	}
	mc, err := dml.kubeClient.DeckhouseV1alpha1().ModuleConfigs().Get(dml.ctx, m.basic.GetName(), v1.GetOptions{})
	if err != nil {
		return err
	}

	if enable {
		obj.Properties.State = "Enabled"
		mc.Status.State = "Enabled"

	} else {
		obj.Properties.State = "Disabled"
		mc.Status.Status = "Disabled"
	}

	_, err = dml.kubeClient.DeckhouseV1alpha1().Modules().Update(dml.ctx, obj, v1.UpdateOptions{})
	if err != nil {
		return err
	}

	_, err = dml.kubeClient.DeckhouseV1alpha1().ModuleConfigs().Update(dml.ctx, mc, v1.UpdateOptions{})

	return err
}
