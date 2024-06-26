/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"trivy-provider/validators"
	"trivy-provider/web"

	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/javadb"
	"github.com/go-chi/chi/v5"
	"github.com/go-logr/zapr"
	"github.com/google/go-containerregistry/pkg/name"
	"go.uber.org/zap"
)

func main() {
	var (
		keyFile        string
		certFile       string
		clientCAFile   string
		host           string = "0.0.0.0"
		port           int    = 8443
		timeoutSeconds int    = 10
	)
	flag.StringVar(&keyFile, "key-file", keyFile, "Path to file containing TLS certificate key.")
	flag.StringVar(&certFile, "cert-file", certFile, "Path to file containing TLS certificate.")
	flag.StringVar(&clientCAFile, "client-ca-file", clientCAFile, "Path to client CA certificate (gatekeeper CA).")
	flag.StringVar(&host, "host", host, "Host for for the server to listen on.")
	flag.IntVar(&port, "port", port, "Port for the server to listen on.")
	flag.IntVar(&timeoutSeconds, "timeout", timeoutSeconds, "Scanning timeout in seconds.")
	flag.Parse()

	zapLog, err := zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic(fmt.Sprintf("unable to initialize logger: %v", err))
	}
	logger := zapr.NewLogger(zapLog).WithName("trivy-provider")

	if err = initJavaDB(); err != nil {
		panic(fmt.Sprintf("Couldn't initialize JavaDB: %v", err))
	}
	logger.Info("JavaDB was successfully initialized")

	tlsConfig, err := newTLSConfig(clientCAFile)
	if err != nil {
		panic(err)
	}

	handler := chi.NewRouter()
	server := &http.Server{
		Addr:      net.JoinHostPort(host, strconv.Itoa(port)),
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	remoteURL := os.Getenv("TRIVY_REMOTE_URL")
	scanner := validators.NewRemoteValidator(remoteURL, logger)
	scanningTimeout := time.Duration(timeoutSeconds) * time.Second
	validator := web.NewHandler(scanner, scanningTimeout, logger)
	handler.HandleFunc("/validate", validator.HandleRequest())

	logger.Info("starting server...")
	if err = server.ListenAndServeTLS(certFile, keyFile); err != nil {
		panic(err)
	}
}

func initJavaDB() error {
	javaDbImage := os.Getenv("TRIVY_JAVA_DB_IMAGE")
	if len(javaDbImage) == 0 {
		javaDbImage = "ghcr.io/aquasecurity/trivy-java-db:1"
	}

	ref, err := name.ParseReference(javaDbImage)
	if err != nil {
		return err
	}

	javadb.Init("/home/javadb", ref, false, true, ftypes.RegistryOptions{Insecure: false})
	if err = javadb.Update(); err != nil {
		return err
	}

	return nil
}

func newTLSConfig(caCertFile string) (*tls.Config, error) {
	if caCertFile == "" {
		return nil, nil
	}

	certPool, err := readCACert(caCertFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS13,
	}, nil
}

func readCACert(caCertFile string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load Gatekeeper's CA certificate %s: %w", caCertFile, err)
	}

	clientCAs := x509.NewCertPool()
	clientCAs.AppendCertsFromPEM(caCert)
	return clientCAs, nil
}
