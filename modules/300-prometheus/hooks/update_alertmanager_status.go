/*
Copyright 2022 Flant JSC

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
	"fmt"

	"github.com/clarketm/json"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"

	"github.com/deckhouse/deckhouse/go_lib/hooks/set_cr_statuses"
)

// hook for setting CR statuses
var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnAfterHelm: &go_hook.OrderedConfig{Order: 10},
}, updateAmStatus)

func updateAmStatus(input *go_hook.HookInput) error {
	addressDeclaredAlertmanagers := make([]alertmanagerAddress, 0)
	serviceDeclaredAlertmanagers := make([]alertmanagerService, 0)
	internalDeclaredAlertmanagers := make([]alertmanagerInternal, 0)

	// get AMs' names
	err := json.Unmarshal([]byte(input.Values.Get("prometheus.internal.alertmanagers.byAddress").String()), &addressDeclaredAlertmanagers)
	if err != nil {
		return fmt.Errorf("cannot unmarshal values: %v", err)
	}

	err = json.Unmarshal([]byte(input.Values.Get("prometheus.internal.alertmanagers.byService").String()), &serviceDeclaredAlertmanagers)
	if err != nil {
		return fmt.Errorf("cannot unmarshal values: %v", err)
	}

	err = json.Unmarshal([]byte(input.Values.Get("prometheus.internal.alertmanagers.internal").String()), &internalDeclaredAlertmanagers)
	if err != nil {
		return fmt.Errorf("cannot unmarshal values: %v", err)
	}

	// update AMs' statuses
	for _, am := range addressDeclaredAlertmanagers {
		input.StatusCollector.UpdateStatus(set_cr_statuses.SetProcessedStatus, "deckhouse.io/v1alpha1", "customalertmanager", "", am.Name)
	}

	for _, am := range serviceDeclaredAlertmanagers {
		input.StatusCollector.UpdateStatus(set_cr_statuses.SetProcessedStatus, "deckhouse.io/v1alpha1", "customalertmanager", "", am.ResourceName)
	}

	for _, am := range internalDeclaredAlertmanagers {
		name := am["name"].(string)
		input.StatusCollector.UpdateStatus(set_cr_statuses.SetProcessedStatus, "deckhouse.io/v1alpha1", "customalertmanager", "", name)
	}

	return nil
}
