/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package main

import (
	"context"
	"fmt"
	"os"
	goruntime "runtime"

	"github.com/deckhouse/deckhouse/ee/modules/025-static-routing-manager/images/agent/api/v1alpha1"
	"github.com/deckhouse/deckhouse/ee/modules/025-static-routing-manager/images/agent/pkg/config"
	"github.com/deckhouse/deckhouse/ee/modules/025-static-routing-manager/images/agent/pkg/controller"
	"github.com/deckhouse/deckhouse/ee/modules/025-static-routing-manager/images/agent/pkg/kubutils"
	"github.com/deckhouse/deckhouse/ee/modules/025-static-routing-manager/images/agent/pkg/logger"

	v1 "k8s.io/api/core/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	resourcesSchemeFuncs = []func(*apiruntime.Scheme) error{
		v1alpha1.AddToScheme,
		v1alpha1.AddInternalToScheme,
		clientgoscheme.AddToScheme,
		extv1.AddToScheme,
		v1.AddToScheme,
	}
)

func main() {
	ctx := context.Background()
	cfgParams, err := config.NewConfig()
	if err != nil {
		fmt.Println("unable to create NewConfig " + err.Error())
		os.Exit(1)
	}

	log, err := logger.NewLogger(cfgParams.Loglevel)
	if err != nil {
		fmt.Println(fmt.Sprintf("unable to create NewLogger, err: %v", err))
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("[main] Go Version:%s ", goruntime.Version()))
	log.Info(fmt.Sprintf("[main] OS/Arch:Go OS/Arch:%s/%s ", goruntime.GOOS, goruntime.GOARCH))

	log.Info("[main] CfgParams has been successfully created")
	log.Info(fmt.Sprintf("[main] %s = %s", config.LogLevelENV, cfgParams.Loglevel))
	log.Info(fmt.Sprintf("[main] %s = %d", config.RequeueIntervalENV, cfgParams.RequeueInterval))
	log.Info(fmt.Sprintf("[main] %s = %d", config.PeriodicReconciliationIntervalENV, cfgParams.PeriodicReconciliationInterval))
	log.Info(fmt.Sprintf("[main] %s = %s", config.ProbeAddressPortENV, cfgParams.ProbeAddressPort))
	log.Info(fmt.Sprintf("[main] %s = %s", config.MetricsAddressPortENV, cfgParams.MetricsAddressPort))
	log.Info(fmt.Sprintf("[main] %s = %s", config.NodeNameENV, cfgParams.NodeName))

	kConfig, err := kubutils.KubernetesDefaultConfigCreate()
	if err != nil {
		log.Error(err, "[main] unable to KubernetesDefaultConfigCreate")
	}
	log.Info("[main] kubernetes config has been successfully created.")

	scheme := runtime.NewScheme()
	for _, f := range resourcesSchemeFuncs {
		err := f(scheme)
		if err != nil {
			log.Error(err, "[main] unable to add scheme to func")
			os.Exit(1)
		}
	}
	log.Info("[main] successfully read scheme CR")

	managerOpts := manager.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: cfgParams.ProbeAddressPort,
		Logger:                 log.GetLogger(),
		Metrics: metricsserver.Options{
			BindAddress: cfgParams.MetricsAddressPort,
		},
	}

	mgr, err := manager.New(kConfig, managerOpts)
	if err != nil {
		log.Error(err, "[main] unable to manager.New")
		os.Exit(1)
	}
	log.Info("[main] successfully created kubernetes manager")

	// metrics := monitoring.GetMetrics("")

	if _, err = controller.RunRoutesReconcilerAgentController(mgr, *cfgParams, *log); err != nil {
		log.Error(err, "[main] unable to controller.RunRoutesReconcilerAgentController")
		os.Exit(1)
	}

	if _, err = controller.RunIPRulesReconcilerAgentController(mgr, *cfgParams, *log); err != nil {
		log.Error(err, "[main] unable to controller.RunIPRulesReconcilerAgentController")
		os.Exit(1)
	}

	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Error(err, "[main] unable to mgr.AddHealthzCheck")
		os.Exit(1)
	}
	log.Info("[main] successfully AddHealthzCheck")

	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.Error(err, "[main] unable to mgr.AddReadyzCheck")
		os.Exit(1)
	}
	log.Info("[main] successfully AddReadyzCheck")

	err = mgr.Start(ctx)
	if err != nil {
		log.Error(err, "[main] unable to mgr.Start")
		os.Exit(1)
	}

	log.Info("[main] successfully starts the manager")
}
