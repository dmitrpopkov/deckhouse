{{- define "bootstrap_script" -}}
#!/usr/bin/env bash

# FIXME
set -x

set -Eeuo pipefail

function detect_bundle() {
  {{- .Files.Get "candi/bashible/detect_bundle.sh" | nindent 2 }}
}

function get_bootstrap() {
  token="$(</var/lib/bashible/bootstrap-token)"
  bundle="$(detect_bundle)"
  node_group_name="{{ .nodeGroupName }}"

  bootstrap_bundle_name="$bundle.$node_group_name"

  while true; do
    for server in {{ .apiserverEndpoints | join " " }}; do
      url="https://$server/apis/bashible.deckhouse.io/v1alpha1/bootstrap/$bootstrap_bundle_name"
      if curl -s -f -x "" -X GET "$url" --header "Authorization: Bearer $token" --cacert "/var/lib/bashible/ca.crt"
      then
        return 0
      else
        >&2 echo "failed to get bootstrap $bootstrap_bundle_name with curl https://$server..."
      fi
    done
    sleep 10
  done
}

bootstrap_object="$(get_bootstrap)"

bash <<< "$(echo "$bootstrap_object" | jq -r .bootstrap)"
{{- end }}
