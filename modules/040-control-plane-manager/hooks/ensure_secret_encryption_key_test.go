/*
Copyright 2021 Flant JSC

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

package hooks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = Describe("Modules :: controlPlaneManager :: hooks :: ensure_secret_encryption_key ::", func() {
	f := HookExecutionConfigInit(`{"controlPlaneManager":{"internal":{}}}`, ``)

	FContext("Empty cluster", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(``))
			f.RunHook()
		})

		It("Must be executed successfully", func() {
			Expect(f).To(ExecuteSuccessfully())
		})

		It("Must create a ConfigMap with hash and set a value", func() {
			secret := f.KubernetesResource("Secret", kubeSystemNS, secretEncryptionKeySecretName)
			Expect(secret.Exists()).To(BeTrue())

			Expect(f.ValuesGet(secretEncryptionKeyValuePath).Exists()).To(BeTrue())
			Expect(f.ValuesGet(secretEncryptionKeyValuePath).String()).To(HaveLen(44))
		})
	})
})
