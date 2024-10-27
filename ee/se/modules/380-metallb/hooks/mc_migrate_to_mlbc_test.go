/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	_ "github.com/flant/addon-operator/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

const (
	config = `
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleConfig
metadata:
  name: metallb
spec:
  enabled: true
  version: 1
  settings:
    speaker:
      nodeSelector:
        node-role.deckhouse.io/metallb: ""
      tolerations:
        - effect: NoExecute
          key: dedicated.deckhouse.io
          operator: Equal
          value: frontend
    addressPools:
      - name: nginx-loadbalancer-pool1
        protocol: layer2
        addresses:
          - 192.168.70.100-192.168.70.110
      - name: nginx-loadbalancer-pool2
        protocol: layer2
        addresses:
          - 192.168.71.100-192.168.72.110
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: zone-a
  namespace: d8-metallb
spec:
  ipAddressPools:
  - pool-1
  - pool-2
  nodeSelectors:
  - matchLabels:
      zone: a
`
	expectedMLBC = `
---
apiVersion: network.deckhouse.io/v1alpha1
kind: MetalLoadBalancerClass
metadata:
  name: l2-default
spec:
  isDefault: true
  type: L2
  addressPool:
  - 192.168.70.100-192.168.70.110
  - 192.168.71.100-192.168.72.110
  nodeSelector:
    node-role.deckhouse.io/metallb: ""
    zone: a
  tolerations:
  - effect: NoExecute
    key: dedicated.deckhouse.io
    operator: Equal
    value: frontend
`
)

var _ = Describe("Metallb hooks :: migrate MC to MetalLoadBalancerClass ::", func() {
	f := HookExecutionConfigInit(`{"metallb":{"internal":{}}}`, "")
	f.RegisterCRD("deckhouse.io", "v1alpha1", "ModuleConfig", false)
	f.RegisterCRD("network.deckhouse.io", "v1alpha1", "MetalLoadBalancerClass", false)
	f.RegisterCRD("metallb.io", "v1beta1", "L2Advertisement", true)
	f.RegisterCRD("metallb.io", "v1beta1", "IPAddressPool", true)

	Context("Empty cluster", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(config))
			f.RunHook()
		})
		It("Should run", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.BindingContexts.Array()).ShouldNot(BeEmpty())
		})
	})

	Context("Cluster with Metallb ModuleConfig", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(config))
			f.RunHook()
		})

		It("Created a new resource based on ModuleConfig", func() {
			Expect(f).To(ExecuteSuccessfully())

			MLBC := f.KubernetesResource("MetalLoadBalancerClass", "", "l2-default")
			Expect(MLBC.ToYaml()).To(MatchYAML(expectedMLBC))
		})
	})
})
