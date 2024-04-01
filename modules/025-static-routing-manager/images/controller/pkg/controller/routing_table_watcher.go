/*
Copyright 2024 Flant JSC

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

package controller

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"static-routing-manager-controller/api/v1alpha1"
	"static-routing-manager-controller/pkg/config"
	"static-routing-manager-controller/pkg/logger"
	"static-routing-manager-controller/pkg/monitoring"
	"time"

	errors2 "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	RoutingTableCtrlName = "static-routing-manager-controller"

	RouteTableIDMin int = 0
	RouteTableIDMax int = 4294967295

	EmptyRouteTableIDReconcile reconcileType = "EmptyRouteTableID"
)

type (
	reconcileType string
)

func RunRoutingTableWatcherController(
	mgr manager.Manager,
	cfg config.Options,
	log logger.Logger,
	metrics monitoring.Metrics,
) (controller.Controller, error) {
	cl := mgr.GetClient()

	c, err := controller.New(RoutingTableCtrlName, mgr, controller.Options{
		Reconciler: reconcile.Func(func(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
			log.Info("[RoutingTableReconciler] starts Reconcile")
			rt := &v1alpha1.RoutingTable{}
			err := cl.Get(ctx, request.NamespacedName, rt)
			if err != nil && !errors2.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("[RoutingTableReconciler] unable to get RoutingTable, name: %s", request.Name))
				return reconcile.Result{}, err
			}

			if rt.Name == "" {
				log.Info(fmt.Sprintf("[RoutingTableReconciler] seems like the RoutingTable for the request %s was deleted. Reconcile retrying will stop.", request.Name))
				return reconcile.Result{}, nil
			}

			shouldRequeue, err := runEventReconcile(ctx, cl, log, rt)
			if err != nil {
				log.Error(err, fmt.Sprintf("[RoutingTableReconciler] an error occured while reconciles the RoutingTable, name: %s", rt.Name))
			}

			if shouldRequeue {
				log.Warning(fmt.Sprintf("[RoutingTableReconciler] Reconciler will requeue the request, name: %s", request.Name))
				return reconcile.Result{
					RequeueAfter: cfg.RequeueInterval * time.Second,
				}, nil
			}

			log.Info("[RoutingTableReconciler] ends Reconcile")
			return reconcile.Result{}, nil
		}),
	})
	if err != nil {
		log.Error(err, "[RunRoutingTableWatcherController] unable to create controller")
		return nil, err
	}

	err = c.Watch(source.Kind(mgr.GetCache(), &v1alpha1.RoutingTable{}), &handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "[RunRoutingTableWatcherController] unable to watch the events")
		return nil, err
	}

	return c, nil
}

func runEventReconcile(ctx context.Context, cl client.Client, log logger.Logger, rt *v1alpha1.RoutingTable) (bool, error) {
	recType, err := identifyReconcileFunc(rt)
	if err != nil {
		log.Error(err, fmt.Sprintf("[runEventReconcile] unable to identify reconcile func for the RoutingTable %s", rt.Name))
		return true, err
	}
	log.Debug(fmt.Sprintf("[runEventReconcile] reconcile operation: %s", recType))
	switch recType {
	case EmptyRouteTableIDReconcile:
		log.Debug(fmt.Sprintf("[runEventReconcile] EmptyRouteTableIDReconcile starts reconciliataion for the RoutingTable, name: %s", rt.Name))
		return reconcileRTGenerateIDFunc(ctx, cl, log, rt)
	default:
		log.Debug(fmt.Sprintf("[runEventReconcile] the RoutingTable %s should not be reconciled", rt.Name))
	}

	return false, nil
}

func identifyReconcileFunc(rt *v1alpha1.RoutingTable) (reconcileType, error) {
	should := shouldReconcileByEmptyRouteTableIDFunc(rt)
	if should {
		return EmptyRouteTableIDReconcile, nil
	}
	return "none", nil
}

func shouldReconcileByEmptyRouteTableIDFunc(rt *v1alpha1.RoutingTable) bool {
	if rt.DeletionTimestamp != nil {
		return false
	}
	reflect.ValueOf(rt.Spec.IpRouteTableID).IsNil()

	if &rt.Spec.IpRouteTableID != nil || &rt.Status.IpRouteTableID != nil {
		return false
	}

	return true
}

func reconcileRTGenerateIDFunc(
	ctx context.Context,
	cl client.Client,
	log logger.Logger,
	rt *v1alpha1.RoutingTable,
) (bool, error) {
	log.Debug(fmt.Sprintf("[reconcileRTGenerateIDFunc] starts the RoutingTable %s validation", rt.Name))

	var newRTId int
	var err error

	if &rt.Spec.IpRouteTableID != nil {
		newRTId = rt.Spec.IpRouteTableID
	} else {
		newRTId, err = generateFreeRoutingTableID(ctx, cl, log)
		if err != nil {
			log.Error(err, fmt.Sprintf("[reconcileRTGenerateIDFunc] unable to generate free RoutingTableID"))
			return true, err
		}
	}

	err = updateRoutingTableIDInStatus(ctx, cl, rt, newRTId)
	if err != nil {
		log.Error(err, fmt.Sprintf("[reconcileRTGenerateIDFunc] unable to update the RoutingTable, name: %s", rt.Name))
		return true, err
	}
	log.Debug(fmt.Sprintf("[reconcileRTGenerateIDFunc] successfully updated the RoutingTable %s status", rt.Name))

	return false, nil
}

func generateFreeRoutingTableID(
	ctx context.Context,
	cl client.Client,
	log logger.Logger,
) (int, error) {
	rtList := &v1alpha1.RoutingTableList{}
	err := cl.List(ctx, rtList)
	if err != nil {
		log.Error(err, "[generateFreeRoutingTableID] unable to list Routing Tables")
		return 65536, err
	}
LABEL:
	for {
		randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
		newRTId := randomizer.Intn(RouteTableIDMax-RouteTableIDMin) + RouteTableIDMin
		for _, rt := range rtList.Items {
			if rt.Status.IpRouteTableID == newRTId {
				continue LABEL
			}
		}
		return newRTId, nil
	}
}

func updateRoutingTableIDInStatus(
	ctx context.Context,
	cl client.Client,
	rt *v1alpha1.RoutingTable,
	newRTId int,
) error {
	rt.Status.IpRouteTableID = newRTId

	// TODO: add retry logic
	err := cl.Update(ctx, rt)
	if err != nil {
		return err
	}

	return nil
}
