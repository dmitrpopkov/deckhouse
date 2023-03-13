package cloudprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	"deckhouse.io/uibackend/util/retry"
)

type ProviderConfig struct {
	Config        map[string]interface{}
	DiscoveryData map[string]interface{}
}

func (c *ProviderConfig) Zones() []string {
	zones, ok, _ := unstructured.NestedStringSlice(c.DiscoveryData, "zones")
	if !ok {
		return make([]string, 0)
	}

	return zones
}

func GetClusterConfig(ctx context.Context, client *kubernetes.Clientset) (*ProviderConfig, error) {
	var secret *v1.Secret
	err := retry.NewLoop("Get provider cluster config", 3, 3*time.Second).Run(func() error {
		var err error
		secret, err = client.CoreV1().Secrets("kube-system").Get(ctx, "d8-provider-cluster-configuration", metav1.GetOptions{})
		return err
	})

	if err != nil {
		return nil, err
	}

	var discoveryData map[string]interface{}
	var config map[string]interface{}

	if clusterConfigurationYAML, ok := secret.Data["cloud-provider-cluster-configuration.yaml"]; ok && len(clusterConfigurationYAML) > 0 {
		if err = yaml.Unmarshal(clusterConfigurationYAML, &config); err != nil {
			return nil, fmt.Errorf("config unmarshal: %v", err)
		}
	}
	if discoveryDataJSON, ok := secret.Data["cloud-provider-discovery-data.json"]; ok && len(discoveryDataJSON) > 0 {
		if err = json.Unmarshal(discoveryDataJSON, &discoveryData); err != nil {
			return nil, fmt.Errorf("cannot unmarshal cloud-provider-discovery-data.json key: %v", err)
		}
	}

	return &ProviderConfig{
		Config:        config,
		DiscoveryData: discoveryData,
	}, nil
}
