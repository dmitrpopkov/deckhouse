// Copyright 2024 Flant JSC
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

package v1

type DynamixCloudProviderDiscoveryData struct {
	APIVersion string       `json:"apiVersion,omitempty"`
	Kind       string       `json:"kind,omitempty"`
	SEPs       []DynamixSEP `json:"storageDomains,omitempty"`
}

type DynamixSEP struct {
	Name      string   `json:"name"`
	Pools     []string `json:"pools"`
	IsEnabled bool     `json:"isEnabled,omitempty"`
	IsDefault bool     `json:"IsDefault,omitempty"`
}