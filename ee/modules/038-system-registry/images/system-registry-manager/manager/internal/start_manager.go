/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package internal

import (
	"context"
	"sync"
	common_config "system-registry-manager/internal/common"
	"system-registry-manager/internal/master"
	"system-registry-manager/internal/worker"
	pkg_logs "system-registry-manager/pkg/logs"
)

const (
	mainProcessName = "main"
)

func updateMainManageContext(ctx context.Context) context.Context {
	ctx = pkg_logs.SetLoggerToContext(ctx, mainProcessName)
	return ctx
}

func StartManager() {
	rootCtx, rootCtxcancel := context.WithCancel(context.Background())
	defer rootCtxcancel()

	rootCtx = updateMainManageContext(rootCtx)
	log := pkg_logs.GetLoggerFromContext(rootCtx)

	cfg := common_config.NewRuntimeConfig(&rootCtxcancel)
	worker := worker.New(rootCtx, cfg)
	master, err := master.New(rootCtx, cfg)
	if err != nil {
		log.Errorf("error, %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(3) // Changed the value to 2 since we have only two worker goroutines

	// Goroutine for handling signals
	go func() {
		defer wg.Done()
		log.Info("Starting os signal handler...")
		StartOsSignalHandler(rootCtx, cfg, worker.Stop)
	}()

	// Start worker goroutine
	go func() {
		defer wg.Done()
		log.Info("Starting worker...")
		worker.Start()
	}()

	// Start master goroutine
	go func() {
		defer wg.Done()
		log.Info("Starting master...")
		master.Start()
	}()

	wg.Wait()
	log.Info("Shutting down...")
}
