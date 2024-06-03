/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package master

import (
	"context"
	"fmt"
	"system-registry-manager/internal/master/info"
	pkg_api "system-registry-manager/pkg/api"
	pkg_logs "system-registry-manager/pkg/logs"
	"time"
)

const (
	workInterval = 10 * time.Second
)

func startMasterWorkflow(ctx context.Context, m *Master) {
	log := pkg_logs.GetLoggerFromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(workInterval)
			err := masterWorkflow(ctx, m)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

func masterWorkflow(ctx context.Context, m *Master) error {
	log := pkg_logs.GetLoggerFromContext(ctx)
	info := info.NewInfo(log)
	workers, err := info.WorkersInfoWaitAll()
	if err != nil {
		return fmt.Errorf("error getting workers information: %v", err)
	}
	for _, worker := range workers {
		resp, err := worker.Client.RequestCheckRegistry(&pkg_api.CheckRegistryRequest{})
		if err != nil {
			return fmt.Errorf("error checking registry with worker %s: %v", worker.PodName, err)
		}
		if !(resp.Data.RegistryFilesState.ManifestsWaitToCreate ||
			resp.Data.RegistryFilesState.ManifestsWaitToUpdate ||
			resp.Data.RegistryFilesState.StaticPodsWaitToCreate ||
			resp.Data.RegistryFilesState.StaticPodsWaitToUpdate ||
			resp.Data.RegistryFilesState.CertificatesWaitToCreate ||
			resp.Data.RegistryFilesState.CertificatesWaitToUpdate) {
			continue
		}
		err = worker.Client.RequestUpdateRegistry(&pkg_api.CheckRegistryRequest{})
		if err != nil {
			return fmt.Errorf("error updating registry with worker %s: %v", worker.PodName, err)
		}
	}

	allInfo, err := info.AllInfoGet()
	if err != nil {
		return err
	}
	for node_name, node_info := range allInfo {
		log.Infof("Node: %s", node_name)
		log.Info("Node info:")
		log.Info(node_info)

		if node_info.SeaweedfsPod != nil {
			log.Info("Seaweedfs not exist")
			continue
		}

		log.Info("Seaweedfs exist")
		client, err := node_info.SeaweedfsPod.CreateClient()
		defer client.ClientClose()
		if err != nil {
			log.Fatal(err)
			continue
		}
		responseClusterCheck, err := client.ClusterCheck()
		if err != nil {
			log.Fatal(err)
			continue
		}
		log.Info(responseClusterCheck)
		responseClusterRaftPs, err := client.ClusterRaftPs()
		if err != nil {
			log.Fatal(err)
			continue
		}
		log.Info(responseClusterRaftPs)
	}
	return nil
}
