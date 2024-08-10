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

package hooks

import (
	"context"
	"fmt"
	"slices"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/pointer"

	"github.com/deckhouse/deckhouse/dhctl/pkg/config"
	"github.com/deckhouse/deckhouse/go_lib/dependency"
	"github.com/deckhouse/deckhouse/go_lib/dependency/k8s"
)

/* Migration:
This migration implements global hook which migrate cni settings from kube-system/d8-cni-configuration secret to appropriate module config.
If secret doesn't exist, migration skipped.
If module config for cni exists, migration skipped.
Migration scheme:
* cni-simple-bridge
	Simply enable module cni-simple-bridge via ModuleConfig.
* cni-flannel
* cni-cilium
*/

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnStartup: &go_hook.OrderedConfig{Order: 20},
}, dependency.WithExternalDependencies(d8cniSecretMigrate))

func d8cniSecretMigrate(input *go_hook.HookInput, dc dependency.Container) error {
	kubeCl, err := dc.GetK8sClient()
	if err != nil {
		return fmt.Errorf("cannot init Kubernetes client: %v", err)
	}

	// skip migration if d8-cni-configuration secret doesn't exist.
	d8cniSecret, err := kubeCl.CoreV1().Secrets("kube-system").Get(context.TODO(), "d8-cni-configuration", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		input.LogEntry.Info("d8-cni-configuration secret does not exist, skipping migration")
		return nil
	}
	if err != nil {
		return err
	}

	moduleConfigs, err := kubeCl.Dynamic().Resource(config.ModuleConfigGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	for _, mc := range moduleConfigs.Items {

		if slices.Contains([]string{"cni-cilium", "cni-flannel", "cni-simple-bridge"}, mc.GetName()) {
			moduleEnabled, exists, err := unstructured.NestedBool(mc.UnstructuredContent(), "spec", "enabled")
			if err != nil {
				return err
			}
			if !exists {
				break
			}
			if !moduleEnabled {
				break
			}
			input.LogEntry.Infof("Module config for %s found, skipping migration", mc.GetName())
			return removeD8CniSecret(input, kubeCl, d8cniSecret)
		}
	}

	cniBytes, ok := d8cniSecret.Data["cni"]
	if !ok {
		return fmt.Errorf("d8-cni-configuration secret does not contain \"cni\" field, skipping migration")
	}

	cniName := string(cniBytes)
	cniModuleConfig := &config.ModuleConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       config.ModuleConfigKind,
			APIVersion: config.ModuleConfigGroup + "/" + config.ModuleConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: cniName,
		},
	}

	switch cniName {
	case "cni-simple-bridge":
		cniModuleConfig.Spec = config.ModuleConfigSpec{
			Enabled: pointer.Bool(true),
		}
	default:
		return fmt.Errorf("unknown cni name: %s", cniName)
	}

	// create secret
	err = createMC(input, kubeCl, cniModuleConfig)
	if err != nil {
		return err
	}

	//remove secret
	return removeD8CniSecret(input, kubeCl, d8cniSecret)
}

// remove secret
func removeD8CniSecret(input *go_hook.HookInput, kubeCl k8s.Client, secret *v1.Secret) error {
	var secretData []byte
	err := secret.Unmarshal(secretData)
	if err != nil {
		return err
	}
	input.LogEntry.Info(string(secretData))
	return kubeCl.CoreV1().Secrets(secret.Namespace).Delete(context.TODO(), secret.Name, metav1.DeleteOptions{})
}

// create Module Config
func createMC(input *go_hook.HookInput, kubeCl k8s.Client, mc *config.ModuleConfig) error {
	obj, err := sdk.ToUnstructured(mc)
	if err != nil {
		return err
	}
	_, err = kubeCl.Dynamic().Resource(config.ModuleConfigGVR).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}
