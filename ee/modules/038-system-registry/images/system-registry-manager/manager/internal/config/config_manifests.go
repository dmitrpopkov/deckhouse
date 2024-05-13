/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultDirMode  = os.FileMode(0755)
	DefaultFileMode = os.FileMode(0755)
	DefaultCertMode = os.FileMode(644)
)

var (
	TmpWorkspaceDir                       = filepath.Join(os.TempDir(), "system-registry-manager/workspace")
	TmpWorkspaceCertsDir                  = filepath.Join(TmpWorkspaceDir, "pki")
	TmpWorkspaceManifestsDir              = filepath.Join(TmpWorkspaceDir, "manifests")
	TmpWorkspaceStaticPodsDir             = filepath.Join(TmpWorkspaceManifestsDir, "static_pods")
	TmpWorkspaceSeaweedManifestsDir       = filepath.Join(TmpWorkspaceManifestsDir, "seaweedfs_config")
	TmpWorkspaceDockerAuthManifestsDir    = filepath.Join(TmpWorkspaceManifestsDir, "auth_config")
	TmpWorkspaceDockerDistribManifestsDir = filepath.Join(TmpWorkspaceManifestsDir, "distribution_config")

	InputCertsDir                  = "/pki"
	InputManifestsDir              = "/manifests"
	InputStaticPodsDir             = filepath.Join(InputManifestsDir, "static_pods")
	InputSeaweedManifestsDir       = filepath.Join(InputManifestsDir, "seaweedfs_config")
	InputDockerAuthManifestsDir    = filepath.Join(InputManifestsDir, "auth_config")
	InputDockerDistribManifestsDir = filepath.Join(InputManifestsDir, "distribution_config")

	DestionationDir                       = "/etc/kubernetes"
	DestinationSystemRegistryDir          = filepath.Join(DestionationDir, "system-registry")
	DestinationCertsDir                   = filepath.Join(DestinationSystemRegistryDir, "pki")
	DestionationDirStaticPodsDir          = filepath.Join(DestionationDir, "manifests")
	DestionationSeaweedManifestsDir       = filepath.Join(DestinationSystemRegistryDir, "seaweedfs_config")
	DestionationDockerAuthManifestsDir    = filepath.Join(DestinationSystemRegistryDir, "auth_config")
	DestionationDockerDistribManifestsDir = filepath.Join(DestinationSystemRegistryDir, "distribution_config")
)

type ManifestSpec struct {
	InputPath string
	TmpPath   string
	DestPath  string
}

type CaCertificateSpec struct {
	InputPath string
	TmpPath   string
}

type CertificateSpec struct {
	TmpGeneratePath string
	DestPath        string
}

type BaseCertificatesSpec struct {
	CACrt     CaCertificateSpec
	CAKey     CaCertificateSpec
	EtcdCACrt CaCertificateSpec
	EtcdCAKey CaCertificateSpec
}

type GeneratedCertificateSpec struct {
	CAKey  CaCertificateSpec
	CACert CaCertificateSpec
	Key    CertificateSpec
	Cert   CertificateSpec
}

type GeneratedCertificatesSpec struct {
	SeaweedEtcdClientCert GeneratedCertificateSpec
	DockerAuthTokenCert   GeneratedCertificateSpec
}

type ManifestsSpec struct {
	GeneratedCertificates GeneratedCertificatesSpec
	Manifests             []ManifestSpec
}
type DataForManifestRendering struct {
	DiscoveredNodeIP string
}

func NewManifestsSpec() *ManifestsSpec {
	baseCertificates := BaseCertificatesSpec{
		CACrt: CaCertificateSpec{
			InputPath: filepath.Join(InputCertsDir, "ca.crt"),
			TmpPath:   filepath.Join(TmpWorkspaceCertsDir, "ca.crt"),
		},
		CAKey: CaCertificateSpec{
			InputPath: filepath.Join(InputCertsDir, "ca.key"),
			TmpPath:   filepath.Join(TmpWorkspaceCertsDir, "ca.key"),
		},
		EtcdCACrt: CaCertificateSpec{
			InputPath: filepath.Join(InputCertsDir, "etcd-ca.crt"),
			TmpPath:   filepath.Join(TmpWorkspaceCertsDir, "etcd-ca.crt"),
		},
		EtcdCAKey: CaCertificateSpec{
			InputPath: filepath.Join(InputCertsDir, "etcd-ca.key"),
			TmpPath:   filepath.Join(TmpWorkspaceCertsDir, "etcd-ca.key"),
		},
	}

	generatedCertificates := GeneratedCertificatesSpec{
		SeaweedEtcdClientCert: GeneratedCertificateSpec{
			CAKey:  baseCertificates.EtcdCAKey,
			CACert: baseCertificates.EtcdCACrt,
			Key: CertificateSpec{
				TmpGeneratePath: filepath.Join(TmpWorkspaceCertsDir, "seaweedfs-etcd-client.key"),
				DestPath:        filepath.Join(DestinationCertsDir, "seaweedfs-etcd-client.key"),
			},
			Cert: CertificateSpec{
				TmpGeneratePath: filepath.Join(TmpWorkspaceCertsDir, "seaweedfs-etcd-client.crt"),
				DestPath:        filepath.Join(DestinationCertsDir, "seaweedfs-etcd-client.crt"),
			},
		},
		DockerAuthTokenCert: GeneratedCertificateSpec{
			CAKey:  baseCertificates.EtcdCAKey,
			CACert: baseCertificates.EtcdCACrt,
			Key: CertificateSpec{
				TmpGeneratePath: filepath.Join(TmpWorkspaceCertsDir, "token.key"),
				DestPath:        filepath.Join(DestinationCertsDir, "token.key"),
			},
			Cert: CertificateSpec{
				TmpGeneratePath: filepath.Join(TmpWorkspaceCertsDir, "token.crt"),
				DestPath:        filepath.Join(DestinationCertsDir, "token.crt"),
			},
		},
	}

	manifestsSpec := ManifestsSpec{
		// BaseCertificates:      baseCertificates,
		GeneratedCertificates: generatedCertificates,
		Manifests: []ManifestSpec{
			{
				InputPath: filepath.Join(InputDockerDistribManifestsDir, "config.yaml"),
				TmpPath:   filepath.Join(TmpWorkspaceDockerDistribManifestsDir, "config.yaml"),
				DestPath:  filepath.Join(DestionationDockerDistribManifestsDir, "config.yaml"),
			},
			{
				InputPath: filepath.Join(InputDockerAuthManifestsDir, "config.yaml"),
				TmpPath:   filepath.Join(TmpWorkspaceDockerAuthManifestsDir, "config.yaml"),
				DestPath:  filepath.Join(DestionationDockerAuthManifestsDir, "config.yaml"),
			},
			{
				InputPath: filepath.Join(InputSeaweedManifestsDir, "filer.toml"),
				TmpPath:   filepath.Join(TmpWorkspaceSeaweedManifestsDir, "filer.toml"),
				DestPath:  filepath.Join(DestionationSeaweedManifestsDir, "filer.toml"),
			},
			{
				InputPath: filepath.Join(InputSeaweedManifestsDir, "master.toml"),
				TmpPath:   filepath.Join(TmpWorkspaceSeaweedManifestsDir, "master.toml"),
				DestPath:  filepath.Join(DestionationSeaweedManifestsDir, "master.toml"),
			},
			{
				InputPath: filepath.Join(InputStaticPodsDir, "system-registry.yaml"),
				TmpPath:   filepath.Join(TmpWorkspaceStaticPodsDir, "system-registry.yaml"),
				DestPath:  filepath.Join(DestionationDirStaticPodsDir, "system-registry.yaml"),
			},
		},
	}
	return &manifestsSpec
}

func GetDataForManifestRendering() DataForManifestRendering {
	cfg := GetConfig()
	return DataForManifestRendering{
		DiscoveredNodeIP: cfg.MyIP,
	}
}
