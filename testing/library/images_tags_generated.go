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
		"cloudDataDiscoverer":       "imageHash-cloudProviderAws-cloudDataDiscoverer",
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
		"cloudDataDiscoverer":       "imageHash-cloudProviderAzure-cloudDataDiscoverer",
	},
	"cloudProviderGcp": map[string]interface{}{
		"cloudControllerManager122": "imageHash-cloudProviderGcp-cloudControllerManager122",
		"cloudControllerManager123": "imageHash-cloudProviderGcp-cloudControllerManager123",
		"cloudControllerManager124": "imageHash-cloudProviderGcp-cloudControllerManager124",
		"cloudControllerManager125": "imageHash-cloudProviderGcp-cloudControllerManager125",
		"cloudControllerManager126": "imageHash-cloudProviderGcp-cloudControllerManager126",
		"cloudDataDiscoverer":       "imageHash-cloudProviderGcp-cloudDataDiscoverer",
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
		"cloudDataDiscoverer":       "imageHash-cloudProviderOpenstack-cloudDataDiscoverer",
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
		"builderArtifact":        "imageHash-cniCilium-builderArtifact",
		"builderRuntimeArtifact": "imageHash-cniCilium-builderRuntimeArtifact",
		"cilium":                 "imageHash-cniCilium-cilium",
		"operator":               "imageHash-cniCilium-operator",
		"virtCilium":             "imageHash-cniCilium-virtCilium",
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
	"delivery": map[string]interface{}{
		"argocd":               "imageHash-delivery-argocd",
		"argocdImageUpdater":   "imageHash-delivery-argocdImageUpdater",
		"redis":                "imageHash-delivery-redis",
		"werfArgocdCmpSidecar": "imageHash-delivery-werfArgocdCmpSidecar",
	},
	"descheduler": map[string]interface{}{
		"descheduler": "imageHash-descheduler-descheduler",
	},
	"documentation": map[string]interface{}{
		"web": "imageHash-documentation-web",
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
		"kruiseStateMetrics":    "imageHash-ingressNginx-kruiseStateMetrics",
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
		"alertsReceiver":              "imageHash-prometheus-alertsReceiver",
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
		"containerd1620":                      "imageHash-registrypackages-containerd1620",
		"containerdAlteros14631El7X8664":      "imageHash-registrypackages-containerdAlteros14631El7X8664",
		"containerdAlteros161831El7X8664":     "imageHash-registrypackages-containerdAlteros161831El7X8664",
		"containerdAstra1461Buster":           "imageHash-registrypackages-containerdAstra1461Buster",
		"containerdAstra16181Buster":          "imageHash-registrypackages-containerdAstra16181Buster",
		"containerdCentos14631El7X8664":       "imageHash-registrypackages-containerdCentos14631El7X8664",
		"containerdCentos14631El8X8664":       "imageHash-registrypackages-containerdCentos14631El8X8664",
		"containerdCentos16731El9X8664":       "imageHash-registrypackages-containerdCentos16731El9X8664",
		"containerdDebian1431Stretch":         "imageHash-registrypackages-containerdDebian1431Stretch",
		"containerdDebian1461Bullseye":        "imageHash-registrypackages-containerdDebian1461Bullseye",
		"containerdDebian1461Buster":          "imageHash-registrypackages-containerdDebian1461Buster",
		"containerdRedos14631El7X8664":        "imageHash-registrypackages-containerdRedos14631El7X8664",
		"containerdRedos161831El7X8664":       "imageHash-registrypackages-containerdRedos161831El7X8664",
		"containerdUbuntu1461Bionic":          "imageHash-registrypackages-containerdUbuntu1461Bionic",
		"containerdUbuntu1461Focal":           "imageHash-registrypackages-containerdUbuntu1461Focal",
		"containerdUbuntu15111Jammy":          "imageHash-registrypackages-containerdUbuntu15111Jammy",
		"crictl122":                           "imageHash-registrypackages-crictl122",
		"crictl123":                           "imageHash-registrypackages-crictl123",
		"crictl124":                           "imageHash-registrypackages-crictl124",
		"crictl125":                           "imageHash-registrypackages-crictl125",
		"crictl126":                           "imageHash-registrypackages-crictl126",
		"d8Curl801":                           "imageHash-registrypackages-d8Curl801",
		"dockerAlteros1903153El7X866473":      "imageHash-registrypackages-dockerAlteros1903153El7X866473",
		"dockerAltlinux2301Alt1X8664":         "imageHash-registrypackages-dockerAltlinux2301Alt1X8664",
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
		"kubeadm12217":                        "imageHash-registrypackages-kubeadm12217",
		"kubeadm12317":                        "imageHash-registrypackages-kubeadm12317",
		"kubeadm12413":                        "imageHash-registrypackages-kubeadm12413",
		"kubeadm1259":                         "imageHash-registrypackages-kubeadm1259",
		"kubeadm1264":                         "imageHash-registrypackages-kubeadm1264",
		"kubectl12217":                        "imageHash-registrypackages-kubectl12217",
		"kubectl12317":                        "imageHash-registrypackages-kubectl12317",
		"kubectl12413":                        "imageHash-registrypackages-kubectl12413",
		"kubectl1259":                         "imageHash-registrypackages-kubectl1259",
		"kubectl1264":                         "imageHash-registrypackages-kubectl1264",
		"kubelet12217":                        "imageHash-registrypackages-kubelet12217",
		"kubelet12317":                        "imageHash-registrypackages-kubelet12317",
		"kubelet12413":                        "imageHash-registrypackages-kubelet12413",
		"kubelet1259":                         "imageHash-registrypackages-kubelet1259",
		"kubelet1264":                         "imageHash-registrypackages-kubelet1264",
		"kubernetesCni120":                    "imageHash-registrypackages-kubernetesCni120",
		"tomlMerge01":                         "imageHash-registrypackages-tomlMerge01",
		"virtWhatDebian1151Deb9u1":            "imageHash-registrypackages-virtWhatDebian1151Deb9u1",
	},
	"runtimeAuditEngine": map[string]interface{}{
		"falco":         "imageHash-runtimeAuditEngine-falco",
		"falcosidekick": "imageHash-runtimeAuditEngine-falcosidekick",
		"rulesLoader":   "imageHash-runtimeAuditEngine-rulesLoader",
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
