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

package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	listenHost  = "0.0.0.0"
	listenPort  = "8080"
	madisonHost = "madison.flant.com"
)

func main() {
	madisonScheme := os.Getenv("MADISON_SCHEME")
	if madisonScheme == "" {
		log.Fatal("MADISON_SCHEME is not set")
	}
	madisonBackend := os.Getenv("MADISON_BACKEND")
	if madisonBackend == "" {
		log.Fatal("MADISON_BACKEND is not set")
	}
	madisonKey := os.Getenv("MADISON_AUTH_KEY")
	if madisonKey == "" {
		log.Fatal("MADISON_AUTH_KEY is not set")
	}

	proxy := newMadisonProxy(madisonScheme, madisonBackend, madisonKey)

	mux := http.NewServeMux()

	mux.HandleFunc("/readyz", readyHandler)
	mux.HandleFunc("/healthz", readyHandler)
	mux.HandleFunc("/", proxy.ServeHTTP)

	s := &http.Server{
		Addr:    listenHost + ":" + listenPort,
		Handler: mux,
	}
	go func() {
		err := s.ListenAndServe()
		if err == nil || err == http.ErrServerClosed {
			log.Info("Shutting down.")
			return
		}
		log.Error(err)
	}()

	// Block to wait for a signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigs

	// 30 sec is the readiness check timeout
	deadline := time.Now().Add(30 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	log.Info("Got signal ", sig)
	err := s.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func newMadisonProxy(madisonScheme, madisonBackend, madisonAuthKey string) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = madisonScheme
			req.URL.Host = madisonBackend
			req.URL.Path = "/api/events/prometheus/" + madisonAuthKey
			req.Host = madisonHost
			req.Header.Set("Host", madisonHost)
		},
	}
}
