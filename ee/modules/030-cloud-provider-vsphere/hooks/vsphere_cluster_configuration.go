/*
Copyright 2021 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"encoding/json"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/deckhouse/deckhouse/dhctl/pkg/config"
	v1 "github.com/deckhouse/deckhouse/ee/modules/030-cloud-provider-vsphere/hooks/internal/v1"
	"github.com/deckhouse/deckhouse/go_lib/hooks/cluster_configuration"
)

var _ = cluster_configuration.RegisterHook(func(input *go_hook.HookInput, metaCfg *config.MetaConfig, providerDiscoveryData *unstructured.Unstructured, secretFound bool) error {

	p := map[string]json.RawMessage{}
	if metaCfg != nil {
		p = metaCfg.ProviderClusterConfig
	}

	var providerClusterConfiguration v1.VsphereProviderClusterConfiguration
	err := convertJSONRawMessageToStruct(p, &providerClusterConfiguration)
	if err != nil {
		return err
	}

	var moduleConfiguration v1.VsphereModuleConfiguration
	err = json.Unmarshal([]byte(input.Values.Get("cloudProviderVsphere").String()), &moduleConfiguration)
	if err != nil {
		return err
	}

	overrideValues(providerClusterConfiguration, moduleConfiguration)
	input.Values.Set("cloudProviderVsphere.internal.providerClusterConfiguration", providerClusterConfiguration)

	var discoveryData v1.VsphereProviderClusterConfiguration
	err = sdk.FromUnstructured(providerDiscoveryData, &discoveryData)
	input.Values.Set("cloudProviderVsphere.internal.providerDiscoveryData", discoveryData)

	return nil
})

func convertJSONRawMessageToStruct(in map[string]json.RawMessage, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, out)
	if err != nil {
		return err
	}
	return nil
}

func overrideValues(p v1.VsphereProviderClusterConfiguration, m v1.VsphereModuleConfiguration) {
	if m.Host != nil {
		p.Provider.Server = m.Host
	}

	if m.Username != nil {
		p.Provider.Username = m.Username
	}

	if m.Password != nil {
		p.Provider.Password = m.Password
	}

	if m.Insecure != nil {
		p.Provider.Insecure = m.Insecure
	}

	if m.RegionTagCategory != nil {
		p.RegionTagCategory = m.RegionTagCategory
	}

	if m.ZoneTagCategory != nil {
		p.ZoneTagCategory = m.ZoneTagCategory
	}

	if m.DisableTimesync != nil {
		p.DisableTimesync = m.DisableTimesync
	}

	if m.ExternalNetworkNames != nil {
		p.ExternalNetworkNames = m.ExternalNetworkNames
	}

	if m.InternalNetworkNames != nil {
		p.ExternalNetworkNames = m.ExternalNetworkNames
	}

	if m.Region != nil {
		p.Region = m.Region
	}

	if m.Zones != nil {
		p.Zones = m.Zones
	}
	/*
			{
				ConfigKey: "vmFolderPath",
				ValueKey:  "vmFolderPath",
			},
			{
				ConfigKey: "sshKeys.0",
				ValueKey:  "sshPublicKey",
			},
		}
	*/
}
