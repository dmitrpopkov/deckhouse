package hooks

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Global hooks :: migrate/panic_on_flant_com_taint ::", func() {

	const (
		initValuesString       = `{}`
		initConfigValuesString = `{}`
	)

	const (
		stateNodeWithDeprecatedTaints = `
---
apiVersion: v1
kind: Node
metadata:
  name: node-0
spec:
  taints:
  - effect: NoExecute
    key: dedicated.flant.com
    value: system
  - effect: NoExecute
    key: dedicated.flant.com
    value: other
---
apiVersion: v1
kind: Node
metadata:
  name: node-1
spec:
  taints:
  - effect: NoExecute
    key: dedicated.flant.com
    value: system
  - effect: NoExecute
    key: dedicated.flant.com
    value: other
`
		stateNodeWithGoodTaints = `
---
apiVersion: v1
kind: Node
metadata:
  name: node-0
spec:
  taints:
  - effect: NoExecute
    key: dedicated.flant.com
    value: production
---
apiVersion: v1
kind: Node
metadata:
  name: node-1
spec:
  taints:
  - effect: NoExecute
    key: dedicated.flant.com
    value: production
`
		stateNodeWithDeprecatedLabels = `
---
apiVersion: v1
kind: Node
metadata:
  name: node-0
  labels:
    node-role.flant.com/system: ""
    node-role.flant.com/frontend: ""
    node-role.flant.com/whatever: ""
---
apiVersion: v1
kind: Node
metadata:
  name: node-1
  labels:
    node-role.flant.com/whatever: ""
`
	)

	f := HookExecutionConfigInit(initValuesString, initConfigValuesString)

	Context("Empty cluster", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(``))
			f.RunHook()
		})

		It("Hook must not fail", func() {
			Expect(f).To(ExecuteSuccessfully())
		})
	})

	Context("Cluster with Node having deprecated flant.com taints", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateNodeWithDeprecatedTaints))
			f.RunHook()
		})

		It("Hook must fail", func() {
			Expect(f).To(Not(ExecuteSuccessfully()))
		})
	})

	Context("Cluster with Node having good flant.com taints", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateNodeWithGoodTaints))
			f.RunHook()
		})

		It("Hook must not fail", func() {
			Expect(f).To(ExecuteSuccessfully())
		})
	})

	Context("Cluster with Node having deprecated flant.com labels", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateNodeWithDeprecatedLabels))
			f.RunHook()
		})

		It("Hook must fail", func() {
			Expect(f).To(Not(ExecuteSuccessfully()))
		})
	})

})
