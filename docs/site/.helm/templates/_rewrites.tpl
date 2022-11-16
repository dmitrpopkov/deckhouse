{{- define "rewrites" }}
rewrite ^/$ /en/ permanent;
rewrite ^/(en|ru)/getting_started.html$ /$1/gs/ permanent;
rewrite ^/ru/terms-of-service\.html /ru/security-policy.html permanent;
rewrite ^/ru/cookie-policy\.html /ru/security-policy.html permanent;
rewrite ^/ru/privacy-policy\.html /ru/security-policy.html permanent;
rewrite ^/en/security-policy\.html /en/privacy-policy.html permanent;
rewrite ^/(en|ru)/documentation/v1/modules/030-cloud-provider-openstack/usage.html$ /$1/documentation/v1/modules/030-cloud-provider-openstack/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/030-cloud-provider-vsphere/usage.html$ /$1/documentation/v1/modules/030-cloud-provider-vsphere/docs/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/030-cloud-provider-aws/usage.html$ /$1/documentation/v1/modules/030-cloud-provider-aws/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/030-cloud-provider-azure/usage.html$ /$1/documentation/v1/modules/030-cloud-provider-azure/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/030-cloud-provider-gcp/usage.html$ /$1/documentation/v1/modules/030-cloud-provider-gcp/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/030-cloud-provider-yandex/usage.html$ /$1/documentation/v1/modules/030-cloud-provider-yandex/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/031-local-path-provisioner/usage.html$ /$1/documentation/v1/modules/031-local-path-provisioner/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/040-control-plane-manager/usage.html$ /$1/documentation/v1/modules/040-control-plane-manager/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/040-node-manager/usage.html$ /$1/documentation/v1/modules/040-node-manager/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/042-kube-dns/usage.html$ /$1/documentation/v1/modules/042-kube-dns/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/050-network-policy-engine/usage.html$ /$1/documentation/v1/modules/050-network-policy-engine/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/099-ceph-csi/usage.html$ /$1/documentation/v1/modules/099-ceph-csi/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/110-istio/usage.html$ /$1/documentation/v1/modules/110-istio/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/302-vertical-pod-autoscaler/usage.html$ /$1/documentation/v1/modules/302-vertical-pod-autoscaler/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/303-prometheus-pushgateway/usage.html$ /$1/documentation/v1/modules/303-prometheus-pushgateway/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/350-node-local-dns/usage.html$ /$1/documentation/v1/modules/350-node-local-dns/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/380-metallb/usage.html$ /$1/documentation/v1/modules/380-metallb/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/400-descheduler/usage.html$ /$1/documentation/v1/modules/400-descheduler/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/402-ingress-nginx/usage.html$ /$1/documentation/v1/modules/402-ingress-nginx/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/450-keepalived/usage.html$ /$1/documentation/v1/modules/450-keepalived/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/460-log-shipper/usage.html$ /$1/documentation/v1/modules/460-log-shipper/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/465-pod-reloader/usage.html$ /$1/documentation/v1/modules/465-pod-reloader/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/500-basic-auth/usage.html$ /$1/documentation/v1/modules/500-basic-auth/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/500-dashboard/usage.html$ /$1/documentation/v1/modules/500-dashboard/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/500-okmeter/usage.html$ /$1/documentation/v1/modules/500-okmeter/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/500-openvpn/usage.html$ /$1/documentation/v1/modules/500-openvpn/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/500-upmeter/usage.html$ /$1/documentation/v1/modules/500-upmeter/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/600-namespace-configurator/usage.html$ /$1/documentation/v1/modules/600-namespace-configurator/examples.html;
rewrite ^/(en|ru)/documentation/v1/modules/810-deckhouse-web/usage.html$ /$1/documentation/v1/modules/810-deckhouse-web/examples.html;
{{- end }}
