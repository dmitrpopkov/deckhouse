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

# Upload pki for system-registry


{{- if and .registry.registryMode (ne .registry.registryMode "Direct") }}

# Prepare vars
registry_pki_path="/etc/kubernetes/system-registry/pki"
etcd_pki_path="/etc/kubernetes/pki/etcd"

bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf get ns d8-system || bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf create ns d8-system

bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf -n d8-system delete secret registry-node-${D8_NODE_HOSTNAME}-pki || true
bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf -n d8-system create secret generic registry-node-${D8_NODE_HOSTNAME}-pki \
{{- if eq .registry.registryStorageMode "S3" }}
  --from-file=seaweedfs.key=$registry_pki_path/seaweedfs.key \
  --from-file=seaweedfs.crt=$registry_pki_path/seaweedfs.crt \
{{- end }}
  --from-file=auth.key=$registry_pki_path/auth.key \
  --from-file=auth.crt=$registry_pki_path/auth.crt \
  --from-file=distribution.key=$registry_pki_path/distribution.key \
  --from-file=distribution.crt=$registry_pki_path/distribution.crt \
  --labels=heritage=deckhouse,module=embedded-registry,type=node-secret

bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf -n d8-system delete secret registry-pki || true
bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf -n d8-system create secret generic registry-pki \
{{- if eq .registry.registryStorageMode "S3" }}
  --from-file=etcd-ca.key=$etcd_pki_path/ca.key \
  --from-file=etcd-ca.crt=$etcd_pki_path/ca.crt \
{{- end }}
  --from-file=registry-ca.key=$registry_pki_path/ca.key \
  --from-file=registry-ca.crt=$registry_pki_path/ca.crt \
  --labels=heritage=deckhouse,module=embedded-registry,type=ca-secret

bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf -n d8-system delete secret registry-user-rw || true
bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf -n d8-system create secret generic registry-user-rw \
  --from-literal=name='{{- .registry.internalRegistryAccess.userRw.name }}' \
  --from-literal=password='{{- .registry.internalRegistryAccess.userRw.password }}' \
  --from-literal=passwordHash='{{- .registry.internalRegistryAccess.userRw.passwordHash }}' \
  --type='system-registry/user' \
  --labels=heritage=deckhouse,module=embedded-registry

bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf -n d8-system delete secret registry-user-ro || true
bb-kubectl --kubeconfig=/etc/kubernetes/admin.conf -n d8-system create secret generic registry-user-ro \
  --from-literal=name='{{- .registry.internalRegistryAccess.userRo.name }}' \
  --from-literal=password='{{- .registry.internalRegistryAccess.userRo.password }}' \
  --from-literal=passwordHash='{{- .registry.internalRegistryAccess.userRo.passwordHash }}' \
  --type='system-registry/user' \
  --labels=heritage=deckhouse,module=embedded-registry

{{- end }}
