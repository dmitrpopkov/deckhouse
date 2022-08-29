/*
Copyright 2021 Flant JSC

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
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	updater2 "github.com/deckhouse/deckhouse/modules/020-deckhouse/hooks/internal/update"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook/metrics"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/pointer"

	"github.com/deckhouse/deckhouse/go_lib/dependency"
	"github.com/deckhouse/deckhouse/go_lib/dependency/cr"
	"github.com/deckhouse/deckhouse/go_lib/hooks/update"
	"github.com/deckhouse/deckhouse/modules/002-deckhouse/hooks/internal/apis/v1alpha1"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: "/modules/deckhouse/update_deckhouse_image",
	Schedule: []go_hook.ScheduleConfig{
		{
			Name:    "update_deckhouse_image",
			Crontab: "*/15 * * * * *",
		},
	},
	Settings: &go_hook.HookConfigSettings{
		EnableSchedulesOnStartup: true,
	},
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "deckhouse_pod",
			ApiVersion: "v1",
			Kind:       "Pod",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-system"},
				},
			},
			LabelSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "deckhouse",
				},
			},
			FieldSelector: &types.FieldSelector{
				MatchExpressions: []types.FieldSelectorRequirement{
					{
						Field:    "status.phase",
						Operator: "Equals",
						Value:    "Running",
					},
				},
			},
			ExecuteHookOnEvents:          pointer.BoolPtr(false),
			ExecuteHookOnSynchronization: pointer.BoolPtr(false),
			FilterFunc:                   filterDeckhousePod,
		},
		{
			Name:                         "releases",
			ApiVersion:                   "deckhouse.io/v1alpha1",
			Kind:                         "DeckhouseRelease",
			ExecuteHookOnEvents:          pointer.BoolPtr(false),
			ExecuteHookOnSynchronization: pointer.BoolPtr(false),
			FilterFunc:                   filterDeckhouseRelease,
		},
		{
			Name:       "updating_cm",
			ApiVersion: "v1",
			Kind:       "ConfigMap",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-system"},
				},
			},
			NameSelector: &types.NameSelector{
				MatchNames: []string{"d8-release-updating"},
			},
			ExecuteHookOnSynchronization: pointer.BoolPtr(false),
			ExecuteHookOnEvents:          pointer.BoolPtr(false),
			FilterFunc:                   filterUpdatingCM,
		},
	},
}, dependency.WithExternalDependencies(updateDeckhouse))

type deckhousePodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
	ImageID   string `json:"imageID"`
	Ready     bool   `json:"ready"`
}

// while cluster bootstrapping we have the tag for deckhouse image like: alpha, beta, early-access, stable, rock-solid
// it is set via dhctl, which does not know anything about releases and tags
// We can use this bootstrap image for applying first release without any requirements (like update windows, canary, etc)
func (dpi deckhousePodInfo) isBootstrapImage() bool {
	colonIndex := strings.LastIndex(dpi.Image, ":")
	if colonIndex == -1 {
		return false
	}

	tag := dpi.Image[colonIndex+1:]

	if tag == "" {
		return false
	}

	switch strings.ToLower(tag) {
	case "alpha", "beta", "early-access", "stable", "rock-solid":
		return true

	default:
		return false
	}
}

const (
	metricReleasesGroup = "d8_releases"
	metricUpdatingGroup = "d8_updating"
)

func updateDeckhouse(input *go_hook.HookInput, dc dependency.Container) error {
	deckhousePod := getDeckhousePod(input.Snapshots["deckhouse_pod"])
	if deckhousePod == nil {
		input.LogEntry.Warn("Deckhouse pod does not exist. Skipping update")
		return nil
	}

	if !input.Values.Exists("deckhouse.releaseChannel") {
		// dev upgrade - by tag
		return tagUpdate(input, dc, deckhousePod)
	}

	// production upgrade
	input.MetricsCollector.Expire(metricReleasesGroup)

	if deckhousePod.Ready {
		input.MetricsCollector.Expire(metricUpdatingGroup)
		if isUpdatingCMExists(input) {
			deleteUpdatingCM(input)
		}
	} else if isUpdatingCMExists(input) {
		input.MetricsCollector.Set("d8_is_updating", 1, nil, metrics.WithGroup(metricUpdatingGroup))
	}

	// initialize updater
	approvalMode := input.Values.Get("deckhouse.update.mode").String()
	updater := updater2.NewDeckhouseUpdater(approvalMode, deckhousePod.Ready, deckhousePod.isBootstrapImage())

	// fetch releases from snapshot and patch initial statuses
	updater.FetchAndPrepareReleases(input)
	if len(updater.releases) == 0 {
		return nil
	}

	// predict next patch for Deploy
	updater.PredictNextRelease()

	// has already Deployed the latest release
	if updater.LastReleaseDeployed() {
		return nil
	}

	// some release is forced, burn everything, apply this patch!
	if updater.HasForceRelease() {
		updater.ApplyForcedRelease(input)
		return nil
	}

	if updater.PredictedReleaseIsPatch() {
		// patch release does not respect update windows or ManualMode
		updater.ApplyPredictedRelease(input)
		return nil
	} else if !updater.inManualMode {
		// update windows works only for Auto deployment mode
		windows, err := getUpdateWindows(input)
		if err != nil {
			return fmt.Errorf("update windows configuration is not valid: %s", err)
		}

		updatePermitted := isUpdatePermitted(windows)

		if !updatePermitted {
			input.LogEntry.Info("Deckhouse update does not get into update windows. Skipping")
			release := updater.PredictedRelease()
			if release != nil {
				updateStatus(input, release, "Release is waiting for update window", v1alpha1.PhasePending)
			}
			return nil
		}
	}

	updater.ApplyPredictedRelease(input)
	return nil
}

// getUpdateWindows return set update windows
func getUpdateWindows(input *go_hook.HookInput) (update.Windows, error) {
	windowsData, exists := input.Values.GetOk("deckhouse.update.windows")
	if !exists {
		return nil, nil
	}

	return update.FromJSON([]byte(windowsData.Raw))
}

// used also in check_deckhouse_release.go
func filterDeckhouseRelease(unstructured *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var release v1alpha1.DeckhouseRelease

	err := sdk.FromUnstructured(unstructured, &release)
	if err != nil {
		return nil, err
	}

	var hasSuspendAnnotation, hasForceAnnotation, hasDisruptionApprovedAnnotation bool

	if v, ok := release.Annotations["release.deckhouse.io/suspended"]; ok {
		if v == "true" {
			hasSuspendAnnotation = true
		}
	}

	if v, ok := release.Annotations["release.deckhouse.io/force"]; ok {
		if v == "true" {
			hasForceAnnotation = true
		}
	}

	if v, ok := release.Annotations["release.deckhouse.io/disruption-approved"]; ok {
		if v == "true" {
			hasDisruptionApprovedAnnotation = true
		}
	}

	var releaseApproved bool
	if v, ok := release.Annotations["release.deckhouse.io/approved"]; ok {
		if v == "true" {
			releaseApproved = true
		}
	} else {
		releaseApproved = release.Approved
	}

	var cooldown *time.Time
	if v, ok := release.Annotations["release.deckhouse.io/cooldown"]; ok {
		cd, err := time.Parse(time.RFC3339, v)
		if err == nil {
			cooldown = &cd
		}
	}

	return updater2.DeckhouseRelease{
		Name:          release.Name,
		Version:       semver.MustParse(release.Spec.Version),
		ApplyAfter:    release.Spec.ApplyAfter,
		CooldownUntil: cooldown,
		Requirements:  release.Spec.Requirements,
		Disruptions:   release.Spec.Disruptions,
		Status: v1alpha1.DeckhouseReleaseStatus{
			Phase:    release.Status.Phase,
			Approved: release.Status.Approved,
			Message:  release.Status.Message,
		},
		ManuallyApproved:                releaseApproved,
		HasSuspendAnnotation:            hasSuspendAnnotation,
		HasForceAnnotation:              hasForceAnnotation,
		HasDisruptionApprovedAnnotation: hasDisruptionApprovedAnnotation,
	}, nil
}

func filterUpdatingCM(unstructured *unstructured.Unstructured) (go_hook.FilterResult, error) {
	return unstructured.GetName(), nil
}

func filterDeckhousePod(unstructured *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var pod corev1.Pod
	err := sdk.FromUnstructured(unstructured, &pod)
	if err != nil {
		return nil, err
	}

	// ignore evicted and shutdown pods
	if pod.Status.Phase == corev1.PodFailed {
		return nil, nil
	}

	var imageName, imageID string

	if len(pod.Spec.Containers) > 0 {
		imageName = pod.Spec.Containers[0].Image
	}

	var ready bool

	if len(pod.Status.ContainerStatuses) > 0 {
		imageID = pod.Status.ContainerStatuses[0].ImageID
		ready = pod.Status.ContainerStatuses[0].Ready
	}

	return deckhousePodInfo{
		Image:     imageName,
		ImageID:   imageID,
		Name:      pod.Name,
		Namespace: pod.Namespace,
		Ready:     ready,
	}, nil
}

func isUpdatePermitted(windows update.Windows) bool {
	if len(windows) == 0 {
		return true
	}

	now := time.Now()

	if os.Getenv("D8_IS_TESTS_ENVIRONMENT") != "" {
		now = time.Date(2021, 01, 01, 13, 30, 00, 00, time.UTC)
	}

	return windows.IsAllowed(now)
}

// tagUpdate update by tag, in dev mode or specified image
func tagUpdate(input *go_hook.HookInput, dc dependency.Container, deckhousePod *deckhousePodInfo) error {
	if deckhousePod.Image == "" && deckhousePod.ImageID == "" {
		// pod is restarting or something like that, try more in a 15 seconds
		return nil
	}

	if deckhousePod.Image == "" || deckhousePod.ImageID == "" {
		input.LogEntry.Debug("Deckhouse pod is not ready. Try to update later")
		return nil
	}

	idSplitIndex := strings.LastIndex(deckhousePod.ImageID, "@")
	if idSplitIndex == -1 {
		return fmt.Errorf("image hash not found: %s", deckhousePod.ImageID)
	}
	imageHash := deckhousePod.ImageID[idSplitIndex+1:]

	imageSplitIndex := strings.LastIndex(deckhousePod.Image, ":")
	if imageSplitIndex == -1 {
		return fmt.Errorf("image tag not found: %s", deckhousePod.Image)
	}
	repo := deckhousePod.Image[:imageSplitIndex]
	tag := deckhousePod.Image[imageSplitIndex+1:]

	regClient, err := dc.GetRegistryClient(repo, cr.WithCA(getCA(input)), cr.WithInsecureSchema(isHTTP(input)))
	if err != nil {
		input.LogEntry.Errorf("Registry (%s) client init failed: %s", repo, err)
		return nil
	}

	input.MetricsCollector.Inc("deckhouse_registry_check_total", map[string]string{})
	input.MetricsCollector.Inc("deckhouse_kube_image_digest_check_total", map[string]string{})

	repoDigest, err := regClient.Digest(tag)
	if err != nil {
		input.MetricsCollector.Inc("deckhouse_registry_check_errors_total", map[string]string{})
		input.LogEntry.Errorf("Registry (%s) get digest failed: %s", repo, err)
		return nil
	}

	input.MetricsCollector.Set("deckhouse_kube_image_digest_check_success", 1.0, map[string]string{})

	if strings.TrimSpace(repoDigest) == strings.TrimSpace(imageHash) {
		return nil
	}

	input.LogEntry.Info("New deckhouse image found. Restarting")

	input.PatchCollector.Delete("v1", "Pod", deckhousePod.Namespace, deckhousePod.Name)

	return nil
}

// Updater

func createUpdatingCM(input *go_hook.HookInput, version string) {
	cm := &corev1.ConfigMap{
		TypeMeta: v1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "d8-release-updating",
			Namespace: "d8-system",
			Labels: map[string]string{
				"heritage": "deckhouse",
			},
		},
		Data: map[string]string{
			"version": version,
		},
	}

	input.PatchCollector.Create(cm, object_patch.UpdateIfExists())
}

func isUpdatingCMExists(input *go_hook.HookInput) bool {
	snap := input.Snapshots["updating_cm"]
	return len(snap) > 0
}

func deleteUpdatingCM(input *go_hook.HookInput) {
	input.PatchCollector.Delete("v1", "ConfigMap", "d8-system", "d8-release-updating", object_patch.InBackground())
}

func getDeckhousePod(snap []go_hook.FilterResult) *deckhousePodInfo {
	var deckhousePod deckhousePodInfo

	switch len(snap) {
	case 0:
		return nil

	case 1:
		deckhousePod = snap[0].(deckhousePodInfo)

	default:
		for _, sn := range snap {
			if sn == nil {
				continue
			}
			deckhousePod = sn.(deckhousePodInfo)
			break
		}
	}

	return &deckhousePod
}
