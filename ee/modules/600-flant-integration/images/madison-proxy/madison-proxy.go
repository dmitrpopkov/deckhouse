/*
Copyright 2021 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package main

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const logFormat = `request:%s %s %s, status:%s, body:%s`

// This type implements the http.RoundTripper interface
type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	var reqBody []byte

	// Send the request, get the response (or the error)
	res, e = lrt.Proxied.RoundTrip(req)

	if req.Body != nil {
		var err error
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			log.Error(err)
		}
	}
	log.Infof(logFormat, req.Method, req.URL, req.Proto, res.Status, string(reqBody))
	return
}

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
	transport := http.DefaultTransport
	transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &httputil.ReverseProxy{
		Transport: LoggingRoundTripper{transport},
		Director: func(req *http.Request) {
			req.URL.Scheme = madisonScheme
			req.URL.Host = madisonBackend
			if req.URL.Path == "/api/v1/alerts" || req.URL.Path == "/api/v2/alerts" {
				req.URL.Path = "/api/events/prometheus/" + madisonAuthKey
			}
			if req.URL.Path == "/readyz" {
				req.URL.Path = "/healthz"
			}
			req.Host = madisonHost
			req.Header.Set("Host", madisonHost)
		},
	}
}
