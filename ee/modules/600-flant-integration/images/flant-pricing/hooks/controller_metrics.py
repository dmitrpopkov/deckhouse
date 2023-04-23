#!/usr/bin/env python3
#
# Copyright 2023 Flant JSC
# Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
#
# This hook is responsible for generating metrics for d8 controllers resource consumption.


from utils import prometheus_query_value, prometheus_metric_builder, prometheus_function_builder
from typing import List, Dict, Any, TypeVar
from dataclasses import dataclass
from shell_operator import hook
from abc import ABC, abstractmethod


@dataclass
class Controller:
    name: str
    namespace: str
    kind: str
    module: str


MetricCollectorType = TypeVar("MetricCollectorType", bound="AbstractMetricCollector")


class AbstractMetricCollector(ABC):

    metric_group = "group_d8_controller_metrics"
    cpu_metric_name = "flant_pricing_controller_average_cpu_usage_seconds"
    memory_metric_name = "flant_pricing_controller_average_memory_working_set_bytes:without_kmem"
    interval = "5m"

    def collect(self, ctx: hook.Context, controllers: List[Controller]):
        '''Export metrics to hook context from Controllers list'''

        ctx.metrics.expire(self.metric_group)
        for ctrl in controllers:
            # Export metrics
            ctx.metrics.collect({
                "name": self.cpu_metric_name,
                "group": self.metric_group,
                "set": self.get_cpu_prometheus(ctrl),
                "labels": self.controller_metric_labels(ctrl),
            })

            ctx.metrics.collect({
                "name": self.memory_metric_name,
                "group": self.metric_group,
                "set": self.get_memory_prometheus(ctrl),
                "labels": self.controller_metric_labels(ctrl),
            })

    @abstractmethod
    def get_cpu_prometheus(self, controller: Controller) -> float:
        raise NotImplementedError('define get_cpu_prometheus to use this base class')

    @abstractmethod
    def get_memory_prometheus(self, controller: Controller) -> float:
        raise NotImplementedError('define get_memory_prometheus to use this base class')

    @staticmethod
    def controller_metric_labels(ctrl: Controller) -> Dict[str, str]:
        '''Helper func with labels for resource metrics'''

        return {
            "name": ctrl.name,
            "module": ctrl.module,
            "kind": ctrl.kind,
        }


class MetricCollector(AbstractMetricCollector):
    resource_consumption_query = '''
    sum (
        ( {} )
        + on(pod) group_left(controller_name, controller_type)
        ( {} * 0 )
    )
    '''

    def get_cpu_prometheus(self, controller: Controller) -> float:
        '''Query prometheus for controller cpu consumption'''

        return prometheus_query_value(
            query=self.resource_consumption_query.format(
                prometheus_function_builder(
                    f="rate",
                    metric=self.__resource_metric("container_cpu_usage_seconds_total", controller),
                    interval=self.interval,
                ),
                self.__controller_metric(controller),
            ),
        )

    def get_memory_prometheus(self, controller: Controller) -> float:
        '''Query prometheus for controller memory consumption'''

        return prometheus_query_value(
            query=self.resource_consumption_query.format(
                prometheus_function_builder(
                    f="avg_over_time",
                    metric=self.__resource_metric("container_memory_working_set_bytes:without_kmem", controller),
                    interval=self.interval,
                ),
                self.__controller_metric(controller),
            ),
        )

    def __controller_metric(controller: Controller) -> str:
        '''
        Generate kube_controller_pod metric from Controller instance
        input: Controller(name="dex", namespace="d8-user-authn, kind="Deployment")
        output: kube_controller_pod{namespace="d8-user-authn",controller_name="dex",controller_type="Deployment"}
        '''

        return prometheus_metric_builder(
            metric_name="kube_controller_pod",
            labels={
                "namespace": controller.namespace,
                "controller_name": controller.name,
                "controller_type": controller.kind
            },
        )

    def __resource_metric(metric_name: str, controller: Controller) -> str:
        '''
        Generate resource metric from Controller instance and metric name
        input: "container_cpu_usage_seconds_total", Controller(name="dex", namespace="d8-user-authn, kind="Deployment")
        output: `container_cpu_usage_seconds_total{namespace="d8-user-authn"}`
        '''

        return prometheus_metric_builder(
            metric_name=metric_name,
            labels={
                "namespace": controller.namespace,
            },
        )


class HookRunner:
    def __init__(self, metric_collector: MetricCollectorType):
        self.metric_collector = metric_collector

    def run(self, ctx: hook.Context):
        # Generate list of Controllers from snapshots
        controllers = self.__process_controllers(ctx.snapshots)

        # Generate metrics from Controllers list
        self.metric_collector.collect(ctx, controllers)

    def __process_controllers(self, snapshots: Dict[str, List[Dict[str, Any]]]) -> List[Controller]:
        '''Generate list of Controllers from binding context snapshots'''

        controllers = []
        for queue_snapshot in snapshots.values():
            for sn in queue_snapshot:
                controllers.append(self.__parse_controller(sn))
        return controllers

    def __parse_controller(self, controller_snapshot: Dict[str, Any]) -> Controller:
        '''
        Generate controller instance from snapshot
        '''

        filter_result = controller_snapshot["filterResult"]
        return Controller(
            kind=filter_result["kind"],
            name=filter_result["name"],
            namespace=filter_result["namespace"],
            module=filter_result["module"],
        )


if __name__ == "__main__":
    metric_collector = MetricCollector()
    hook_runner = HookRunner(metric_collector)
    hook.run(hook_runner.run, configpath="controller_metrics.yaml")
