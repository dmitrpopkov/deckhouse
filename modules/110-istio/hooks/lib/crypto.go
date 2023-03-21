/*
Copyright 2023 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

type IstioCA struct {
	Key   string `json:"key"`
	Cert  string `json:"cert"`
	Chain string `json:"chain"`
	Root  string `json:"root"`
}

type Keypair struct {
	Pub  string `json:"pub"`
	Priv string `json:"priv"`
}
