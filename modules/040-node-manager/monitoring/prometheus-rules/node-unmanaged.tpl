- name: d8.node-unmanaged
  rules:
    - alert: D8NodeIsUnmanaged
      expr: max by (node) (d8_unmanaged_nodes_on_cluster) > 0
      for: 10m
      labels:
        tier: cluster
        severity_level: "9"
      annotations:
        plk_markup_format: "markdown"
        plk_protocol_version: "1"
        plk_create_group_if_not_exists__d8_cluster_has_unmanaged_nodes: "D8ClusterHasUnmanagedNodes,tier=cluster,prometheus=deckhouse,kubernetes=~kubernetes"
        plk_grouped_by__d8_cluster_has_unmanaged_nodes: "D8ClusterHasUnmanagedNodes,tier=cluster,prometheus=deckhouse,kubernetes=~kubernetes"
    {{- if .Values.global.modules.publicDomainTemplate }}
        summary: The {{`{{ $labels.node }}`}} Node is not managed by the [node-manager](https://deckhouse.io/documentation/v1/modules/040-node-manager/) module.
        description: |
          The {{`{{ $labels.node }}`}} Node is not managed by the [node-manager](https://deckhouse.io/documentation/v1/modules/040-node-manager/faq.html#how-to-put-an-existing-cluster-node-under-the-node-managers-control) module.
    {{- else }}
        summary: The {{`{{ $labels.node }}`}} Node is not managed by the `node-manager`.
        description: |
          The {{`{{ $labels.node }}`}} Node is not managed by the `node-manager`.
    {{- end }}

          The recommended actions are as follows:
          - Follow these instructions to clean up the node and add it to the cluster: https://deckhouse.io/documentation/v1/modules/040-node-manager/faq.html#how-to-clean-up-a-node-for-adding-to-the-cluster
