/*
Copyright 2022 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package istio_versions

import "encoding/json"

type IstioVersionsMap map[string]IstioVersionInfo

type IstioVersionInfo struct {
	FullVersion string `json:"fullVersion"`
	Revision    string `json:"revision"`
	ImageSuffix string `json:"imageSuffix"`
}

func (vm IstioVersionsMap) GetVersionByRevision(rev string) string {
	for ver, istioVerInfo := range vm {
		if istioVerInfo.Revision == rev {
			return ver
		}
	}
	return ""
}

func (vm IstioVersionsMap) IsRevisionSupported(rev string) bool {
	for _, istioVerInfo := range vm {
		if istioVerInfo.Revision == rev {
			return true
		}
	}
	return false
}

func (vm IstioVersionsMap) GetFullVersionByRevision(rev string) string {
	for _, istioVerInfo := range vm {
		if istioVerInfo.Revision == rev {
			return istioVerInfo.FullVersion
		}
	}
	return ""
}

func (vm IstioVersionsMap) GetAllVersions() []string {
	versions := make([]string, len(vm))
	for ver := range vm {
		versions = append(versions, ver)
	}
	return versions
}

func VersionMapStrToVersionMap(versionMapRaw string) IstioVersionsMap {
	versionMap := make(IstioVersionsMap)
	_ = json.Unmarshal([]byte(versionMapRaw), &versionMap)
	return versionMap
}
