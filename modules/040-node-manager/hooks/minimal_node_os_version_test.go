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

package hooks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/deckhouse/deckhouse/go_lib/dependency/requirements"
	. "github.com/deckhouse/deckhouse/testing/hooks"
)

const (
	node1 = `
---
apiVersion: v1
kind: Node
metadata:
  name: node1
  labels:
    node.deckhouse.io/group: group
status:
  nodeInfo:
    osImage: Ubuntu 20.04.3 LTS
`
	node2 = `
---
apiVersion: v1
kind: Node
metadata:
  name: node2
  labels:
    node.deckhouse.io/group: group2
status:
  nodeInfo:
    osImage: Ubuntu 18.04.5 LTS
`
	node3 = `
---
apiVersion: v1
kind: Node
metadata:
  name: node3
  labels:
    node.deckhouse.io/group: group3
status:
  nodeInfo:
    osImage: CentOS Linux 7 (Core)
`
)

var _ = Describe("node-manager :: minimal_node_os_version ", func() {
	f := HookExecutionConfigInit(`{}`, `{}`)

	Context("Nodes objects is not find", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(""))
			f.RunHook()
		})

		It("Should have no minimal version", func() {
			Expect(f).To(ExecuteSuccessfully())
			_, exists := requirements.GetValue(minVersionValuesKey)
			Expect(exists).To(BeFalse())
		})
	})

	Context("One node with Ubuntu OS", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(node1))
			f.RunHook()
		})

		It("Should have minimal version", func() {
			Expect(f).To(ExecuteSuccessfully())
			value, exists := requirements.GetValue(minVersionValuesKey)
			Expect(exists).To(BeTrue())
			Expect(value).To(BeEquivalentTo("20.4.3"))
		})
	})

	Context("One node with Centos OS", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(node3))
			f.RunHook()
		})

		It("Should have no minimal version", func() {
			Expect(f).To(ExecuteSuccessfully())
			_, exists := requirements.GetValue(minVersionValuesKey)
			Expect(exists).To(BeFalse())
		})
	})

	Context("Two nodes with Ubuntu OS", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(node1 + node2))
			f.RunHook()
		})

		It("Should have minimal version", func() {
			Expect(f).To(ExecuteSuccessfully())
			value, exists := requirements.GetValue(minVersionValuesKey)
			Expect(exists).To(BeTrue())
			Expect(value).To(BeEquivalentTo("18.4.5"))
		})
	})

	Context("Two nodes with Ubuntu OS and one node with CentOS", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(node1 + node2 + node3))
			f.RunHook()
		})

		It("Should have minimal version", func() {
			Expect(f).To(ExecuteSuccessfully())
			value, exists := requirements.GetValue(minVersionValuesKey)
			Expect(exists).To(BeTrue())
			Expect(value).To(BeEquivalentTo("18.4.5"))
		})
	})
})
