// Copyright 2024 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/deckhouse/deckhouse/go_lib/registry-packages-proxy/proxy"

	"registry-packages-proxy/internal/app"
	"registry-packages-proxy/internal/credentials"
)

func main() {
	kpApp := kingpin.New("registry packages proxy", "A proxy for registry packages")
	kpApp.HelpFlag.Short('h')

	app.InitFlags(kpApp)

	kpApp.Action(func(_ *kingpin.ParseContext) error {
		ctx := context.Background()

		ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		app.RegisterMetrics()

		logger := app.InitLogger()
		client := app.InitClient(logger)
		dynamicClient := app.InitDynamicClient(logger)

		listener, err := net.Listen("tcp", app.ListenAddress)
		if err != nil {
			return errors.Wrap(err, "failed to listen")
		}

		watcher := credentials.NewWatcher(client, dynamicClient, app.RegistrySecretDiscoveryPeriod, logger)

		go watcher.Watch(ctx)

		cache, err := app.NewCache(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to create cache")
		}
		defer cache.Close()

		if cache == nil {
			logger.Info("Cache is disabled")
		}

		proxy, err := proxy.NewProxy(listener, app.BuildRouter(), watcher, proxy.Options{
			Cache:  cache,
			Logger: logger,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create proxy")
		}

		logger.Infof("Listening on %s", listener.Addr().String())

		err = proxy.Serve()
		if err != nil {
			return errors.Wrap(err, "failed to serve")
		}

		return nil
	})

	_, err := kpApp.Parse(os.Args[1:])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
