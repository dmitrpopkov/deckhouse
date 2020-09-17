/*

User-stories:
1. There are special nodes for kube-dns in cluster — hook must fit kube-dns deployment to this nodes and masters. If there is kube-dns-autoscaler in cluster then hook must keep replicas.
2. There aren't dedicated dns-nodes, but there are special system-nodes in cluster — hook must fit kube-dns deployment to this nodes and masters. If there is kube-dns-autoscaler in cluster then hook must keep replicas.
3. There aren't special nodes — hook must fit kube-dns deployment to this nodes. Replicas must be counted by formula: ([([2,<count_master_nodes>,<original_replicas>] | max), ([2, '<count_master_nodes + count_nonspecific_nodes>'] | max)] | min).
4. kube-dns deployment should aim to fit pods to different nodes.
5. If there are empty fields in affinity then hook must delete them.

*/

package hooks

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Kube DNS :: policy ::", func() {
	const (
		initValuesString = `
kubeDns:
  enableLogs: false
  internal:
    replicas: 2
    enablePodAntiAffinity: false
`
		initConfigValuesString = `{}`
		stateMaster            = `

`
	)
	kubeMasterNodes := func(quantity int) string {
		result := ``
		for i := 1; i <= quantity; i++ {
			result += fmt.Sprintf(`---
apiVersion: v1
kind: Node
metadata:
  name: master-%d
  labels:
    node-role.kubernetes.io/master: ""
`, i)
		}
		return result
	}
	kubeDnsNodes := func(quantity int) string {
		result := ``
		for i := 1; i <= quantity; i++ {
			result += fmt.Sprintf(`---
apiVersion: v1
kind: Node
metadata:
  name: dns-%d
  labels:
    node-role.deckhouse.io/kube-dns: ""
`, i)
		}
		return result
	}

	f := HookExecutionConfigInit(initValuesString, initConfigValuesString)

	Context("With only Master node in a cluster", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateMaster))
			f.RunHook()
		})

		It("", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.BindingContexts.Array()).ShouldNot(BeEmpty())

			Expect(f.ValuesGet("kubeDns.internal.replicas").Int()).To(Equal(int64(2)))
			Expect(f.ValuesGet("kubeDns.internal.specificNodeType").Exists()).To(BeFalse())
			Expect(f.ValuesGet("kubeDns.internal.enablePodAntiAffinity").Bool()).To(BeFalse())
		})
	})

	Context("With 2 Master nodes and 2 specific nodes in a cluster", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(kubeMasterNodes(3) + kubeDnsNodes(1)))
			f.RunHook()
		})

		It("", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.BindingContexts.Array()).ShouldNot(BeEmpty())

			Expect(f.ValuesGet("kubeDns.internal.replicas").Int()).To(Equal(int64(4)))
			Expect(f.ValuesGet("kubeDns.internal.specificNodeType").String()).To(Equal("kube-dns"))
			Expect(f.ValuesGet("kubeDns.internal.enablePodAntiAffinity").Bool()).To(BeTrue())
		})
	})

	Context("With Master node and 4 specific nodes in a cluster", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(kubeMasterNodes(1) + kubeDnsNodes(4)))
			f.RunHook()
		})

		It("", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.BindingContexts.Array()).ShouldNot(BeEmpty())

			Expect(f.ValuesGet("kubeDns.internal.replicas").Int()).To(Equal(int64(3))) // because of the master quantity + 2 limit
			Expect(f.ValuesGet("kubeDns.internal.specificNodeType").String()).To(Equal("kube-dns"))
			Expect(f.ValuesGet("kubeDns.internal.enablePodAntiAffinity").Bool()).To(BeTrue())
		})
	})
})
