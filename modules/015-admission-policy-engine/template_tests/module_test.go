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

package template_tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/helm"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

const (
	globalValues = `
deckhouseVersion: test
enabledModules: ["vertical-pod-autoscaler-crd", "prometheus", "operator-prometheus-crd"]
clusterConfiguration:
  apiVersion: deckhouse.io/v1
  kind: ClusterConfiguration
  clusterDomain: cluster.local
  clusterType: Static
  kubernetesVersion: "1.21"
  podSubnetCIDR: 10.111.0.0/16
  podSubnetNodeCIDRPrefix: "24"
  serviceSubnetCIDR: 10.222.0.0/16
discovery:
  clusterMasterCount: 3
  prometheusScrapeInterval: 30
  d8SpecificNodeCountByRole:
    system: 1
modulesImages:
  registry: registry.deckhouse.io/deckhouse/fe
  registryDockercfg: Y2ZnCg==
  registryAddress: registry.deckhouse.io
  registryPath: /deckhouse/fe
  registryCA: CACACA
  registryScheme: https
  tags:
    common:
      kubeRbacProxy: hash
    admissionPolicyEngine:
      gatekeeper: hash
modules:
  placement: {}
`
)

var _ = Describe("Module :: admissionPolicyEngine :: helm template ::", func() {
	f := SetupHelmConfig(`{admissionPolicyEngine: {podSecurityStandards: {}, internal: {webhook: {ca: YjY0ZW5jX3N0cmluZwo=, crt: YjY0ZW5jX3N0cmluZwo=, key: YjY0ZW5jX3N0cmluZwo=}}}}`)

	Context("Cluster with deckhouse on master node", func() {
		BeforeEach(func() {
			f.ValuesSetFromYaml("global", globalValues)
			// f.ValuesSetFromYaml("deckhouse", moduleValuesForMasterNode)
			f.HelmRender()
		})

		nsName := "d8-admission-policy-engine"

		It("Everything must render properly", func() {
			Expect(f.RenderError).ShouldNot(HaveOccurred())
			sa := f.KubernetesResource("ServiceAccount", nsName, "admission-policy-engine")
			dp := f.KubernetesResource("Deployment", nsName, "gatekeeper-controller-manager")
			Expect(sa.Exists()).To(BeTrue())
			Expect(dp.Exists()).To(BeTrue())
		})
	})
})
