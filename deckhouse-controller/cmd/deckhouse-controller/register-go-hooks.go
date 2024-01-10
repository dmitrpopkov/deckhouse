// Code generated by "register.go" DO NOT EDIT.
package main

import (
	_ "github.com/flant/addon-operator/sdk"

	_ "github.com/deckhouse/deckhouse/ee/be/modules/350-node-local-dns/hooks"
	_ "github.com/deckhouse/deckhouse/ee/be/modules/600-flant-integration/hooks"
	_ "github.com/deckhouse/deckhouse/ee/be/modules/600-flant-integration/hooks/madison"
	_ "github.com/deckhouse/deckhouse/ee/be/modules/600-flant-integration/hooks/migrate"
	_ "github.com/deckhouse/deckhouse/ee/be/modules/600-flant-integration/hooks/pricing"
	_ "github.com/deckhouse/deckhouse/ee/fe/modules/340-monitoring-applications/hooks"
	_ "github.com/deckhouse/deckhouse/ee/fe/modules/500-basic-auth/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/030-cloud-provider-openstack/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/030-cloud-provider-vcd/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/030-cloud-provider-vsphere/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/110-istio/hooks/ee"
	_ "github.com/deckhouse/deckhouse/ee/modules/110-istio/hooks/ee/lib/crd"
	_ "github.com/deckhouse/deckhouse/ee/modules/150-user-authn/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/160-multitenancy-manager/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/160-multitenancy-manager/hooks/apis/deckhouse.io/v1alpha1"
	_ "github.com/deckhouse/deckhouse/ee/modules/160-multitenancy-manager/hooks/apis/deckhouse.io/v1alpha2"
	_ "github.com/deckhouse/deckhouse/ee/modules/300-prometheus/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/380-metallb/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/450-keepalived/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/450-network-gateway/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/500-operator-trivy/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/502-delivery/hooks"
	_ "github.com/deckhouse/deckhouse/ee/modules/502-delivery/hooks/https"
	_ "github.com/deckhouse/deckhouse/ee/modules/502-delivery/hooks/werf_sources"
	_ "github.com/deckhouse/deckhouse/ee/modules/650-runtime-audit-engine/hooks"
	_ "github.com/deckhouse/deckhouse/global-hooks"
	_ "github.com/deckhouse/deckhouse/global-hooks/deckhouse-config"
	_ "github.com/deckhouse/deckhouse/global-hooks/discovery"
	_ "github.com/deckhouse/deckhouse/global-hooks/migrate"
	_ "github.com/deckhouse/deckhouse/global-hooks/resources"
	_ "github.com/deckhouse/deckhouse/modules/000-common/hooks"
	_ "github.com/deckhouse/deckhouse/modules/002-deckhouse/hooks"
	_ "github.com/deckhouse/deckhouse/modules/005-external-module-manager/hooks"
	_ "github.com/deckhouse/deckhouse/modules/010-cloud-data-crd/hooks"
	_ "github.com/deckhouse/deckhouse/modules/010-operator-prometheus-crd/hooks"
	_ "github.com/deckhouse/deckhouse/modules/010-prometheus-crd/hooks"
	_ "github.com/deckhouse/deckhouse/modules/010-snapshot-controller-crd/hooks"
	_ "github.com/deckhouse/deckhouse/modules/010-user-authn-crd/hooks"
	_ "github.com/deckhouse/deckhouse/modules/010-vertical-pod-autoscaler-crd/hooks"
	_ "github.com/deckhouse/deckhouse/modules/011-flow-schema/hooks"
	_ "github.com/deckhouse/deckhouse/modules/015-admission-policy-engine/hooks"
	_ "github.com/deckhouse/deckhouse/modules/021-cni-cilium/hooks"
	_ "github.com/deckhouse/deckhouse/modules/021-kube-proxy/hooks"
	_ "github.com/deckhouse/deckhouse/modules/030-cloud-provider-aws/hooks"
	_ "github.com/deckhouse/deckhouse/modules/030-cloud-provider-azure/hooks"
	_ "github.com/deckhouse/deckhouse/modules/030-cloud-provider-gcp/hooks"
	_ "github.com/deckhouse/deckhouse/modules/030-cloud-provider-yandex/hooks"
	_ "github.com/deckhouse/deckhouse/modules/031-ceph-csi/hooks"
	_ "github.com/deckhouse/deckhouse/modules/031-local-path-provisioner/hooks"
	_ "github.com/deckhouse/deckhouse/modules/035-cni-flannel/hooks"
	_ "github.com/deckhouse/deckhouse/modules/035-cni-simple-bridge/hooks"
	_ "github.com/deckhouse/deckhouse/modules/040-control-plane-manager/hooks"
	_ "github.com/deckhouse/deckhouse/modules/040-control-plane-manager/requirements"
	_ "github.com/deckhouse/deckhouse/modules/040-node-manager/hooks"
	_ "github.com/deckhouse/deckhouse/modules/040-node-manager/hooks/migration"
	_ "github.com/deckhouse/deckhouse/modules/040-node-manager/hooks/pkg/schema"
	_ "github.com/deckhouse/deckhouse/modules/040-node-manager/requirements"
	_ "github.com/deckhouse/deckhouse/modules/042-kube-dns/hooks"
	_ "github.com/deckhouse/deckhouse/modules/045-snapshot-controller/hooks"
	_ "github.com/deckhouse/deckhouse/modules/101-cert-manager/hooks"
	_ "github.com/deckhouse/deckhouse/modules/110-istio/hooks"
	_ "github.com/deckhouse/deckhouse/modules/110-istio/hooks/https"
	_ "github.com/deckhouse/deckhouse/modules/110-istio/hooks/lib"
	_ "github.com/deckhouse/deckhouse/modules/110-istio/hooks/lib/crd"
	_ "github.com/deckhouse/deckhouse/modules/110-istio/hooks/lib/istio_versions"
	_ "github.com/deckhouse/deckhouse/modules/110-istio/requirements"
	_ "github.com/deckhouse/deckhouse/modules/140-user-authz/hooks"
	_ "github.com/deckhouse/deckhouse/modules/150-user-authn/hooks"
	_ "github.com/deckhouse/deckhouse/modules/150-user-authn/hooks/https"
	_ "github.com/deckhouse/deckhouse/modules/300-prometheus/hooks"
	_ "github.com/deckhouse/deckhouse/modules/300-prometheus/hooks/https"
	_ "github.com/deckhouse/deckhouse/modules/301-prometheus-metrics-adapter/hooks"
	_ "github.com/deckhouse/deckhouse/modules/302-vertical-pod-autoscaler/hooks"
	_ "github.com/deckhouse/deckhouse/modules/340-extended-monitoring/hooks"
	_ "github.com/deckhouse/deckhouse/modules/340-monitoring-custom/hooks"
	_ "github.com/deckhouse/deckhouse/modules/340-monitoring-kubernetes/hooks"
	_ "github.com/deckhouse/deckhouse/modules/340-monitoring-kubernetes/requirements"
	_ "github.com/deckhouse/deckhouse/modules/340-monitoring-ping/hooks"
	_ "github.com/deckhouse/deckhouse/modules/400-descheduler/hooks"
	_ "github.com/deckhouse/deckhouse/modules/402-ingress-nginx/hooks"
	_ "github.com/deckhouse/deckhouse/modules/402-ingress-nginx/requirements"
	_ "github.com/deckhouse/deckhouse/modules/460-log-shipper/hooks"
	_ "github.com/deckhouse/deckhouse/modules/462-loki/hooks"
	_ "github.com/deckhouse/deckhouse/modules/500-cilium-hubble/hooks"
	_ "github.com/deckhouse/deckhouse/modules/500-cilium-hubble/hooks/https"
	_ "github.com/deckhouse/deckhouse/modules/500-dashboard/hooks"
	_ "github.com/deckhouse/deckhouse/modules/500-dashboard/hooks/https"
	_ "github.com/deckhouse/deckhouse/modules/500-okmeter/hooks"
	_ "github.com/deckhouse/deckhouse/modules/500-openvpn/hooks"
	_ "github.com/deckhouse/deckhouse/modules/500-openvpn/hooks/https"
	_ "github.com/deckhouse/deckhouse/modules/500-upmeter/hooks"
	_ "github.com/deckhouse/deckhouse/modules/500-upmeter/hooks/dynamic_probe"
	_ "github.com/deckhouse/deckhouse/modules/500-upmeter/hooks/https"
	_ "github.com/deckhouse/deckhouse/modules/500-upmeter/hooks/smokemini"
	_ "github.com/deckhouse/deckhouse/modules/600-namespace-configurator/hooks"
	_ "github.com/deckhouse/deckhouse/modules/600-secret-copier/hooks"
	_ "github.com/deckhouse/deckhouse/modules/810-documentation/hooks"
	_ "github.com/deckhouse/deckhouse/modules/810-documentation/hooks/https"
)
