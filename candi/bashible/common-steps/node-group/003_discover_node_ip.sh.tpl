# Copyright 2021 Flant JSC
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

# Ensure we have file
touch /var/lib/bashible/discovered-node-ip

{{- if ne .nodeGroup.nodeType "Static" }}
  {{ if eq .runType "ClusterBootstrap" }}
    {{- if and .clusterBootstrap.cloud .clusterBootstrap.cloud.nodeIP }}
echo {{ .clusterBootstrap.cloud.nodeIP }} > /var/lib/bashible/discovered-node-ip
    {{- end }}
  # For CloudEphemeral, CloudPermanent or CloudStatic node we try to discover IP from Node object
  {{- else }}
if [ -f /etc/kubernetes/kubelet.conf ] ; then
  if node="$(bb-kubectl --kubeconfig=/etc/kubernetes/kubelet.conf get node $HOSTNAME -o json 2> /dev/null)" ; then
    echo "$node" | jq -r '([.status.addresses[] | select(.type == "InternalIP") | .address] + [.status.addresses[] | select(.type == "ExternalIP") | .address])[0] // ""' > /var/lib/bashible/discovered-node-ip
  else
    bb-log-error "Unable to discover node IP for node object: No access to API server"
    exit 1
  fi
fi
  {{- end }}
{{- end }}

{{- if eq .nodeGroup.nodeType "Static" }}
# if node has one interface, use it's network as internal_network_cidr
internal_network_cidrs=""
if [[ "$(ip route show scope link | wc -l)" -eq 1 ]]; then
  internal_network_cidrs=$(ip route show scope link | awk '{print $1}')"
fi

  {{- if and (hasKey .nodeGroup "static") (hasKey .nodeGroup.static "internalNetworkCIDRs")}}
internal_network_cidrs="{{ .nodeGroup.static.internalNetworkCIDRs }}"
  {{- end }}

if [[ -z "$internal_network_cidrs" ]]; then
  bb-log-error "Cannot discover internal network CIDRs. Node has more than one interface, and nodeGroup.static.internalNetworkCIDRs is not set."
  bb-log-error "Please set nodeGroup.static.internalNetworkCIDRs to one of the node networks:"
  for network in $(ip route show scope link | awk '{print $1}'); do
    bb-log-error "  - $network"
  done
  exit 1
fi

function is_ip_in_cidr() {
  ip="$1"
  IFS="/" read net_address net_prefix <<< "$2"

  IFS=. read -r a b c d <<< "$ip"
  ip_dec="$((a * 256 ** 3 + b * 256 ** 2 + c * 256 + d))"

  IFS=. read -r a b c d <<< "$net_address"
  net_address_dec="$((a * 256 ** 3 + b * 256 ** 2 + c * 256 + d))"

  netmask=$(((0xFFFFFFFF << (32 - net_prefix)) & 0xFFFFFFFF))

  test $((netmask & ip_dec)) -eq $((netmask & net_address_dec))
}

if bb-is-ubuntu-version? 22.04 || bb-is-ubuntu-version? 20.04 || bb-is-ubuntu-version? 18.04; then
  ip_in_system=$(ip -f inet -br -j addr | jq -r '.[] | .addr_info[] | .local')
else
  ip_in_system=$(ip -f inet -br addr | grep -E '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}' -o)
fi

for cidr in $internal_network_cidrs; do
  for ip in $ip_in_system; do
    if is_ip_in_cidr "$ip" "$cidr"; then
      echo $ip > /var/lib/bashible/discovered-node-ip
      exit 0
    fi
  done
done
{{- end }}

{{- if or (eq .runType "ClusterBootstrap") (eq .nodeGroup.nodeType "Static") }}
if [ -z "$(cat /var/lib/bashible/discovered-node-ip)" ] ; then
  bb-log-error "Failed to discover node_ip but it's required for cluster bootstrap or static cluster nodes"
  exit 1
fi
{{- end }}
