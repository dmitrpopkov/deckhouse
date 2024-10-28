[filer.options]
recursive_delete = false # do we really need for registry?

[etcd]
enabled = true
servers = "{{ .etcd.addresses | join "," }}"

key_prefix = "seaweedfs_meta."
tls_ca_file= "/kubernetes_pki/etcd/ca.crt"
tls_client_crt_file="/system_registry_pki/seaweedfs-etcd-client.crt"
tls_client_key_file="/system_registry_pki/seaweedfs-etcd-client.key"
