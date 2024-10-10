// Copyright 2022 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"

	"github.com/Masterminds/semver/v3"
	"github.com/deckhouse/deckhouse/go_lib/dependency"
	"github.com/deckhouse/deckhouse/go_lib/dependency/cr"
	"github.com/ettle/strcase"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"gopkg.in/yaml.v2"
)

type ModuleService struct {
	dc dependency.Container

	registry        string
	registryOptions []cr.Option
}

func NewModuleService(registryAddress string, registryConfig *RegistryConfig) *ModuleService {
	return &ModuleService{
		dc:              dependency.NewDependencyContainer(),
		registry:        registryAddress,
		registryOptions: GenerateRegistryOptions(registryConfig),
	}
}

func (svc *ModuleService) ListModules() ([]string, error) {
	regCli, err := svc.dc.GetRegistryClient(svc.registry, svc.registryOptions...)
	if err != nil {
		return nil, fmt.Errorf("get registry client: %v", err)
	}

	ls, err := regCli.ListTags()
	if err != nil {
		return nil, fmt.Errorf("list tags: %v", err)
	}

	return ls, err
}

func (svc *ModuleService) ListModuleTags(moduleName string, fullList bool) ([]string, error) {
	regCli, err := svc.dc.GetRegistryClient(path.Join(svc.registry, moduleName), svc.registryOptions...)
	if err != nil {
		return nil, fmt.Errorf("get registry client: %v", err)
	}

	ls, err := regCli.ListTags()
	if err != nil {
		return nil, fmt.Errorf("list tags: %v", err)
	}

	return ls, err
}

type moduleReleaseMetadata struct {
	Version *semver.Version `json:"version"`

	Changelog map[string]any
}

func (svc *ModuleService) GetModuleRelease(moduleName, releaseChannel string) (*moduleReleaseMetadata, error) {
	regCli, err := svc.dc.GetRegistryClient(path.Join(svc.registry, moduleName, "release"), svc.registryOptions...)
	if err != nil {
		return nil, fmt.Errorf("get registry client: %v", err)
	}

	img, err := regCli.Image(strcase.ToKebab(releaseChannel))
	if err != nil {
		return nil, fmt.Errorf("fetch image error: %v", err)
	}

	moduleMetadata, err := svc.fetchModuleReleaseMetadata(img)
	if err != nil {
		return nil, fmt.Errorf("fetch module release metadata error: %v", err)
	}

	if moduleMetadata.Version == nil {
		return nil, fmt.Errorf("module release %q metadata malformed: no version found", moduleName)
	}

	return moduleMetadata, nil
}

func (svc *ModuleService) fetchModuleReleaseMetadata(img v1.Image) (*moduleReleaseMetadata, error) {
	var meta = new(moduleReleaseMetadata)

	rc := mutate.Extract(img)
	defer rc.Close()

	rr := &releaseReader{
		versionReader:   bytes.NewBuffer(nil),
		changelogReader: bytes.NewBuffer(nil),
	}

	err := rr.untarModuleMetadata(rc)
	if err != nil {
		return nil, err
	}

	if rr.versionReader.Len() > 0 {
		err = json.NewDecoder(rr.versionReader).Decode(&meta)
		if err != nil {
			return nil, err
		}
	}

	if rr.changelogReader.Len() > 0 {
		var changelog map[string]any

		err = yaml.NewDecoder(rr.changelogReader).Decode(&changelog)

		if err != nil {
			meta.Changelog = make(map[string]any)
			return nil, nil
		}

		meta.Changelog = changelog
	}

	return meta, nil
}
