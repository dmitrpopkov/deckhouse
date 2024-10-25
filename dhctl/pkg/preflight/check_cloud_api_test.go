package preflight

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/deckhouse/deckhouse/dhctl/pkg/config"
)

func TestGetCloudApiURLFromMetaConfig(t *testing.T) {
	tests := []struct {
		name               string
		providerName       string
		providerConfigJSON string
		expectedURL        string
	}{
		{
			name:         "OpenStack provider",
			providerName: "OpenStack",
			providerConfigJSON: `{
				"authURL": "https://openstack.example.com/v3/auth",
				"domainName": "provider.local",
				"tenantID": "tenantID",
				"username": "username",
				"password": "password",
				"region": "eu-3"
			}`,
			expectedURL: "https://openstack.example.com/v3/auth",
		},
		{
			name:         "vSphere provider",
			providerName: "vSphere",
			providerConfigJSON: `{
				"server": "https://vsphere.example.com/sdk",
				"username": "vsphereUser",
				"password": "vspherePass"
			}`,
			expectedURL: "https://vsphere.example.com/sdk",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := require.New(t)
			clusterConfigYAML := `
---
apiVersion: deckhouse.io/v1
kind: InitConfiguration
deckhouse:
  imagesRepo: registry.deckhouse.io/deckhouse/ce
  releaseChannel: Alpha
`
			metaConfig, err := config.ParseConfigFromData(clusterConfigYAML)
			s.NoError(err)
			metaConfig.ProviderName = tt.providerName
			metaConfig.ProviderClusterConfig = map[string]json.RawMessage{
				"provider": json.RawMessage(tt.providerConfigJSON),
			}
			s.Equal(tt.providerName, metaConfig.ProviderName)
			cloudApiURL, err := getCloudApiURLFromMetaConfig(metaConfig)
			s.NoError(err)
			s.Equal(tt.expectedURL, cloudApiURL)
		})
	}
}
