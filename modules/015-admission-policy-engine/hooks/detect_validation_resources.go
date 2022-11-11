package hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/yaml"

	"github.com/deckhouse/deckhouse/go_lib/dependency"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "constraint-exporter-cm",
			ApiVersion: "v1",
			Kind:       "ConfigMap",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-admission-policy-engine"},
				},
			},
			NameSelector: &types.NameSelector{
				MatchNames: []string{"constraint-exporter"},
			},
			FilterFunc: filterExporterCM,
		},
	},
}, dependency.WithExternalDependencies(handleValidationKinds))

func handleValidationKinds(input *go_hook.HookInput, dc dependency.Container) error {
	snap := input.Snapshots["constraint-exporter-cm"]
	if len(snap) == 0 {
		input.LogEntry.Info("no exporter cm found")
		return nil
	}

	kindsRaw := snap[0].(string)

	var matchKinds []matchKind

	err := yaml.Unmarshal([]byte(kindsRaw), &matchKinds)
	if err != nil {
		return err
	}

	res := make([]matchResoyrce, 0, len(matchKinds))

	k8s, _ := dc.GetK8sClient()
	apiRes, err := restmapper.GetAPIGroupResources(k8s.Discovery())

	rmapper := restmapper.NewDiscoveryRESTMapper(apiRes)

	for _, mk := range matchKinds {
		groups := make([]string, 0, len(mk.APIGroups))
		resources := make([]string, 0, len(mk.Kinds))

		uniqGroups := make(map[string]struct{})
		uniqResources := make(map[string]struct{})

		for _, apiGroup := range mk.APIGroups {
			for _, kind := range mk.Kinds {
				rm, err := rmapper.RESTMapping(schema.GroupKind{
					Group: apiGroup,
					Kind:  kind,
				})
				if err != nil {
					input.LogEntry.Warnf("Resource mapping failed. Group: %s, Kind: %s. Error: %s", apiGroup, kind, err)
					continue
				}

				uniqGroups[rm.Resource.Group] = struct{}{}
				uniqResources[rm.Resource.Resource] = struct{}{}
			}
		}

		for k := range uniqGroups {
			groups = append(groups, k)
		}

		for k := range uniqResources {
			resources = append(resources, k)
		}

		res = append(res, matchResoyrce{
			APIGroups: groups,
			Resources: resources,
		})
	}

	input.LogEntry.Infof("Find matchKinds: %s", matchKinds)
	input.LogEntry.Infof("make resources: %s", res)

	input.Values.Set("admissionPolicyEngine.internal.trackedResources", res)

	return nil
}

func filterExporterCM(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var cm corev1.ConfigMap

	err := sdk.FromUnstructured(obj, &cm)
	if err != nil {
		return nil, err
	}

	return cm.Data["kinds.yaml"], nil
}

type matchKind struct {
	APIGroups []string `json:"apiGroups"`
	Kinds     []string `json:"kinds"`
}

type matchResoyrce struct {
	APIGroups []string `json:"apiGroups"`
	Resources []string `json:"resources"`
}
