# Copyright 2021 Flant CJSC
# Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE

module "network_security" {
  source = "../../../terraform-modules/network-security"
  prefix = local.prefix
  remote_ip_prefix = "0.0.0.0/0"
  enabled = local.network_security
}

module "keypair" {
  source = "../../../terraform-modules/keypair"
  prefix = local.prefix
  ssh_public_key = var.providerClusterConfiguration.sshPublicKey
}

data "openstack_compute_availability_zones_v2" "zones" {}

data "openstack_networking_network_v2" "external" {
  name = local.external_network_name
}

resource "openstack_networking_network_v2" "internal" {
  name = local.prefix
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "internal" {
  name = local.prefix
  network_id = openstack_networking_network_v2.internal.id
  cidr = local.internal_network_cidr
  ip_version = 4
  gateway_ip = cidrhost(local.internal_network_cidr, 1)
  enable_dhcp = "true"
  allocation_pool {
    start = cidrhost(local.internal_network_cidr, 2)
    end = cidrhost(local.internal_network_cidr, 254)
  }
  dns_nameservers = var.providerClusterConfiguration.standard.internalNetworkDNSServers
}

resource "openstack_networking_router_v2" "router" {
  name = local.prefix
  admin_state_up = "true"
  external_network_id = data.openstack_networking_network_v2.external.id
}

resource "openstack_networking_router_interface_v2" "router" {
  router_id = openstack_networking_router_v2.router.id
  subnet_id = openstack_networking_subnet_v2.internal.id
}
