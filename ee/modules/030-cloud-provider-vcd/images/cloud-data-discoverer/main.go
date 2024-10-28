/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package main

import (
	"os"

	"github.com/alecthomas/kingpin"

	cloud_data "github.com/deckhouse/deckhouse/go_lib/cloud-data"
	"github.com/deckhouse/deckhouse/go_lib/cloud-data/app"
)

func main() {
	kpApp := kingpin.New("vcd cloud data discoverer", "A tool for discovery data from cloud provider")
	kpApp.HelpFlag.Short('h')

	app.InitFlags(kpApp)

	logger := app.InitLogger()

	kpApp.Action(func(context *kingpin.ParseContext) error {
		client := app.InitClient(logger)
		dynamicClient := app.InitDynamicClient(logger)
		discoverer := NewDiscoverer(logger)

		r := cloud_data.NewReconciler(discoverer, app.ListenAddress, app.DiscoveryPeriod, logger, client, dynamicClient)
		r.Start()

		return nil
	})

	_, err := kpApp.Parse(os.Args[1:])
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
