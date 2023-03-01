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
	"sort"
	"time"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/deckhouse/deckhouse/go_lib/dependency"
)

// We have to block main queue until Deckhouse release will be deployed
// otherwise helm release could stuck in the pending-upgrade state and deckhouse will be broken

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnAfterHelm: &go_hook.OrderedConfig{Order: 10},
}, dependency.WithExternalDependencies(waitForDeckhouseReleaseToBeDeployed))

var (
	pollTimeout = 30 * time.Second
)

func waitForDeckhouseReleaseToBeDeployed(input *go_hook.HookInput, dc dependency.Container) error {
	var lastErr error
	err := wait.Poll(time.Second, pollTimeout, func() (done bool, err error) {
		input.LogEntry.Infof("waiting for deckhouse release to be deployed")
		ok, err := isReleaseDeployed(input, dc)
		if err != nil {
			lastErr = err
			return false, err
		}

		return ok, nil
	})

	if err != nil {
		return fmt.Errorf("timeout waiting for deckhouse release to be deployed. last error: %v", lastErr)
	}

	return nil
}

func isReleaseDeployed(input *go_hook.HookInput, dc dependency.Container) (bool, error) {
	kubeClient, err := dc.GetK8sClient()
	if err != nil {
		input.LogEntry.Errorf("%v", err)
		return false, err
	}

	releases, err := kubeClient.CoreV1().Secrets("d8-system").List(context.TODO(), v1.ListOptions{LabelSelector: "name=deckhouse"})
	if err != nil {
		input.LogEntry.Errorf("%v", err)
		return false, err
	}

	sort.Slice(releases.Items, func(i, j int) bool {
		// reversed rule to have the latest release on the [0] place
		return releases.Items[i].CreationTimestamp.After(releases.Items[j].CreationTimestamp.Time)
	})

	if len(releases.Items) == 0 {
		return true, nil
	}

	latestRelease := releases.Items[0]

	if latestRelease.Labels["status"] == "deployed" {
		return true, nil
	}

	return false, nil
}
