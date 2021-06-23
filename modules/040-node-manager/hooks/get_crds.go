package hooks

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
	cljson "github.com/clarketm/json"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/deckhouse/deckhouse/modules/040-node-manager/hooks/internal/v1alpha2"
)

const (
	CRITypeDocker           = "Docker"
	CRITypeContainerd       = "Containerd"
	NodeGroupDefaultCRIType = CRITypeContainerd
)

type InstanceClassCrdInfo struct {
	Name string
	Spec interface{}
}

func applyInstanceClassCrdFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	return InstanceClassCrdInfo{
		Name: obj.GetName(),
		Spec: obj.Object["spec"],
	}, nil
}

type NodeGroupCrdInfo struct {
	Name            string
	Spec            v1alpha2.NodeGroupSpec
	ManualRolloutID string
}

// applyNodeGroupCrdFilter returns name, spec and manualRolloutID from the NodeGroup
func applyNodeGroupCrdFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	nodeGroup := new(v1alpha2.NodeGroup)
	err := sdk.FromUnstructured(obj, nodeGroup)
	if err != nil {
		return nil, err
	}

	return NodeGroupCrdInfo{
		Name:            nodeGroup.GetName(),
		Spec:            nodeGroup.Spec,
		ManualRolloutID: nodeGroup.GetAnnotations()["manual-rollout-id"],
	}, nil
}

type MachineDeploymentCrdInfo struct {
	Name string
	Zone string
}

func applyMachineDeploymentCrdFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	return MachineDeploymentCrdInfo{
		Name: obj.GetName(),
		Zone: obj.GetAnnotations()["zone"],
	}, nil
}

func applyCloudProviderSecretKindZonesFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	secretData, err := DecodeDataFromSecret(obj)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"instanceClassKind": secretData["instanceClassKind"],
		"zones":             secretData["zones"],
	}, nil
}

func applyCloudInstanceSecretKindFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	secretData, err := DecodeDataFromSecret(obj)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"instanceClassKind": secretData["instanceClassKind"],
	}, nil
}

var getCRDsHookConfig = &go_hook.HookConfig{
	Queue:        "/modules/node-manager",
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
	Kubernetes: []go_hook.KubernetesConfig{
		// A binding with dynamic kind has index 0 for simplicity.
		{
			Name:       "ics",
			ApiVersion: "",
			Kind:       "",
			FilterFunc: applyInstanceClassCrdFilter,
		},
		{
			Name:       "ngs",
			ApiVersion: "deckhouse.io/v1alpha2",
			Kind:       "NodeGroup",
			FilterFunc: applyNodeGroupCrdFilter,
		},
		{
			Name:       "machine_deployments",
			ApiVersion: "machine.sapcloud.io/v1alpha1",
			Kind:       "MachineDeployment",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-cloud-instance-manager"},
				},
			},
			LabelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "heritage",
						Operator: "In",
						Values:   []string{"deckhouse"},
					},
				},
			},
			FilterFunc: applyMachineDeploymentCrdFilter,
		},
		// kube-system/Secret/d8-node-manager-cloud-provider
		{
			Name:       "cloud_provider_secret",
			ApiVersion: "v1",
			Kind:       "Secret",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"kube-system"},
				},
			},
			NameSelector: &types.NameSelector{
				MatchNames: []string{"d8-node-manager-cloud-provider"},
			},
			FilterFunc: applyCloudProviderSecretKindZonesFilter,
		},
	},
	Schedule: []go_hook.ScheduleConfig{
		{
			Name:    "sync",
			Crontab: "*/10 * * * *",
		},
	},
}

var _ = sdk.RegisterFunc(getCRDsHookConfig, getCRDsHandler)

func getCRDsHandler(input *go_hook.HookInput) error {
	// Detect InstanceClass kind and change binding if needed.
	kindInUse, kindFromSecret := detectInstanceClassKind(input, getCRDsHookConfig)

	// Kind is changed, so objects in "dynamic-kind" can be ignored. Update kind and stop the hook.
	if kindInUse != kindFromSecret {
		if kindFromSecret == "" {
			input.LogEntry.Infof("InstanceClassKind has changed from '%s' to '': disable binding 'ics'", kindInUse)
			*input.BindingActions = append(*input.BindingActions, go_hook.BindingAction{
				Name:       "ics",
				Action:     "Disable",
				Kind:       "",
				ApiVersion: "",
			})
		} else {
			input.LogEntry.Infof("InstanceClassKind has changed from '%s' to '%s': update kind for binding 'ics'", kindInUse, kindFromSecret)
			*input.BindingActions = append(*input.BindingActions, go_hook.BindingAction{
				Name:   "ics",
				Action: "UpdateKind",
				Kind:   kindFromSecret,
				// TODO Set apiVersion to exact value? Should it be in a Secret?
				// ApiVersion: "deckhouse.io/v1alpha1",
				ApiVersion: "",
			})
		}
		// Save new kind as current kind.
		getCRDsHookConfig.Kubernetes[0].Kind = kindFromSecret
		// Binding changed, hook will be restarted with new objects in "ics" snapshot.
		return nil
	}

	// TODO What should we do with a broken semver?
	// Read kubernetes version either from clusterConfiguration or from discovery.
	var globalTargetKubernetesVersion *semver.Version
	var err error

	versionValue, has := input.Values.GetOk("global.discovery.kubernetesVersion")
	if has {
		globalTargetKubernetesVersion, err = semver.NewVersion(versionValue.String())
		if err != nil {
			return fmt.Errorf("global.discovery.kubernetesVersion contains a malformed semver: %s: %v", versionValue.String(), err)
		}
	}

	versionValue, has = input.Values.GetOk("global.clusterConfiguration.kubernetesVersion")
	if has {
		globalTargetKubernetesVersion, err = semver.NewVersion(versionValue.String())
		if err != nil {
			return fmt.Errorf("global.clusterConfiguration.kubernetesVersion contains a malformed semver: %s: %v", versionValue.String(), err)
		}
	}

	controlPlaneKubeVersions := make([]*semver.Version, 0)
	if input.Values.Exists("global.discovery.kubernetesVersions") {
		for _, verItem := range input.Values.Get("global.discovery.kubernetesVersions").Array() {
			ver, _ := semver.NewVersion(verItem.String())
			controlPlaneKubeVersions = append(controlPlaneKubeVersions, ver)
		}
	}

	controlPlaneMinVersion := SemverMin(controlPlaneKubeVersions)

	// Default zones. Take them from input.Snapshots["machine_deployments"]
	// and from input.Snapshots["cloud_provider_secret"].zones
	defaultZones := make([]string, 0)
	defaultZonesMap := make(map[string]struct{})
	for _, machineInfoItem := range input.Snapshots["machine_deployments"] {
		machineInfo := machineInfoItem.(MachineDeploymentCrdInfo)
		if _, has := defaultZonesMap[machineInfo.Zone]; !has {
			defaultZonesMap[machineInfo.Zone] = struct{}{}
			defaultZones = append(defaultZones, machineInfo.Zone)
		}
	}
	if len(input.Snapshots["cloud_provider_secret"]) > 0 {
		secretInfo := input.Snapshots["cloud_provider_secret"][0].(map[string]interface{})
		zonesUntyped := secretInfo["zones"]

		input.LogEntry.Infof("cloud_provider_secret.zones: %T %#v", zonesUntyped, zonesUntyped)

		zonesTyped := make([]string, 0)

		switch v := zonesUntyped.(type) {
		case []string:
			zonesTyped = append(zonesTyped, v...)
		case []interface{}:
			for _, zoneUntyped := range v {
				if s, ok := zoneUntyped.(string); ok {
					zonesTyped = append(zonesTyped, s)
				}
			}
		case string:
			zonesTyped = append(zonesTyped, v)
		}

		for _, zoneStr := range zonesTyped {
			if _, has := defaultZonesMap[zoneStr]; !has {
				defaultZonesMap[zoneStr] = struct{}{}
				defaultZones = append(defaultZones, zoneStr)
			}
		}
	}

	// Save timestamp for updateEpoch.
	timestamp := epochTimestampAccessor()

	finalNodeGroups := make([]interface{}, 0)

	for _, v := range input.Snapshots["ngs"] {
		nodeGroup := v.(NodeGroupCrdInfo)
		ngForValues := nodeGroupForValues(nodeGroup.Spec.DeepCopy())
		// Copy manualRolloutID and name.
		ngForValues["name"] = nodeGroup.Name
		ngForValues["manualRolloutID"] = nodeGroup.ManualRolloutID

		if nodeGroup.Spec.NodeType == "Static" {
			if staticValue, has := input.Values.GetOk("nodeManager.internal.static"); has {
				static := staticValue.Map()
				if len(static) > 0 {
					ngForValues["static"] = static
				}
			}
		}

		if nodeGroup.Spec.NodeType == "Cloud" && kindInUse != "" {
			instanceClasses := make(map[string]interface{})

			for _, icsItem := range input.Snapshots["ics"] {
				ic := icsItem.(InstanceClassCrdInfo)
				instanceClasses[ic.Name] = ic.Spec
			}

			// check #1 — .spec.cloudInstances.classReference.kind should be allowed in our cluster
			nodeGroupInstanceClassKind := nodeGroup.Spec.CloudInstances.ClassReference.Kind
			if nodeGroupInstanceClassKind != kindInUse {
				errorMsg := fmt.Sprintf("Wrong classReference: Kind %s is not allowed, the only allowed kind is %s.", nodeGroupInstanceClassKind, kindInUse)

				if input.Values.Exists("nodeManager.internal.nodeGroups") {
					savedNodeGroups := input.Values.Get("nodeManager.internal.nodeGroups").Array()
					for _, savedNodeGroup := range savedNodeGroups {
						ng := savedNodeGroup.Map()
						if ng["name"].String() == nodeGroup.Name {
							finalNodeGroups = append(finalNodeGroups, savedNodeGroup.Value().(map[string]interface{}))
							errorMsg += " Earlier stored version of NG is in use now!"
						}
					}
				}

				input.LogEntry.Errorf("Bad NodeGroup '%s': %s", nodeGroup.Name, errorMsg)
				setNodeGroupErrorStatus(input.ObjectPatcher, nodeGroup.Name, errorMsg)
				continue
			}

			// check #2 — .spec.cloudInstances.classReference should be valid
			nodeGroupInstanceClassName := nodeGroup.Spec.CloudInstances.ClassReference.Name
			isKnownClassName := false
			for className := range instanceClasses {
				if className == nodeGroupInstanceClassName {
					isKnownClassName = true
					break
				}
			}
			if !isKnownClassName {
				errorMsg := fmt.Sprintf("Wrong classReference: There is no valid instance class %s of type %s.", nodeGroupInstanceClassName, nodeGroupInstanceClassKind)

				if input.Values.Exists("nodeManager.internal.nodeGroups") {
					savedNodeGroups := input.Values.Get("nodeManager.internal.nodeGroups").Array()
					for _, savedNodeGroup := range savedNodeGroups {
						ng := savedNodeGroup.Map()
						if ng["name"].String() == nodeGroup.Name {
							finalNodeGroups = append(finalNodeGroups, savedNodeGroup.Value().(map[string]interface{}))
							errorMsg += " Earlier stored version of NG is in use now!"
						}
					}
				}

				input.LogEntry.Errorf("Bad NodeGroup '%s': %s", nodeGroup.Name, errorMsg)
				setNodeGroupErrorStatus(input.ObjectPatcher, nodeGroup.Name, errorMsg)
				continue
			}

			// check #3 — zones should be valid
			if len(defaultZonesMap) > 0 {
				// All elements in nodeGroup.Spec.CloudInstances.Zones
				// should contain in defaultZonesMap.
				containCount := 0
				unknownZones := make([]string, 0)
				for _, zone := range nodeGroup.Spec.CloudInstances.Zones {
					if _, has := defaultZonesMap[zone]; has {
						containCount++
					} else {
						unknownZones = append(unknownZones, zone)
					}
				}
				if containCount != len(nodeGroup.Spec.CloudInstances.Zones) {
					errorMsg := fmt.Sprintf("unknown cloudInstances.zones: %v", unknownZones)
					input.LogEntry.Errorf("Bad NodeGroup '%s': %s", nodeGroup.Name, errorMsg)

					setNodeGroupErrorStatus(input.ObjectPatcher, nodeGroup.Name, errorMsg)
					continue
				}
			}

			// Put instanceClass.spec into values.
			ngForValues["instanceClass"] = instanceClasses[nodeGroupInstanceClassName]

			var zones []string
			if nodeGroup.Spec.CloudInstances.Zones != nil {
				zones = nodeGroup.Spec.CloudInstances.Zones
			}
			if zones == nil {
				zones = defaultZones
			}

			if ngForValues["cloudInstances"] == nil {
				ngForValues["cloudInstances"] = v1alpha2.CloudInstances{}
			}
			cloudInstances := ngForValues["cloudInstances"].(v1alpha2.CloudInstances)
			cloudInstances.Zones = zones
			ngForValues["cloudInstances"] = cloudInstances
		}

		// Determine effective Kubernetes version.
		effectiveKubeVer := globalTargetKubernetesVersion
		if controlPlaneMinVersion != nil {
			if effectiveKubeVer == nil || effectiveKubeVer.GreaterThan(controlPlaneMinVersion) {
				// Nodes should not be above control plane
				effectiveKubeVer = controlPlaneMinVersion
			}
		}
		ngForValues["kubernetesVersion"] = SemverMajMin(effectiveKubeVer)

		// Detect CRI type. Default CRI type is 'Docker' for Kubernetes version less than 1.19.
		v1_19_0, _ := semver.NewVersion("1.19.0")
		defaultCRIType := NodeGroupDefaultCRIType
		if effectiveKubeVer.LessThan(v1_19_0) {
			defaultCRIType = CRITypeDocker
		}

		if criValue, has := input.Values.GetOk("global.clusterConfiguration.defaultCRI"); has {
			defaultCRIType = criValue.String()
		}

		newCRIType := nodeGroup.Spec.CRI.Type
		if newCRIType == "" {
			newCRIType = defaultCRIType
		}

		switch newCRIType {
		case CRITypeDocker:
			// cri is NotManaged if .spec.cri.docker.manage is explicitly set to false.
			if nodeGroup.Spec.CRI.Docker != nil && nodeGroup.Spec.CRI.Docker.Manage != nil && !*nodeGroup.Spec.CRI.Docker.Manage {
				newCRIType = "NotManaged"
			}
		case CRITypeContainerd:
			// Containerd requires Kubernetes version 1.19+.
			if effectiveKubeVer.LessThan(v1_19_0) {
				return fmt.Errorf("cri type Containerd is allowed only for kubernetes 1.19+")
			}
		}

		if ngForValues["cri"] == nil {
			ngForValues["cri"] = v1alpha2.CRI{}
		}
		cri := ngForValues["cri"].(v1alpha2.CRI)
		cri.Type = newCRIType
		ngForValues["cri"] = cri

		// Calculate update epoch
		// updateEpoch is a value that changes every 4 hour for a particular NodeGroup in the cluster.
		// Values are spread over 4 hour window to update nodes at different times.
		// Also, updateEpoch value is a unix time of the next update.
		updateEpoch := calculateUpdateEpoch(timestamp,
			input.Values.Get("global.discovery.clusterUUID").String(),
			nodeGroup.Name)
		ngForValues["updateEpoch"] = updateEpoch

		// Reset status error for current NodeGroup.
		setNodeGroupErrorStatus(input.ObjectPatcher, nodeGroup.Name, "")

		ngBytes, _ := cljson.Marshal(ngForValues)
		finalNodeGroups = append(finalNodeGroups, json.RawMessage(ngBytes))
	}

	input.Values.Set("nodeManager.internal.nodeGroups", finalNodeGroups)
	return nil
}

func nodeGroupForValues(nodeGroupSpec *v1alpha2.NodeGroupSpec) map[string]interface{} {
	res := make(map[string]interface{})

	res["nodeType"] = nodeGroupSpec.NodeType
	if !nodeGroupSpec.CRI.IsEmpty() {
		res["cri"] = nodeGroupSpec.CRI
	}
	if !nodeGroupSpec.CloudInstances.IsEmpty() {
		res["cloudInstances"] = nodeGroupSpec.CloudInstances
	}
	if !nodeGroupSpec.NodeTemplate.IsEmpty() {
		res["nodeTemplate"] = nodeGroupSpec.NodeTemplate
	}
	if !nodeGroupSpec.Chaos.IsEmpty() {
		res["chaos"] = nodeGroupSpec.Chaos
	}
	if !nodeGroupSpec.OperatingSystem.IsEmpty() {
		res["operatingSystem"] = nodeGroupSpec.OperatingSystem
	}
	if !nodeGroupSpec.Disruptions.IsEmpty() {
		res["disruptions"] = nodeGroupSpec.Disruptions
	}
	if !nodeGroupSpec.Kubelet.IsEmpty() {
		res["kubelet"] = nodeGroupSpec.Kubelet
	}

	return res
}

func setNodeGroupErrorStatus(patcher go_hook.ObjectPatcher, nodeGroupName, message string) {
	statusErrorPatch, _ := json.Marshal(map[string]interface{}{
		"status": map[string]interface{}{
			"error": message,
		},
	})
	// Patches are deferred, error is always nil.
	_ = patcher.MergePatchObject(statusErrorPatch, "deckhouse.io/v1alpha2", "NodeGroup", "", nodeGroupName, "/status")
}

var epochTimestampAccessor = func() int64 {
	return time.Now().Unix()
}

var detectInstanceClassKind = func(input *go_hook.HookInput, config *go_hook.HookConfig) (inUse string, fromSecret string) {
	if len(input.Snapshots["cloud_provider_secret"]) > 0 {
		if secretInfo, ok := input.Snapshots["cloud_provider_secret"][0].(map[string]interface{}); ok {
			if kind, ok := secretInfo["instanceClassKind"].(string); ok {
				fromSecret = kind
			}
		}
	}

	return getCRDsHookConfig.Kubernetes[0].Kind, fromSecret
}

const EpochWindowSize int64 = 4 * 60 * 60 // 4 hours
// calculateUpdateEpoch returns an end point of the drifted 4 hour window for given cluster and timestamp.
//
// epoch is the unix timestamp of an end time of the drifted 4 hour window.
//   A0---D0---------------------A1---D1---------------------A2---D2-----
//   A - points for windows in absolute time
//   D - points for drifted windows
//
// Epoch for timestamps A0 <= ts <= D0 is D0
//
// Epoch for timestamps D0 < ts <= D1 is D1
func calculateUpdateEpoch(ts int64, clusterUUID string, nodeGroupName string) string {
	hasher := fnv.New64a()
	// error is always nil here
	_, _ = hasher.Write([]byte(clusterUUID))
	_, _ = hasher.Write([]byte(nodeGroupName))
	drift := int64(hasher.Sum64() % uint64(EpochWindowSize))

	// Near zero timestamps. It should not happen, isn't it?
	if ts <= drift {
		return strconv.FormatInt(drift, 10)
	}

	// Get the start of the absolute time window (non-drifted). Correct timestamp be 1 second
	// to get correct window start when timestamp is equal to the end of the drifted window.
	absWindowStart := ((ts - drift - 1) / EpochWindowSize) * EpochWindowSize
	epoch := absWindowStart + EpochWindowSize + drift
	return strconv.FormatInt(epoch, 10)
}
