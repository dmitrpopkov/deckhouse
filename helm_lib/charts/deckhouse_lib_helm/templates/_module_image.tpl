{{- /* Usage: {{ include "helm_lib_module_image" (list . "<container-name>") }} */ -}}
{{- /* returns image name */ -}}
{{- define "helm_lib_module_image" }}
  {{- $context := index . 0 }} {{- /* Template context with .Values, .Chart, etc */ -}}
  {{- $containerName := index . 1 | trimAll "\"" }} {{- /* Container name */ -}}
  {{- $moduleName := $context.Chart.Name | replace "-" "_" | camelcase | untitle }}
  {{- $imageHash := index $context.Values.global.modulesImages.tags $moduleName $containerName }}
  {{- if not $imageHash }}
  {{- $error := (printf "Image %s.%s has no tag" $moduleName $containerName ) }}
  {{- fail $error }}
  {{- end }}
  {{- $registryBase := $context.Values.global.modulesImages.registry.base }}
  {{- if $context.Values | dig  $moduleName "registry" "base" "" }}
    {{- $registryBase = (printf "%s/%s" (index $context.Values $moduleName "registry" "base") $moduleName) }}
  {{- end }}
  {{- printf "%s:%s" $registryBase $imageHash }}
{{- end }}

{{- /* Usage: {{ include "helm_lib_module_image_no_fail" (list . "<container-name>") }} */ -}}
{{- /* returns image name if found */ -}}
{{- define "helm_lib_module_image_no_fail" }}
  {{- $context := index . 0 }} {{- /* Template context with .Values, .Chart, etc */ -}}
  {{- $containerName := index . 1 | trimAll "\"" }} {{- /* Container name */ -}}
  {{- $moduleName := $context.Chart.Name | replace "-" "_" | camelcase | untitle }}
  {{- $imageHash := index $context.Values.global.modulesImages.tags $moduleName $containerName }}
  {{- if $imageHash }}
    {{- $registryBase := $context.Values.global.modulesImages.registry.base }}
    {{- if $context.Values | dig  $moduleName "registry" "base" "" }}
      {{- $registryBase = (printf "%s/%s" (index $context.Values $moduleName "registry" "base") $moduleName) }}
    {{- end }}
    {{- printf "%s:%s" $registryBase $imageHash }}
  {{- end }}
{{- end }}

{{- /* Usage: {{ include "helm_lib_module_common_image" (list . "<container-name>") }} */ -}}
{{- /* returns image name from common module */ -}}
{{- define "helm_lib_module_common_image" }}
  {{- $context := index . 0 }} {{- /* Template context with .Values, .Chart, etc */ -}}
  {{- $containerName := index . 1 | trimAll "\"" }} {{- /* Container name */ -}}
  {{- $imageHash := index $context.Values.global.modulesImages.tags "common" $containerName }}
  {{- if not $imageHash }}
  {{- $error := (printf "Image %s.%s has no tag" "common" $containerName ) }}
  {{- fail $error }}
  {{- end }}
  {{- printf "%s:%s" $context.Values.global.modulesImages.registry.base $imageHash }}
{{- end }}

{{- /* Usage: {{ include "helm_lib_module_common_image_no_fail" (list . "<container-name>") }} */ -}}
{{- /* returns image name from common module if found */ -}}
{{- define "helm_lib_module_common_image_no_fail" }}
  {{- $context := index . 0 }} {{- /* Template context with .Values, .Chart, etc */ -}}
  {{- $containerName := index . 1 | trimAll "\"" }} {{- /* Container name */ -}}
  {{- $imageHash := index $context.Values.global.modulesImages.tags "common" $containerName }}
  {{- if $imageHash }}
  {{- printf "%s:%s" $context.Values.global.modulesImages.registry.base $imageHash }}
  {{- end }}
{{- end }}
