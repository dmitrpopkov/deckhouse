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

image="{{ printf "%s%s:%s" $.registry.address $.registry.path (index $.images.controlPlaneManager "kubernetesApiProxy") }}"
if docker version >/dev/null 2>/dev/null; then
  HOME=/ docker pull "$image"
# Use ctr because it respects http_proxy and https_proxy environment variables and allows installing in air-gapped environments using simple http proxy.
# If ctr is not installed for some reason, fallback to crictl. This is just a small safety precaution.
elif ctr --version >/dev/null 2>/dev/null; then
  ctr -n=k8s.io images pull "$image"
elif crictl version >/dev/null 2>/dev/null; then
  crictl pull "$image"
fi
