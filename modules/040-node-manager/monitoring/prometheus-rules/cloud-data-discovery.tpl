- name: cloud-data-discoverer.general
  rules:
  - alert: D8CloudDataDiscovererCloudRequestError
    for: 1h
    expr: max by(job, namespace)(cloud_data_discovery_cloud_request_error == 1)
    labels:
      severity_level: "6"
      d8_module: node-manager
      d8_component: cloud-data-discoverer
    annotations:
      plk_protocol_version: "1"
      plk_markup_format: "markdown"
      summary: Cloud data discoverer cannot get data from cloud
      plk_create_group_if_not_exists__malfunctioning: "D8CloudDataDiscovererMalfunctioning,tier=cluster,prometheus=deckhouse,kubernetes=~kubernetes"
      plk_grouped_by__malfunctioning: "D8CloudDataDiscovererMalfunctioning,tier=cluster,prometheus=deckhouse,kubernetes=~kubernetes"
      description: |
        Cloud data discoverer cannot get data from cloud. See cloud data discoverer logs for more information:
        `kubectl -n {{`{{ $labels.namespace }}`}} logs deploy/cloud-data-discoverer`

  - alert: D8CloudDataDiscovererSaveError
    for: 1h
    expr: max by(job, namespace)(cloud_data_discovery_update_resource_error == 1)
    labels:
      severity_level: "6"
      d8_module: node-manager
      d8_component: cloud-data-discoverer
    annotations:
      plk_protocol_version: "1"
      plk_markup_format: "markdown"
      summary: Cloud data discoverer cannot save data to k8s resource
      plk_create_group_if_not_exists__malfunctioning: "D8CloudDataDiscovererMalfunctioning,tier=cluster,prometheus=deckhouse,kubernetes=~kubernetes"
      plk_grouped_by__malfunctioning: "D8CloudDataDiscovererMalfunctioning,tier=cluster,prometheus=deckhouse,kubernetes=~kubernetes"
      description: |
        Cloud data discoverer cannot save data to k8s resource. See cloud data discoverer logs for more information:
        `kubectl -n {{`{{ $labels.namespace }}`}} logs deploy/cloud-data-discoverer`

{{/* TODO remove condition with label namespace after release 1.61 when all yandex clusters applied cloud-migrator in discoverer  */}}
  - alert: ClusterHasOrphanedDisks
    for: 1h
    expr: max by(job, id, name)(cloud_data_discovery_orphaned_disk_info{namespace!="d8-cloud-provider-yandex"} == 1)
    labels:
      severity_level: "6"
      d8_module: node-manager
      d8_component: cloud-data-discoverer
    annotations:
      plk_protocol_version: "1"
      plk_markup_format: "markdown"
      summary: Cloud data discoverer finds disks in the cloud for which there is no PersistentVolume in the cluster
      plk_create_group_if_not_exists__main: "ClusterHasCloudDataDiscovererAlerts,tier=cluster,prometheus=deckhouse,kubernetes=~kubernetes"
      plk_grouped_by__main: "ClusterHasCloudDataDiscovererAlerts,tier=cluster,prometheus=deckhouse,kubernetes=~kubernetes"
      description: |
        Cloud data discoverer finds disks in the cloud for which there is no PersistentVolume in the cluster. You can manually delete these disks from your cloud:
          ID: {{`{{ $labels.id }}`}}, Name: {{`{{ $labels.name }}`}}
