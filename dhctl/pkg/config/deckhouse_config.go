package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/deckhouse/deckhouse/dhctl/pkg/app"
	"github.com/deckhouse/deckhouse/dhctl/pkg/log"
)

const (
	initConfigurationError = `%s field in InitConfiguration is deprecated.
Please use ModuleConfig 'deckhouse' section in configuration. Example:
---
apiVersion: deckhouse.io/v1alpha1
kind: ClusterConfiguration
...
apiVersion: deckhouse.io/v1alpha1
kind: InitConfiguration
...
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleConfig
metadata:
  name: deckhouse
spec:
  settings:
    %s
`
)

type DeckhouseInstaller struct {
	Registry              RegistryData
	LogLevel              string
	Bundle                string
	DevBranch             string
	UUID                  string
	KubeDNSAddress        string
	ClusterConfig         []byte
	ProviderClusterConfig []byte
	StaticClusterConfig   []byte
	TerraformState        []byte
	NodesTerraformState   map[string][]byte
	CloudDiscovery        []byte
	ModuleConfigs         []*ModuleConfig

	KubeadmBootstrap   bool
	MasterNodeSelector bool
}

func (c *DeckhouseInstaller) GetImage(forceVersionTag bool) string {
	registryNameTemplate := "%s%s:%s"
	tag := c.DevBranch
	if forceVersionTag {
		versionTag, foundValidTag := readVersionTagFromInstallerContainer()
		if foundValidTag {
			tag = versionTag
		}
	}

	if tag == "" {
		panic("Probably you use development image. please use devBranch")
	}

	return fmt.Sprintf(registryNameTemplate, c.Registry.Address, c.Registry.Path, tag)
}

func (c *DeckhouseInstaller) IsRegistryAccessRequired() bool {
	return c.Registry.DockerCfg != ""
}

func readVersionTagFromInstallerContainer() (string, bool) {
	rawFile, err := os.ReadFile(app.VersionFile)
	if err != nil {
		log.WarnF(
			"Could not read %s: %v\nWill fall back to installation from release channel or dev branch.",
			app.VersionFile, err,
		)
		return "", false
	}

	tag := string(rawFile)
	if _, err = semver.NewVersion(tag); err != nil {
		return "", false
	}

	return tag, true
}

func PrepareDeckhouseInstallConfig(metaConfig *MetaConfig) (*DeckhouseInstaller, error) {
	clusterConfig, err := metaConfig.ClusterConfigYAML()
	if err != nil {
		return nil, fmt.Errorf("marshal cluster config: %v", err)
	}

	providerClusterConfig, err := metaConfig.ProviderClusterConfigYAML()
	if err != nil {
		return nil, fmt.Errorf("marshal provider config: %v", err)
	}

	staticClusterConfig, err := metaConfig.StaticClusterConfigYAML()
	if err != nil {
		return nil, fmt.Errorf("marshal static config: %v", err)
	}

	bundle := "Default"
	logLevel := "Info"

	// todo after release 1.55 remove it and from openapi schema
	deprecatedFields := make([]string, 0, 3)
	deprecatedFieldsExamples := make([]string, 0, 3)
	if metaConfig.DeckhouseConfig.ReleaseChannel != "" {
		deprecatedFields = append(deprecatedFields, "releaseChannel")
		deprecatedFieldsExamples = append(deprecatedFieldsExamples, "releaseChannel: Stable")
	}

	if metaConfig.DeckhouseConfig.Bundle != "" {
		bundle = metaConfig.DeckhouseConfig.Bundle
		deprecatedFields = append(deprecatedFields, "bundle")
		deprecatedFieldsExamples = append(deprecatedFieldsExamples, "bundle: Default")
	}

	if metaConfig.DeckhouseConfig.LogLevel != "" {
		logLevel = metaConfig.DeckhouseConfig.LogLevel
		deprecatedFields = append(deprecatedFields, "logLevel")
		deprecatedFieldsExamples = append(deprecatedFieldsExamples, "logLevel: Info")
	}

	if len(deprecatedFields) > 0 {
		log.WarnF(initConfigurationError, strings.Join(deprecatedFields, ","), strings.Join(deprecatedFieldsExamples, "\n    "))
	}

	if len(metaConfig.DeckhouseConfig.ConfigOverrides) > 0 {
		log.WarnLn(`
Config overrides is deprecated. Please use module config:
---
apiVersion: deckhouse.io/v1alpha1
kind: ClusterConfiguration
...
apiVersion: deckhouse.io/v1alpha1
kind: InitConfiguration
...
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleConfig
metadata:
  name: global
spec:
  settings:
    highAvailability: false
    modules:
      publicDomainTemplate: '%s.example.com'
  version: 1
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleConfig
metadata:
  name: cni-flannel
spec:
  enabled: true
---
...
`)
		if len(metaConfig.ModuleConfigs) > 0 {
			return nil, fmt.Errorf("Cannot use ModuleConfig's and configOverrides. Please use ModuleConfig's")
		}
		mcs, err := ConvertInitConfigurationToModuleConfigs(metaConfig)
		if err != nil {
			return nil, err
		}

		metaConfig.ModuleConfigs = mcs
	}

	// find deckhouse module config for extract release
	for _, mc := range metaConfig.ModuleConfigs {
		if mc.GetName() != "deckhouse" {
			continue
		}
		logLevelRaw, ok := mc.Spec.Settings["logLevel"]
		if ok {
			logLevel = logLevelRaw.(string)
		}
		bundleRaw, ok := mc.Spec.Settings["bundle"]
		if ok {
			bundle = bundleRaw.(string)
		}
	}

	installConfig := DeckhouseInstaller{
		UUID:                  metaConfig.UUID,
		Registry:              metaConfig.Registry,
		DevBranch:             metaConfig.DeckhouseConfig.DevBranch,
		Bundle:                bundle,
		LogLevel:              logLevel,
		KubeDNSAddress:        metaConfig.ClusterDNSAddress,
		ProviderClusterConfig: providerClusterConfig,
		StaticClusterConfig:   staticClusterConfig,
		ClusterConfig:         clusterConfig,
		ModuleConfigs:         metaConfig.ModuleConfigs,
	}

	return &installConfig, nil
}
