/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/deckhouse/deckhouse/go_lib/dependency"
	. "github.com/deckhouse/deckhouse/testing/hooks"
)

const (
	values = `
internal: {}
`
	service1 = `
---
apiVersion: v1
kind: Service
metadata:
  name: nginx_serv1
  namespace: nginx1
spec:
  clusterIP: 1.2.3.4
  ports:
  - port: 7473
    protocol: TCP
    targetPort: 7473
  externalTrafficPolicy: Local
  internalTrafficPolicy: Cluster
  selector:
    app: nginx1
  type: LoadBalancer
status:
  conditions:
  - message: 2 of 3 public IPs were assigned
    reason: NotAllIPsAssigned
    status: "False"
    type: AllPublicIPsAssigned
  - message: status_mlbc
    reason: LoadBalancerClassBound
    status: "True"
    type: network.deckhouse.io/load-balancer-class
`
	service2 = `
---
apiVersion: v1
kind: Service
metadata:
  name: nginx_serv2
  namespace: nginx2
  annotations:
    network.deckhouse.io/l2-load-balancer-ips: "7.7.7.7"
spec:
  clusterIP: 2.3.4.5
  ports:
  - port: 7474
    protocol: TCP
    targetPort: 7474
  externalTrafficPolicy: Local
  internalTrafficPolicy: Cluster
  selector:
    app: nginx2
  type: LoadBalancer
`
	service3 = `
---
apiVersion: v1
kind: Service
metadata:
  name: nginx_serv3
  namespace: nginx3
  annotations:
    metallb.universe.tf/address-pool: "test"
spec:
  clusterIP: 4.5.6.7
  ports:
  - port: 7475
    protocol: TCP
    targetPort: 7475
  externalTrafficPolicy: Local
  internalTrafficPolicy: Cluster
  selector:
    app: nginx3
  type: LoadBalancer
status:
  loadBalancer:
    ingress:
    - ip: 12.12.12.12
    - ip: 13.13.13.13`
)

var _ = Describe("Metallb hooks :: prepare services to migrate ::", func() {
	f := HookExecutionConfigInit(`{"metallb":{"internal":{}}}`, "")

	Context("There are 3 services, but only one  must be patched", func() {
		BeforeEach(func() {
			f.KubeStateSet("")
			f.ValuesSetFromYaml("metallb", []byte(values))

			for _, clusterService := range []string{service1, service2, service3} {
				var service v1.Service
				err := yaml.Unmarshal([]byte(clusterService), &service)
				Expect(err).To(BeNil())
				_, err = dependency.TestDC.MustGetK8sClient().
					CoreV1().
					Services(service.GetNamespace()).
					Create(context.TODO(), &service, metav1.CreateOptions{})
				Expect(err).To(BeNil())
			}

			f.RunHook()
		})

		It("Should service is patched", func() {
			Expect(f).To(ExecuteSuccessfully())
			k8sClient := f.BindingContextController.FakeCluster().Client
			patchedService, err := k8sClient.CoreV1().Services("nginx3").Get(
				context.TODO(),
				"nginx_serv3",
				metav1.GetOptions{},
			)

			Expect(err).To(BeNil())
			Expect(patchedService.Annotations["network.deckhouse.io/l2-load-balancer-ips"]).To(
				Equal("12.12.12.12,13.13.13.13"),
			)
		})
	})
})
