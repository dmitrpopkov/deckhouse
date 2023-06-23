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
	"time"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
	Queue:        "/modules/prometheus/set_lookback_delta",
}, setLookbackDelta)

const minLookbackDelta = 10 * time.Second

func setLookbackDelta(input *go_hook.HookInput) error {
	var lookBackDelta time.Duration

	scrapeInterval, err := time.ParseDuration(input.Values.Get("prometheus.scrapeInterval").String())
	if err != nil {
		return nil
	}
	if scrapeInterval*2 < minLookbackDelta {
		lookBackDelta = minLookbackDelta
	} else {
		lookBackDelta = scrapeInterval * 2
	}
	input.Values.Set("prometheus.internal.prometheusMain.lookbackDelta", lookBackDelta.String())
	return nil
}
