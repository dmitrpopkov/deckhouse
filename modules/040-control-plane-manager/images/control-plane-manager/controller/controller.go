/*
Copyright 2023 Flant JSC

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

package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})

	var err error
	config, err = NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(config.ExitChannel)
	}()

	server := &http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/readyz", readyzHandler)

	go func() {
		err := server.ListenAndServe()
		if err == nil || err == http.ErrServerClosed {
			return
		}
		log.Fatal(err)
	}()

	if err := removeOrphanFiles(); err != nil {
		log.Warn(err)
	}

	if err := annotateNode(); err != nil {
		log.Fatal(err)
	}

	if err := waitNodeApproval(); err != nil {
		log.Fatal(err)
	}

	if err := waitImageHolderContainers(); err != nil {
		log.Fatal(err)
	}

	if err := checkEtcdManifest(); err != nil {
		log.Fatal(err)
	}

	if err := checkKubeletConfig(); err != nil {
		log.Fatal(err)
	}

	if err := installKubeadmConfig(); err != nil {
		log.Fatal(err)
	}

	if err := installBasePKIfiles(); err != nil {
		log.Fatal(err)
	}

	if err := fillTmpDirWithPKIData(); err != nil {
		log.Fatal(err)
	}

	if err := renewCertificates(); err != nil {
		log.Fatal(err)
	}

	if err := renewKubeconfigs(); err != nil {
		log.Fatal(err)
	}

	if err := updateRootKubeconfig(); err != nil {
		log.Fatal(err)
	}

	if err := installExtraFiles(); err != nil {
		log.Fatal(err)
	}

	if err := convergeComponents(); err != nil {
		log.Fatal(err)
	}

	if err := config.writeLastAppliedConfigurationChecksum(); err != nil {
		log.Fatal(err)
	}

	if err := cleanup(); err != nil {
		log.Warn(err)
	}

	controlPlaneManagerIsReady = true
	// pause loop
	<-config.ExitChannel

	if err := server.Close(); err != nil {
		log.Fatalf("HTTP close error: %v", err)
	}
}
