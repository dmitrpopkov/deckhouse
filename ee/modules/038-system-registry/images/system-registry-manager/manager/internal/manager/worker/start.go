/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package worker

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	common "system-registry-manager/internal/manager/common"
	pkg_logs "system-registry-manager/pkg/logs"
	"time"
)

const (
	processName     = "worker"
	serverAddr      = "127.0.0.1:8097"
	shutdownTimeout = 5 * time.Second
)

type Worker struct {
	commonCfg *common.RuntimeConfig
	rootCtx   context.Context
	log       *log.Entry

	server *http.Server
}

func New(rootCtx context.Context, rCfg *common.RuntimeConfig) *Worker {
	rootCtx = pkg_logs.SetLoggerToContext(rootCtx, processName)
	log := pkg_logs.GetLoggerFromContext(rootCtx)

	return &Worker{
		commonCfg: rCfg,
		rootCtx:   rootCtx,
		log:       log,
		server:    createServer(),
	}
}

func (m *Worker) Start() {
	m.log.Info("Worker starting...")
	if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		defer m.commonCfg.StopManager()
		m.log.Errorf("error starting server: %v", err)
	}
}

func (m *Worker) Stop() {
	m.log.Info("Worker shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := m.server.Shutdown(ctx); err != nil {
		m.log.Errorf("error shutting down server: %v", err)
	}
	m.log.Info("Worker shutdown")
}

func createServer() *http.Server {
	server := &http.Server{
		Addr: serverAddr,
	}
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/readyz", readyzHandler)
	return server
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
