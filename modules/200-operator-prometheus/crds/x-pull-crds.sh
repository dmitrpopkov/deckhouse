#!/usr/bin/env bash

#
#Copyright 2023 Flant JSC
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.
#

version="0.68.0"

repo="https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/v${version}/example/prometheus-operator-crd/monitoring.coreos.com"

crds=("alertmanagers.yaml alertmanagerconfigs.yaml podmonitors.yaml probes.yaml prometheuses.yaml prometheusrules.yaml servicemonitors.yaml thanosrulers.yaml scrapeconfigs.yaml prometheusagents.yaml")

for crd in $crds
do
  file="${repo}_${crd}"
  curl -s ${file} > ${crd}
done

