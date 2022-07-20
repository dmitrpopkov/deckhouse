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

{{- $kubernetesVersion := printf "%s%s" (.kubernetesVersion | toString) (index .k8s .kubernetesVersion "patch" | toString) | replace "." "" }}
{{- $kubernetesCniVersion := index .k8s .kubernetesVersion "cniVersion" | toString | replace "." "" }}
bb-rp-remove kubeadm
bb-rp-install "kubernetes-cni:{{ index .images.registrypackages (printf "kubernetesCniUbuntu%s" $kubernetesCniVersion) | toString }}" "kubectl:{{ index .images.registrypackages (printf "kubectlUbuntu%s" $kubernetesVersion) | toString }}"

old_kubelet_hash=""
echo ${BB_RP_INSTALLED_PACKAGES_STORE}
if [ -f "${BB_RP_INSTALLED_PACKAGES_STORE}/kubelet/tag" ]; then
  old_kubelet_hash=$(<"${BB_RP_INSTALLED_PACKAGES_STORE}/kubelet/tag")
fi

# weird trick to prevent kubelet start after install
# https://jpetazzo.github.io/2013/10/06/policy-rc-d-do-not-start-services-automatically/
echo exit 101 > /usr/sbin/policy-rc.d
chmod +x /usr/sbin/policy-rc.d
bb-rp-install "kubelet:{{ index .images.registrypackages (printf "kubeletUbuntu%s" $kubernetesVersion) | toString }}"
rm -f /usr/sbin/policy-rc.d
new_kubelet_hash=$(<"${BB_RP_INSTALLED_PACKAGES_STORE}/kubelet/tag")
if [ "${old_kubelet_hash} -ne "${new_kubelet_hash}" ]; then
  bb-flag-set kubelet-need-restart
fi

mkdir -p /etc/kubernetes/manifests
mkdir -p /etc/systemd/system/kubelet.service.d
mkdir -p /etc/kubernetes/pki
