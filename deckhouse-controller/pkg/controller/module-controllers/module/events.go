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
	d8config "github.com/deckhouse/deckhouse/go_lib/deckhouse-config"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"

	"github.com/flant/addon-operator/pkg/module_manager/models/modules/events"
)

func (r *reconciler) runModuleEventLoop(ctx context.Context) error {
	for event := range r.moduleManager.GetModuleEventsChannel() {
		switch event.EventType {
		case events.ModuleRegistered:
			// TODO(ipaqsa) why?
			// add module name as a possible name for validation module config webhook
			d8config.Service().AddPossibleName(event.ModuleName)
			continue
		case events.ModuleConfigChanged:
			if err := r.refreshModuleConfigAndModule(ctx, event.ModuleName); err != nil {
				r.log.Error(err, "failed to handle ModuleConfigChanged event: failed to refresh module config and module", "module", event.ModuleName)
			}
			continue
		case events.ModuleEnabled:
			if err := r.handleEnabledDisabledEvent(ctx, event.ModuleName, true); err != nil {
				r.log.Error(err, "failed to handle ModuleEnabled event: failed to enable module", "module", event.ModuleName)
			}
			continue
		case events.ModuleDisabled:
			if err := r.handleEnabledDisabledEvent(ctx, event.ModuleName, false); err != nil {
				r.log.Error(err, "failed to handle ModuleDisabled event: failed to disable module", "module", event.ModuleName)
			}
			continue
		case events.ModuleStateChanged:
			if err := r.refreshModuleByModuleConfig(ctx, event.ModuleName); err != nil {
				r.log.Error(err, "failed to handle ModuleStateChanged event: failed to refresh module", "module", event.ModuleName)
			}
			continue
		default:
			r.log.Warn("unknown module event type", "type", event.EventType, "module", event.ModuleName)
			continue
		}
	}
	return nil
}

func (r *reconciler) handleEnabledDisabledEvent(ctx context.Context, moduleName string, enable bool) error {
	module := new(v1alpha1.Module)
	if err := r.client.Get(ctx, client.ObjectKey{Name: moduleName}, module); err != nil {
		return err
	}
	return r.setModuleEnabled(ctx, module, enable)
}
