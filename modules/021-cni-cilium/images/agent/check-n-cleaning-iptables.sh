#!/bin/bash

# Copyright 2024 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -Eeuo pipefail

function are_there_cilium_rules_in_legacy_iptables() {
  if iptables-legacy-save | grep -E "cilium|CILIUM" 2>&1 >/dev/null; then
    echo "### There are cilium rules in iptables-legacy"
    return 0
  fi
  return 1
}

function delete_cilium_legacy_iptables_rules_and_chains() {
  echo "### Start removing cilium rules and chains from iptables-legacy"
  for table in $(iptables-legacy-save | grep -E "^\*" | sed s/\*//g); do
    iptables-legacy --table $table --list-rules | grep -E "^-A.*(cilium|CILIUM)" | sed "s/-A/iptables-legacy --table $table -D/pe"
    iptables-legacy --table $table --list-rules | grep -E "^-N.*(cilium|CILIUM)" | sed "s/-N/iptables-legacy --table $table -X/pe"
  done
  echo "### Cilium rules and chains have been removed from iptables-legacy"
}

function flush_common_legacy_iptables_rules_and_chains() {
  for table in $(iptables-legacy-save | grep -E "^\*" | sed s/\*//g); do
    iptables-legacy --table $table -F
    iptables-legacy --table $table -X
    ip6tables-legacy --table $table -F
    ip6tables-legacy --table $table -X
  done
  echo "### Common chains have been flushed in iptables-legacy"
}

function delete_legacy_iptables() {
  for x in _raw _mangle _security _nat _filter; do
    modprobe -r "iptable${x}"
    modprobe -r "ip6table${x}"
  done
  echo "### iptables-legacy have been deleted"
}

function is_current_iptables_mode_eq_nft() {
  if [[ -v CURRENT_IPTABLES_MODE ]] && [[ "${CURRENT_IPTABLES_MODE}" = "nft" ]] ; then
    return 0
  fi
  return 1
}

if is_current_iptables_mode_eq_nft && are_there_cilium_rules_in_legacy_iptables; then
  delete_cilium_legacy_iptables_rules_and_chains
  flush_common_legacy_iptables_rules_and_chains
  delete_legacy_iptables
fi

echo "### The script has completed successfully"