// Code generated by "tools/images_tags.go" DO NOT EDIT.
// To generate run 'make generate'
package library

var DefaultImagesDigests = map[string]interface{}{
	"admissionPolicyEngine": map[string]interface{}{
		"constraintExporter": "imageHash-admissionPolicyEngine-constraintExporter",
		"gatekeeper":         "imageHash-admissionPolicyEngine-gatekeeper",
	},
	"basicAuth": map[string]interface{}{
		"nginx": "imageHash-basicAuth-nginx",
	},
	"cephCsi": map[string]interface{}{
		"cephcsi": "imageHash-cephCsi-cephcsi",
	},
	"certManager": map[string]interface{}{
		"certManagerAcmeSolver": "imageHash-certManager-certManagerAcmeSolver",
		"certManagerCainjector": "imageHash-certManager-certManagerCainjector",
		"certManagerController": "imageHash-certManager-certManagerController",
		"certManagerWebhook":    "imageHash-certManager-certManagerWebhook",
	},
	"chrony": map[string]interface{}{
		"chrony": "imageHash-chrony-chrony",
	},
	"ciliumHubble": map[string]interface{}{
		"relay":      "imageHash-ciliumHubble-relay",
		"uiBackend":  "imageHash-ciliumHubble-uiBackend",
		"uiFrontend": "imageHash-ciliumHubble-uiFrontend",
	},
	"cloudProviderAws": map[string]interface{}{
		"cloudControllerManager122": "imageHash-cloudProviderAws-cloudControllerManager122",
		"cloudControllerManager123": "imageHash-cloudProviderAws-cloudControllerManager123",
		"cloudControllerManager124": "imageHash-cloudProviderAws-cloudControllerManager124",
		"cloudControllerManager125": "imageHash-cloudProviderAws-cloudControllerManager125",
		"cloudControllerManager126": "imageHash-cloudProviderAws-cloudControllerManager126",
		"ebsCsiPlugin":              "imageHash-cloudProviderAws-ebsCsiPlugin",
		"nodeTerminationHandler":    "imageHash-cloudProviderAws-nodeTerminationHandler",
	},
	"cloudProviderAzure": map[string]interface{}{
		"azurediskCsi":              "imageHash-cloudProviderAzure-azurediskCsi",
		"cloudControllerManager122": "imageHash-cloudProviderAzure-cloudControllerManager122",
		"cloudControllerManager123": "imageHash-cloudProviderAzure-cloudControllerManager123",
		"cloudControllerManager124": "imageHash-cloudProviderAzure-cloudControllerManager124",
		"cloudControllerManager125": "imageHash-cloudProviderAzure-cloudControllerManager125",
		"cloudControllerManager126": "imageHash-cloudProviderAzure-cloudControllerManager126",
	},
	"cloudProviderGcp": map[string]interface{}{
		"cloudControllerManager122": "imageHash-cloudProviderGcp-cloudControllerManager122",
		"cloudControllerManager123": "imageHash-cloudProviderGcp-cloudControllerManager123",
		"cloudControllerManager124": "imageHash-cloudProviderGcp-cloudControllerManager124",
		"cloudControllerManager125": "imageHash-cloudProviderGcp-cloudControllerManager125",
		"cloudControllerManager126": "imageHash-cloudProviderGcp-cloudControllerManager126",
		"pdCsiPlugin":               "imageHash-cloudProviderGcp-pdCsiPlugin",
	},
	"cloudProviderOpenstack": map[string]interface{}{
		"cinderCsiPlugin122":        "imageHash-cloudProviderOpenstack-cinderCsiPlugin122",
		"cinderCsiPlugin123":        "imageHash-cloudProviderOpenstack-cinderCsiPlugin123",
		"cinderCsiPlugin124":        "imageHash-cloudProviderOpenstack-cinderCsiPlugin124",
		"cinderCsiPlugin125":        "imageHash-cloudProviderOpenstack-cinderCsiPlugin125",
		"cinderCsiPlugin126":        "imageHash-cloudProviderOpenstack-cinderCsiPlugin126",
		"cloudControllerManager122": "imageHash-cloudProviderOpenstack-cloudControllerManager122",
		"cloudControllerManager123": "imageHash-cloudProviderOpenstack-cloudControllerManager123",
		"cloudControllerManager124": "imageHash-cloudProviderOpenstack-cloudControllerManager124",
		"cloudControllerManager125": "imageHash-cloudProviderOpenstack-cloudControllerManager125",
		"cloudControllerManager126": "imageHash-cloudProviderOpenstack-cloudControllerManager126",
	},
	"cloudProviderVsphere": map[string]interface{}{
		"cloudControllerManager122": "imageHash-cloudProviderVsphere-cloudControllerManager122",
		"cloudControllerManager123": "imageHash-cloudProviderVsphere-cloudControllerManager123",
		"cloudControllerManager124": "imageHash-cloudProviderVsphere-cloudControllerManager124",
		"cloudControllerManager125": "imageHash-cloudProviderVsphere-cloudControllerManager125",
		"cloudControllerManager126": "imageHash-cloudProviderVsphere-cloudControllerManager126",
		"vsphereCsiPlugin":          "imageHash-cloudProviderVsphere-vsphereCsiPlugin",
		"vsphereCsiPluginLegacy":    "imageHash-cloudProviderVsphere-vsphereCsiPluginLegacy",
	},
	"cloudProviderYandex": map[string]interface{}{
		"cloudControllerManager122": "imageHash-cloudProviderYandex-cloudControllerManager122",
		"cloudControllerManager123": "imageHash-cloudProviderYandex-cloudControllerManager123",
		"cloudControllerManager124": "imageHash-cloudProviderYandex-cloudControllerManager124",
		"cloudControllerManager125": "imageHash-cloudProviderYandex-cloudControllerManager125",
		"cloudControllerManager126": "imageHash-cloudProviderYandex-cloudControllerManager126",
		"cloudMetricsExporter":      "imageHash-cloudProviderYandex-cloudMetricsExporter",
		"yandexCsiPlugin":           "imageHash-cloudProviderYandex-yandexCsiPlugin",
	},
	"cniCilium": map[string]interface{}{
		"cilium":     "imageHash-cniCilium-cilium",
		"operator":   "imageHash-cniCilium-operator",
		"virtCilium": "imageHash-cniCilium-virtCilium",
	},
	"cniFlannel": map[string]interface{}{
		"flanneld": "imageHash-cniFlannel-flanneld",
	},
	"cniSimpleBridge": map[string]interface{}{
		"simpleBridge": "imageHash-cniSimpleBridge-simpleBridge",
	},
	"common": map[string]interface{}{
		"alpine":                    "imageHash-common-alpine",
		"checkKernelVersion":        "imageHash-common-checkKernelVersion",
		"csiExternalAttacher122":    "imageHash-common-csiExternalAttacher122",
		"csiExternalAttacher123":    "imageHash-common-csiExternalAttacher123",
		"csiExternalAttacher124":    "imageHash-common-csiExternalAttacher124",
		"csiExternalAttacher125":    "imageHash-common-csiExternalAttacher125",
		"csiExternalAttacher126":    "imageHash-common-csiExternalAttacher126",
		"csiExternalProvisioner122": "imageHash-common-csiExternalProvisioner122",
		"csiExternalProvisioner123": "imageHash-common-csiExternalProvisioner123",
		"csiExternalProvisioner124": "imageHash-common-csiExternalProvisioner124",
		"csiExternalProvisioner125": "imageHash-common-csiExternalProvisioner125",
		"csiExternalProvisioner126": "imageHash-common-csiExternalProvisioner126",
		"csiExternalResizer122":     "imageHash-common-csiExternalResizer122",
		"csiExternalResizer123":     "imageHash-common-csiExternalResizer123",
		"csiExternalResizer124":     "imageHash-common-csiExternalResizer124",
		"csiExternalResizer125":     "imageHash-common-csiExternalResizer125",
		"csiExternalResizer126":     "imageHash-common-csiExternalResizer126",
		"csiExternalSnapshotter122": "imageHash-common-csiExternalSnapshotter122",
		"csiExternalSnapshotter123": "imageHash-common-csiExternalSnapshotter123",
		"csiExternalSnapshotter124": "imageHash-common-csiExternalSnapshotter124",
		"csiExternalSnapshotter125": "imageHash-common-csiExternalSnapshotter125",
		"csiExternalSnapshotter126": "imageHash-common-csiExternalSnapshotter126",
		"csiLivenessprobe122":       "imageHash-common-csiLivenessprobe122",
		"csiLivenessprobe123":       "imageHash-common-csiLivenessprobe123",
		"csiLivenessprobe124":       "imageHash-common-csiLivenessprobe124",
		"csiLivenessprobe125":       "imageHash-common-csiLivenessprobe125",
		"csiLivenessprobe126":       "imageHash-common-csiLivenessprobe126",
		"csiNodeDriverRegistrar122": "imageHash-common-csiNodeDriverRegistrar122",
		"csiNodeDriverRegistrar123": "imageHash-common-csiNodeDriverRegistrar123",
		"csiNodeDriverRegistrar124": "imageHash-common-csiNodeDriverRegistrar124",
		"csiNodeDriverRegistrar125": "imageHash-common-csiNodeDriverRegistrar125",
		"csiNodeDriverRegistrar126": "imageHash-common-csiNodeDriverRegistrar126",
		"kubeRbacProxy":             "imageHash-common-kubeRbacProxy",
		"pause":                     "imageHash-common-pause",
	},
	"containerizedDataImporter": map[string]interface{}{
		"cdiApiserver":    "imageHash-containerizedDataImporter-cdiApiserver",
		"cdiCloner":       "imageHash-containerizedDataImporter-cdiCloner",
		"cdiController":   "imageHash-containerizedDataImporter-cdiController",
		"cdiImporter":     "imageHash-containerizedDataImporter-cdiImporter",
		"cdiOperator":     "imageHash-containerizedDataImporter-cdiOperator",
		"cdiUploadproxy":  "imageHash-containerizedDataImporter-cdiUploadproxy",
		"cdiUploadserver": "imageHash-containerizedDataImporter-cdiUploadserver",
	},
	"controlPlaneManager": map[string]interface{}{
		"controlPlaneManager122":   "imageHash-controlPlaneManager-controlPlaneManager122",
		"controlPlaneManager123":   "imageHash-controlPlaneManager-controlPlaneManager123",
		"controlPlaneManager124":   "imageHash-controlPlaneManager-controlPlaneManager124",
		"controlPlaneManager125":   "imageHash-controlPlaneManager-controlPlaneManager125",
		"controlPlaneManager126":   "imageHash-controlPlaneManager-controlPlaneManager126",
		"etcd":                     "imageHash-controlPlaneManager-etcd",
		"kubeApiserver122":         "imageHash-controlPlaneManager-kubeApiserver122",
		"kubeApiserver123":         "imageHash-controlPlaneManager-kubeApiserver123",
		"kubeApiserver124":         "imageHash-controlPlaneManager-kubeApiserver124",
		"kubeApiserver125":         "imageHash-controlPlaneManager-kubeApiserver125",
		"kubeApiserver126":         "imageHash-controlPlaneManager-kubeApiserver126",
		"kubeApiserverHealthcheck": "imageHash-controlPlaneManager-kubeApiserverHealthcheck",
		"kubeControllerManager122": "imageHash-controlPlaneManager-kubeControllerManager122",
		"kubeControllerManager123": "imageHash-controlPlaneManager-kubeControllerManager123",
		"kubeControllerManager124": "imageHash-controlPlaneManager-kubeControllerManager124",
		"kubeControllerManager125": "imageHash-controlPlaneManager-kubeControllerManager125",
		"kubeControllerManager126": "imageHash-controlPlaneManager-kubeControllerManager126",
		"kubeScheduler122":         "imageHash-controlPlaneManager-kubeScheduler122",
		"kubeScheduler123":         "imageHash-controlPlaneManager-kubeScheduler123",
		"kubeScheduler124":         "imageHash-controlPlaneManager-kubeScheduler124",
		"kubeScheduler125":         "imageHash-controlPlaneManager-kubeScheduler125",
		"kubeScheduler126":         "imageHash-controlPlaneManager-kubeScheduler126",
		"kubernetesApiProxy":       "imageHash-controlPlaneManager-kubernetesApiProxy",
	},
	"dashboard": map[string]interface{}{
		"dashboard":      "imageHash-dashboard-dashboard",
		"metricsScraper": "imageHash-dashboard-metricsScraper",
	},
	"deckhouse": map[string]interface{}{
		"imagesCopier":   "imageHash-deckhouse-imagesCopier",
		"webhookHandler": "imageHash-deckhouse-webhookHandler",
	},
	"deckhouseConfig": map[string]interface{}{
		"deckhouseConfigWebhook": "imageHash-deckhouseConfig-deckhouseConfigWebhook",
	},
	"deckhouseWeb": map[string]interface{}{
		"web": "imageHash-deckhouseWeb-web",
	},
	"delivery": map[string]interface{}{
		"argocd":               "imageHash-delivery-argocd",
		"argocdImageUpdater":   "imageHash-delivery-argocdImageUpdater",
		"redis":                "imageHash-delivery-redis",
		"werfArgocdCmpSidecar": "imageHash-delivery-werfArgocdCmpSidecar",
	},
	"descheduler": map[string]interface{}{
		"descheduler": "imageHash-descheduler-descheduler",
	},
	"extendedMonitoring": map[string]interface{}{
		"certExporter":               "imageHash-extendedMonitoring-certExporter",
		"eventsExporter":             "imageHash-extendedMonitoring-eventsExporter",
		"extendedMonitoringExporter": "imageHash-extendedMonitoring-extendedMonitoringExporter",
		"imageAvailabilityExporter":  "imageHash-extendedMonitoring-imageAvailabilityExporter",
	},
	"flantIntegration": map[string]interface{}{
		"flantPricing": "imageHash-flantIntegration-flantPricing",
		"grafanaAgent": "imageHash-flantIntegration-grafanaAgent",
		"madisonProxy": "imageHash-flantIntegration-madisonProxy",
	},
	"ingressNginx": map[string]interface{}{
		"controller11":          "imageHash-ingressNginx-controller11",
		"controller16":          "imageHash-ingressNginx-controller16",
		"kruise":                "imageHash-ingressNginx-kruise",
		"nginxExporter":         "imageHash-ingressNginx-nginxExporter",
		"protobufExporter":      "imageHash-ingressNginx-protobufExporter",
		"proxyFailover":         "imageHash-ingressNginx-proxyFailover",
		"proxyFailoverIptables": "imageHash-ingressNginx-proxyFailoverIptables",
	},
	"istio": map[string]interface{}{
		"apiProxy":          "imageHash-istio-apiProxy",
		"kiali":             "imageHash-istio-kiali",
		"metadataDiscovery": "imageHash-istio-metadataDiscovery",
		"metadataExporter":  "imageHash-istio-metadataExporter",
		"operatorV1x12x6":   "imageHash-istio-operatorV1x12x6",
		"operatorV1x13x7":   "imageHash-istio-operatorV1x13x7",
		"operatorV1x16x2":   "imageHash-istio-operatorV1x16x2",
		"pilotV1x12x6":      "imageHash-istio-pilotV1x12x6",
		"pilotV1x13x7":      "imageHash-istio-pilotV1x13x7",
		"pilotV1x16x2":      "imageHash-istio-pilotV1x16x2",
		"proxyv2V1x12x6":    "imageHash-istio-proxyv2V1x12x6",
		"proxyv2V1x13x7":    "imageHash-istio-proxyv2V1x13x7",
		"proxyv2V1x16x2":    "imageHash-istio-proxyv2V1x16x2",
	},
	"keepalived": map[string]interface{}{
		"keepalived": "imageHash-keepalived-keepalived",
	},
	"kubeDns": map[string]interface{}{
		"coredns":                           "imageHash-kubeDns-coredns",
		"resolvWatcher":                     "imageHash-kubeDns-resolvWatcher",
		"stsPodsHostsAppenderInitContainer": "imageHash-kubeDns-stsPodsHostsAppenderInitContainer",
		"stsPodsHostsAppenderWebhook":       "imageHash-kubeDns-stsPodsHostsAppenderWebhook",
	},
	"kubeProxy": map[string]interface{}{
		"initContainer": "imageHash-kubeProxy-initContainer",
		"kubeProxy122":  "imageHash-kubeProxy-kubeProxy122",
		"kubeProxy123":  "imageHash-kubeProxy-kubeProxy123",
		"kubeProxy124":  "imageHash-kubeProxy-kubeProxy124",
		"kubeProxy125":  "imageHash-kubeProxy-kubeProxy125",
		"kubeProxy126":  "imageHash-kubeProxy-kubeProxy126",
	},
	"linstor": map[string]interface{}{
		"drbdDriverLoader":          "imageHash-linstor-drbdDriverLoader",
		"drbdReactor":               "imageHash-linstor-drbdReactor",
		"linstorAffinityController": "imageHash-linstor-linstorAffinityController",
		"linstorCsi":                "imageHash-linstor-linstorCsi",
		"linstorPoolsImporter":      "imageHash-linstor-linstorPoolsImporter",
		"linstorSchedulerAdmission": "imageHash-linstor-linstorSchedulerAdmission",
		"linstorSchedulerExtender":  "imageHash-linstor-linstorSchedulerExtender",
		"linstorServer":             "imageHash-linstor-linstorServer",
		"piraeusHaController":       "imageHash-linstor-piraeusHaController",
		"piraeusOperator":           "imageHash-linstor-piraeusOperator",
		"spaas":                     "imageHash-linstor-spaas",
	},
	"localPathProvisioner": map[string]interface{}{
		"helper":               "imageHash-localPathProvisioner-helper",
		"localPathProvisioner": "imageHash-localPathProvisioner-localPathProvisioner",
	},
	"logShipper": map[string]interface{}{
		"vector": "imageHash-logShipper-vector",
	},
	"loki": map[string]interface{}{
		"loki": "imageHash-loki-loki",
	},
	"metallb": map[string]interface{}{
		"controller": "imageHash-metallb-controller",
		"speaker":    "imageHash-metallb-speaker",
	},
	"monitoringKubernetes": map[string]interface{}{
		"ebpfExporter":                      "imageHash-monitoringKubernetes-ebpfExporter",
		"kubeStateMetrics":                  "imageHash-monitoringKubernetes-kubeStateMetrics",
		"kubeletEvictionThresholdsExporter": "imageHash-monitoringKubernetes-kubeletEvictionThresholdsExporter",
		"nodeExporter":                      "imageHash-monitoringKubernetes-nodeExporter",
	},
	"monitoringPing": map[string]interface{}{
		"monitoringPing": "imageHash-monitoringPing-monitoringPing",
	},
	"networkGateway": map[string]interface{}{
		"dnsmasq": "imageHash-networkGateway-dnsmasq",
		"snat":    "imageHash-networkGateway-snat",
	},
	"networkPolicyEngine": map[string]interface{}{
		"kubeRouter": "imageHash-networkPolicyEngine-kubeRouter",
	},
	"nodeLocalDns": map[string]interface{}{
		"coredns":      "imageHash-nodeLocalDns-coredns",
		"iptablesLoop": "imageHash-nodeLocalDns-iptablesLoop",
	},
	"nodeManager": map[string]interface{}{
		"bashibleApiserver":        "imageHash-nodeManager-bashibleApiserver",
		"clusterAutoscaler":        "imageHash-nodeManager-clusterAutoscaler",
		"earlyOom":                 "imageHash-nodeManager-earlyOom",
		"machineControllerManager": "imageHash-nodeManager-machineControllerManager",
	},
	"openvpn": map[string]interface{}{
		"easyrsaMigrator": "imageHash-openvpn-easyrsaMigrator",
		"openvpn":         "imageHash-openvpn-openvpn",
		"ovpnAdmin":       "imageHash-openvpn-ovpnAdmin",
		"pmacct":          "imageHash-openvpn-pmacct",
	},
	"operatorPrometheus": map[string]interface{}{
		"prometheusConfigReloader": "imageHash-operatorPrometheus-prometheusConfigReloader",
		"prometheusOperator":       "imageHash-operatorPrometheus-prometheusOperator",
	},
	"operatorTrivy": map[string]interface{}{
		"nodeCollector": "imageHash-operatorTrivy-nodeCollector",
		"operator":      "imageHash-operatorTrivy-operator",
		"trivy":         "imageHash-operatorTrivy-trivy",
	},
	"podReloader": map[string]interface{}{
		"podReloader": "imageHash-podReloader-podReloader",
	},
	"prometheus": map[string]interface{}{
		"alertmanager":                "imageHash-prometheus-alertmanager",
		"alertsReceiverWebhook":       "imageHash-prometheus-alertsReceiverWebhook",
		"grafana":                     "imageHash-prometheus-grafana",
		"grafanaDashboardProvisioner": "imageHash-prometheus-grafanaDashboardProvisioner",
		"prometheus":                  "imageHash-prometheus-prometheus",
		"trickster":                   "imageHash-prometheus-trickster",
	},
	"prometheusMetricsAdapter": map[string]interface{}{
		"k8sPrometheusAdapter":   "imageHash-prometheusMetricsAdapter-k8sPrometheusAdapter",
		"prometheusReverseProxy": "imageHash-prometheusMetricsAdapter-prometheusReverseProxy",
	},
	"prometheusPushgateway": map[string]interface{}{
		"pushgateway": "imageHash-prometheusPushgateway-pushgateway",
	},
	"registrypackages": map[string]interface{}{
		"conntrackDebian1462":                 "imageHash-registrypackages-conntrackDebian1462",
		"containerdAlteros14631El7X8664":      "imageHash-registrypackages-containerdAlteros14631El7X8664",
		"containerdAlteros161831El7X8664":     "imageHash-registrypackages-containerdAlteros161831El7X8664",
		"containerdAltlinux1619":              "imageHash-registrypackages-containerdAltlinux1619",
		"containerdAstra1461Buster":           "imageHash-registrypackages-containerdAstra1461Buster",
		"containerdAstra16181Buster":          "imageHash-registrypackages-containerdAstra16181Buster",
		"containerdCentos14631El7X8664":       "imageHash-registrypackages-containerdCentos14631El7X8664",
		"containerdCentos14631El8X8664":       "imageHash-registrypackages-containerdCentos14631El8X8664",
		"containerdCentos161831El7X8664":      "imageHash-registrypackages-containerdCentos161831El7X8664",
		"containerdCentos161831El8X8664":      "imageHash-registrypackages-containerdCentos161831El8X8664",
		"containerdCentos161831El9X8664":      "imageHash-registrypackages-containerdCentos161831El9X8664",
		"containerdCentos16731El9X8664":       "imageHash-registrypackages-containerdCentos16731El9X8664",
		"containerdDebian1431Stretch":         "imageHash-registrypackages-containerdDebian1431Stretch",
		"containerdDebian1461Bullseye":        "imageHash-registrypackages-containerdDebian1461Bullseye",
		"containerdDebian1461Buster":          "imageHash-registrypackages-containerdDebian1461Buster",
		"containerdDebian16181Bullseye":       "imageHash-registrypackages-containerdDebian16181Bullseye",
		"containerdDebian16181Buster":         "imageHash-registrypackages-containerdDebian16181Buster",
		"containerdRedos14631El7X8664":        "imageHash-registrypackages-containerdRedos14631El7X8664",
		"containerdRedos161831El7X8664":       "imageHash-registrypackages-containerdRedos161831El7X8664",
		"containerdUbuntu1461Bionic":          "imageHash-registrypackages-containerdUbuntu1461Bionic",
		"containerdUbuntu1461Focal":           "imageHash-registrypackages-containerdUbuntu1461Focal",
		"containerdUbuntu15111Jammy":          "imageHash-registrypackages-containerdUbuntu15111Jammy",
		"containerdUbuntu16181Bionic":         "imageHash-registrypackages-containerdUbuntu16181Bionic",
		"containerdUbuntu16181Focal":          "imageHash-registrypackages-containerdUbuntu16181Focal",
		"containerdUbuntu16181Jammy":          "imageHash-registrypackages-containerdUbuntu16181Jammy",
		"crictl122":                           "imageHash-registrypackages-crictl122",
		"crictl123":                           "imageHash-registrypackages-crictl123",
		"crictl124":                           "imageHash-registrypackages-crictl124",
		"crictl125":                           "imageHash-registrypackages-crictl125",
		"crictl126":                           "imageHash-registrypackages-crictl126",
		"d8Curl7800":                          "imageHash-registrypackages-d8Curl7800",
		"dockerAlteros1903153El7X866473":      "imageHash-registrypackages-dockerAlteros1903153El7X866473",
		"dockerAltlinux201021Alt1X8664":       "imageHash-registrypackages-dockerAltlinux201021Alt1X8664",
		"dockerAstra520101230DebianBuster":    "imageHash-registrypackages-dockerAstra520101230DebianBuster",
		"dockerCentos1903153El7X86647":        "imageHash-registrypackages-dockerCentos1903153El7X86647",
		"dockerCentos1903153El8X86648":        "imageHash-registrypackages-dockerCentos1903153El8X86648",
		"dockerCentos2010173El9X86649":        "imageHash-registrypackages-dockerCentos2010173El9X86649",
		"dockerDebian519031530DebianStretch":  "imageHash-registrypackages-dockerDebian519031530DebianStretch",
		"dockerDebian520101230DebianBullseye": "imageHash-registrypackages-dockerDebian520101230DebianBullseye",
		"dockerDebian520101230DebianBuster":   "imageHash-registrypackages-dockerDebian520101230DebianBuster",
		"dockerRedos1903153El7X866473":        "imageHash-registrypackages-dockerRedos1903153El7X866473",
		"dockerUbuntu519031330UbuntuBionic":   "imageHash-registrypackages-dockerUbuntu519031330UbuntuBionic",
		"dockerUbuntu519031330UbuntuFocal":    "imageHash-registrypackages-dockerUbuntu519031330UbuntuFocal",
		"dockerUbuntu520101430UbuntuJammy":    "imageHash-registrypackages-dockerUbuntu520101430UbuntuJammy",
		"inotifyToolsCentos73149":             "imageHash-registrypackages-inotifyToolsCentos73149",
		"inotifyToolsCentos831419":            "imageHash-registrypackages-inotifyToolsCentos831419",
		"inotifyToolsCentos9322101":           "imageHash-registrypackages-inotifyToolsCentos9322101",
		"jq16":                                "imageHash-registrypackages-jq16",
		"kubeadmAlteros12217":                 "imageHash-registrypackages-kubeadmAlteros12217",
		"kubeadmAlteros12317":                 "imageHash-registrypackages-kubeadmAlteros12317",
		"kubeadmAlteros12412":                 "imageHash-registrypackages-kubeadmAlteros12412",
		"kubeadmAlteros1258":                  "imageHash-registrypackages-kubeadmAlteros1258",
		"kubeadmAlteros1263":                  "imageHash-registrypackages-kubeadmAlteros1263",
		"kubeadmAltlinux12217":                "imageHash-registrypackages-kubeadmAltlinux12217",
		"kubeadmAltlinux12317":                "imageHash-registrypackages-kubeadmAltlinux12317",
		"kubeadmAltlinux12412":                "imageHash-registrypackages-kubeadmAltlinux12412",
		"kubeadmAltlinux1258":                 "imageHash-registrypackages-kubeadmAltlinux1258",
		"kubeadmAltlinux1263":                 "imageHash-registrypackages-kubeadmAltlinux1263",
		"kubeadmAstra12217":                   "imageHash-registrypackages-kubeadmAstra12217",
		"kubeadmAstra12317":                   "imageHash-registrypackages-kubeadmAstra12317",
		"kubeadmAstra12412":                   "imageHash-registrypackages-kubeadmAstra12412",
		"kubeadmAstra1258":                    "imageHash-registrypackages-kubeadmAstra1258",
		"kubeadmAstra1263":                    "imageHash-registrypackages-kubeadmAstra1263",
		"kubeadmCentos12217":                  "imageHash-registrypackages-kubeadmCentos12217",
		"kubeadmCentos12317":                  "imageHash-registrypackages-kubeadmCentos12317",
		"kubeadmCentos12412":                  "imageHash-registrypackages-kubeadmCentos12412",
		"kubeadmCentos1258":                   "imageHash-registrypackages-kubeadmCentos1258",
		"kubeadmCentos1263":                   "imageHash-registrypackages-kubeadmCentos1263",
		"kubeadmDebian12217":                  "imageHash-registrypackages-kubeadmDebian12217",
		"kubeadmDebian12317":                  "imageHash-registrypackages-kubeadmDebian12317",
		"kubeadmDebian12412":                  "imageHash-registrypackages-kubeadmDebian12412",
		"kubeadmDebian1258":                   "imageHash-registrypackages-kubeadmDebian1258",
		"kubeadmDebian1263":                   "imageHash-registrypackages-kubeadmDebian1263",
		"kubeadmRedos12217":                   "imageHash-registrypackages-kubeadmRedos12217",
		"kubeadmRedos12317":                   "imageHash-registrypackages-kubeadmRedos12317",
		"kubeadmRedos12412":                   "imageHash-registrypackages-kubeadmRedos12412",
		"kubeadmRedos1258":                    "imageHash-registrypackages-kubeadmRedos1258",
		"kubeadmRedos1263":                    "imageHash-registrypackages-kubeadmRedos1263",
		"kubeadmUbuntu12217":                  "imageHash-registrypackages-kubeadmUbuntu12217",
		"kubeadmUbuntu12317":                  "imageHash-registrypackages-kubeadmUbuntu12317",
		"kubeadmUbuntu12412":                  "imageHash-registrypackages-kubeadmUbuntu12412",
		"kubeadmUbuntu1258":                   "imageHash-registrypackages-kubeadmUbuntu1258",
		"kubeadmUbuntu1263":                   "imageHash-registrypackages-kubeadmUbuntu1263",
		"kubectlAlteros12217":                 "imageHash-registrypackages-kubectlAlteros12217",
		"kubectlAlteros12317":                 "imageHash-registrypackages-kubectlAlteros12317",
		"kubectlAlteros12412":                 "imageHash-registrypackages-kubectlAlteros12412",
		"kubectlAlteros1258":                  "imageHash-registrypackages-kubectlAlteros1258",
		"kubectlAlteros1263":                  "imageHash-registrypackages-kubectlAlteros1263",
		"kubectlAltlinux12217":                "imageHash-registrypackages-kubectlAltlinux12217",
		"kubectlAltlinux12317":                "imageHash-registrypackages-kubectlAltlinux12317",
		"kubectlAltlinux12412":                "imageHash-registrypackages-kubectlAltlinux12412",
		"kubectlAltlinux1258":                 "imageHash-registrypackages-kubectlAltlinux1258",
		"kubectlAltlinux1263":                 "imageHash-registrypackages-kubectlAltlinux1263",
		"kubectlAstra12217":                   "imageHash-registrypackages-kubectlAstra12217",
		"kubectlAstra12317":                   "imageHash-registrypackages-kubectlAstra12317",
		"kubectlAstra12412":                   "imageHash-registrypackages-kubectlAstra12412",
		"kubectlAstra1258":                    "imageHash-registrypackages-kubectlAstra1258",
		"kubectlAstra1263":                    "imageHash-registrypackages-kubectlAstra1263",
		"kubectlCentos12217":                  "imageHash-registrypackages-kubectlCentos12217",
		"kubectlCentos12317":                  "imageHash-registrypackages-kubectlCentos12317",
		"kubectlCentos12412":                  "imageHash-registrypackages-kubectlCentos12412",
		"kubectlCentos1258":                   "imageHash-registrypackages-kubectlCentos1258",
		"kubectlCentos1263":                   "imageHash-registrypackages-kubectlCentos1263",
		"kubectlDebian12217":                  "imageHash-registrypackages-kubectlDebian12217",
		"kubectlDebian12317":                  "imageHash-registrypackages-kubectlDebian12317",
		"kubectlDebian12412":                  "imageHash-registrypackages-kubectlDebian12412",
		"kubectlDebian1258":                   "imageHash-registrypackages-kubectlDebian1258",
		"kubectlDebian1263":                   "imageHash-registrypackages-kubectlDebian1263",
		"kubectlRedos12217":                   "imageHash-registrypackages-kubectlRedos12217",
		"kubectlRedos12317":                   "imageHash-registrypackages-kubectlRedos12317",
		"kubectlRedos12412":                   "imageHash-registrypackages-kubectlRedos12412",
		"kubectlRedos1258":                    "imageHash-registrypackages-kubectlRedos1258",
		"kubectlRedos1263":                    "imageHash-registrypackages-kubectlRedos1263",
		"kubectlUbuntu12217":                  "imageHash-registrypackages-kubectlUbuntu12217",
		"kubectlUbuntu12317":                  "imageHash-registrypackages-kubectlUbuntu12317",
		"kubectlUbuntu12412":                  "imageHash-registrypackages-kubectlUbuntu12412",
		"kubectlUbuntu1258":                   "imageHash-registrypackages-kubectlUbuntu1258",
		"kubectlUbuntu1263":                   "imageHash-registrypackages-kubectlUbuntu1263",
		"kubeletAlteros12217":                 "imageHash-registrypackages-kubeletAlteros12217",
		"kubeletAlteros12317":                 "imageHash-registrypackages-kubeletAlteros12317",
		"kubeletAlteros12412":                 "imageHash-registrypackages-kubeletAlteros12412",
		"kubeletAlteros1258":                  "imageHash-registrypackages-kubeletAlteros1258",
		"kubeletAlteros1263":                  "imageHash-registrypackages-kubeletAlteros1263",
		"kubeletAltlinux12217":                "imageHash-registrypackages-kubeletAltlinux12217",
		"kubeletAltlinux12317":                "imageHash-registrypackages-kubeletAltlinux12317",
		"kubeletAltlinux12412":                "imageHash-registrypackages-kubeletAltlinux12412",
		"kubeletAltlinux1258":                 "imageHash-registrypackages-kubeletAltlinux1258",
		"kubeletAltlinux1263":                 "imageHash-registrypackages-kubeletAltlinux1263",
		"kubeletAstra12217":                   "imageHash-registrypackages-kubeletAstra12217",
		"kubeletAstra12317":                   "imageHash-registrypackages-kubeletAstra12317",
		"kubeletAstra12412":                   "imageHash-registrypackages-kubeletAstra12412",
		"kubeletAstra1258":                    "imageHash-registrypackages-kubeletAstra1258",
		"kubeletAstra1263":                    "imageHash-registrypackages-kubeletAstra1263",
		"kubeletCentos12217":                  "imageHash-registrypackages-kubeletCentos12217",
		"kubeletCentos12317":                  "imageHash-registrypackages-kubeletCentos12317",
		"kubeletCentos12412":                  "imageHash-registrypackages-kubeletCentos12412",
		"kubeletCentos1258":                   "imageHash-registrypackages-kubeletCentos1258",
		"kubeletCentos1263":                   "imageHash-registrypackages-kubeletCentos1263",
		"kubeletDebian12217":                  "imageHash-registrypackages-kubeletDebian12217",
		"kubeletDebian12317":                  "imageHash-registrypackages-kubeletDebian12317",
		"kubeletDebian12412":                  "imageHash-registrypackages-kubeletDebian12412",
		"kubeletDebian1258":                   "imageHash-registrypackages-kubeletDebian1258",
		"kubeletDebian1263":                   "imageHash-registrypackages-kubeletDebian1263",
		"kubeletRedos12217":                   "imageHash-registrypackages-kubeletRedos12217",
		"kubeletRedos12317":                   "imageHash-registrypackages-kubeletRedos12317",
		"kubeletRedos12412":                   "imageHash-registrypackages-kubeletRedos12412",
		"kubeletRedos1258":                    "imageHash-registrypackages-kubeletRedos1258",
		"kubeletRedos1263":                    "imageHash-registrypackages-kubeletRedos1263",
		"kubeletUbuntu12217":                  "imageHash-registrypackages-kubeletUbuntu12217",
		"kubeletUbuntu12317":                  "imageHash-registrypackages-kubeletUbuntu12317",
		"kubeletUbuntu12412":                  "imageHash-registrypackages-kubeletUbuntu12412",
		"kubeletUbuntu1258":                   "imageHash-registrypackages-kubeletUbuntu1258",
		"kubeletUbuntu1263":                   "imageHash-registrypackages-kubeletUbuntu1263",
		"kubernetesCniAlteros111":             "imageHash-registrypackages-kubernetesCniAlteros111",
		"kubernetesCniAlteros120":             "imageHash-registrypackages-kubernetesCniAlteros120",
		"kubernetesCniAltlinux111":            "imageHash-registrypackages-kubernetesCniAltlinux111",
		"kubernetesCniAltlinux120":            "imageHash-registrypackages-kubernetesCniAltlinux120",
		"kubernetesCniAstra111":               "imageHash-registrypackages-kubernetesCniAstra111",
		"kubernetesCniAstra120":               "imageHash-registrypackages-kubernetesCniAstra120",
		"kubernetesCniCentos111":              "imageHash-registrypackages-kubernetesCniCentos111",
		"kubernetesCniCentos120":              "imageHash-registrypackages-kubernetesCniCentos120",
		"kubernetesCniDebian111":              "imageHash-registrypackages-kubernetesCniDebian111",
		"kubernetesCniDebian120":              "imageHash-registrypackages-kubernetesCniDebian120",
		"kubernetesCniRedos111":               "imageHash-registrypackages-kubernetesCniRedos111",
		"kubernetesCniRedos120":               "imageHash-registrypackages-kubernetesCniRedos120",
		"kubernetesCniUbuntu111":              "imageHash-registrypackages-kubernetesCniUbuntu111",
		"kubernetesCniUbuntu120":              "imageHash-registrypackages-kubernetesCniUbuntu120",
		"tomlMerge01":                         "imageHash-registrypackages-tomlMerge01",
		"virtWhatDebian1151Deb9u1":            "imageHash-registrypackages-virtWhatDebian1151Deb9u1",
	},
	"runtimeAuditEngine": map[string]interface{}{
		"falco":             "imageHash-runtimeAuditEngine-falco",
		"falcoDriverLoader": "imageHash-runtimeAuditEngine-falcoDriverLoader",
		"falcosidekick":     "imageHash-runtimeAuditEngine-falcosidekick",
		"rulesLoader":       "imageHash-runtimeAuditEngine-rulesLoader",
	},
	"snapshotController": map[string]interface{}{
		"snapshotController":        "imageHash-snapshotController-snapshotController",
		"snapshotValidationWebhook": "imageHash-snapshotController-snapshotValidationWebhook",
	},
	"terraformManager": map[string]interface{}{
		"baseTerraformManager":      "imageHash-terraformManager-baseTerraformManager",
		"terraformManagerAws":       "imageHash-terraformManager-terraformManagerAws",
		"terraformManagerAzure":     "imageHash-terraformManager-terraformManagerAzure",
		"terraformManagerGcp":       "imageHash-terraformManager-terraformManagerGcp",
		"terraformManagerOpenstack": "imageHash-terraformManager-terraformManagerOpenstack",
		"terraformManagerVsphere":   "imageHash-terraformManager-terraformManagerVsphere",
		"terraformManagerYandex":    "imageHash-terraformManager-terraformManagerYandex",
	},
	"upmeter": map[string]interface{}{
		"smokeMini": "imageHash-upmeter-smokeMini",
		"status":    "imageHash-upmeter-status",
		"upmeter":   "imageHash-upmeter-upmeter",
		"webui":     "imageHash-upmeter-webui",
	},
	"userAuthn": map[string]interface{}{
		"cfssl":                 "imageHash-userAuthn-cfssl",
		"crowdBasicAuthProxy":   "imageHash-userAuthn-crowdBasicAuthProxy",
		"dex":                   "imageHash-userAuthn-dex",
		"dexAuthenticator":      "imageHash-userAuthn-dexAuthenticator",
		"dexAuthenticatorRedis": "imageHash-userAuthn-dexAuthenticatorRedis",
		"kubeconfigGenerator":   "imageHash-userAuthn-kubeconfigGenerator",
	},
	"userAuthz": map[string]interface{}{
		"webhook": "imageHash-userAuthz-webhook",
	},
	"verticalPodAutoscaler": map[string]interface{}{
		"admissionController": "imageHash-verticalPodAutoscaler-admissionController",
		"recommender":         "imageHash-verticalPodAutoscaler-recommender",
		"updater":             "imageHash-verticalPodAutoscaler-updater",
	},
	"virtualization": map[string]interface{}{
		"imageAlpineBionic": "imageHash-virtualization-imageAlpineBionic",
		"imageAlpineFocal":  "imageHash-virtualization-imageAlpineFocal",
		"imageAlpineJammy":  "imageHash-virtualization-imageAlpineJammy",
		"imageUbuntuBionic": "imageHash-virtualization-imageUbuntuBionic",
		"imageUbuntuFocal":  "imageHash-virtualization-imageUbuntuFocal",
		"imageUbuntuJammy":  "imageHash-virtualization-imageUbuntuJammy",
		"libguestfs":        "imageHash-virtualization-libguestfs",
		"virtApi":           "imageHash-virtualization-virtApi",
		"virtController":    "imageHash-virtualization-virtController",
		"virtExportproxy":   "imageHash-virtualization-virtExportproxy",
		"virtExportserver":  "imageHash-virtualization-virtExportserver",
		"virtHandler":       "imageHash-virtualization-virtHandler",
		"virtLauncher":      "imageHash-virtualization-virtLauncher",
		"virtOperator":      "imageHash-virtualization-virtOperator",
		"vmiRouter":         "imageHash-virtualization-vmiRouter",
	},
}
