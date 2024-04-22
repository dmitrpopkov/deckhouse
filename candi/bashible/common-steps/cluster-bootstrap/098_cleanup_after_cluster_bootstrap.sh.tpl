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

rm -rf /var/lib/bashible/kubeadm
bb-package-remove kubeadm

{{- if and .registryMode (ne .registryMode "Direct") }}
bb-package-remove seaweedfs dockerAuth dockerDistribution etcd
rm -rf $IGNITER_DIR
{{- end }}

rm -f /tmp/bootstrap.sh
rm -rf /tmp/candi-bundle*
rm -f /tmp/bundle*.tar
rm -f /tmp/kubeadm-config*

# needs for dhctl bootstrap restarting check
echo "OK" > /opt/deckhouse/first-control-plane-bashible-ran
