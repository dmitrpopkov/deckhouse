{{- /*
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
*/}}
#!/bin/bash
export LANG=C
{{- if .proxy }}
  {{- if .proxy.httpProxy }}
HTTP_PROXY={{ .proxy.httpProxy | quote }}
  {{- end }}
  {{- if .proxy.httpsProxy }}
HTTPS_PROXY={{ .proxy.httpsProxy | quote }}
  {{- end }}
  {{- if .proxy.noProxy }}
NO_PROXY={{ .proxy.noProxy | join "," | quote }}
  {{- end }}
export HTTP_PROXY=${HTTP_PROXY}
export http_proxy=${HTTP_PROXY}
export HTTPS_PROXY=${HTTPS_PROXY}
export https_proxy=${HTTPS_PROXY}
export NO_PROXY=${NO_PROXY}
export no_proxy=${NO_PROXY}
{{- end }}
apt update
export DEBIAN_FRONTEND=noninteractive
until apt install jq netcat-openbsd curl -y; do
  echo "Error installing packages"
  apt update
  sleep 10
done
mkdir -p /var/lib/bashible/
