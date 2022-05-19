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

data "yandex_compute_image" "nat_image" {
  count = var.should_create_nat_instance ? 1 : 0
  family = "nat-instance-ubuntu"
}

data "yandex_vpc_subnet" "internal_subnet" {
  subnet_id = var.nat_instance_internal_subnet_id == null ? (local.should_create_subnets ? yandex_vpc_subnet.kube_c[0].id : data.yandex_vpc_subnet.kube_c[0].id) : var.nat_instance_internal_subnet_id
}

data "yandex_vpc_subnet" "external_subnet" {
  count = var.nat_instance_external_subnet_id == null ? 0 : 1
  subnet_id = var.nat_instance_external_subnet_id
}

locals {
  internal_subnet_zone = data.yandex_vpc_subnet.internal_subnet.zone
  external_subnet_zone = var.nat_instance_external_subnet_id == null ? null : join("", data.yandex_vpc_subnet.external_subnet.*.zone) # https://github.com/hashicorp/terraform/issues/23222#issuecomment-547462883

  zone_to_cidr = tomap({
    "ru-central1-a" = local.should_create_subnets ? local.kube_a_v4_cidr_block : (local.is_existing_subnet_a ? data.yandex_vpc_subnet.kube_a[0] : null)
    "ru-central1-b" = local.should_create_subnets ? local.kube_b_v4_cidr_block : (local.is_existing_subnet_b ? data.yandex_vpc_subnet.kube_b[0] : null)
    "ru-central1-c" = local.should_create_subnets ? local.kube_c_v4_cidr_block : (local.is_existing_subnet_c ? data.yandex_vpc_subnet.kube_c[0] : null)
  })

  with_internal_nat_instance_internal_cidr = var.nat_instance_internal_subnet_id == null ? null : data.yandex_vpc_subnet.kube_c[0].v4_cidr_blocks[0]
  with_external_nat_instance_internal_cidr = var.nat_instance_external_subnet_id == null ? null : local.zone_to_cidr[local.external_subnet_zone]
  manual_subnets_nat_instance_internal_cidr = local.should_create_subnets ? null : data.yandex_vpc_subnet.kube_c[0].v4_cidr_blocks[0]
  nat_instance_cidr = coalesce(local.with_internal_nat_instance_internal_cidr, local.with_external_nat_instance_internal_cidr, local.manual_subnets_nat_instance_internal_cidr, local.kube_c_v4_cidr_block)

  nat_instance_internal_address_calculated = var.should_create_nat_instance ? (var.nat_instance_internal_address == null ? cidrhost(local.nat_instance_cidr, 10) : var.nat_instance_internal_address) : null

  assign_external_ip_address = var.nat_instance_external_subnet_id == null ? true : false

  user_data = <<-EOT
    #!/bin/bash

    echo "${var.nat_instance_ssh_key}" >> /home/ubuntu/.ssh/authorized_keys

    cat > /etc/netplan/49-second-interface.yaml <<"EOF"
    network:
      version: 2
      ethernets:
        eth1:
          dhcp4: yes
          dhcp4-overrides:
            use-dns: false
            use-ntp: false
            use-hostname: false
            use-routes: false
            use-domains: false
    EOF

    netplan apply
  EOT
}

resource "yandex_compute_instance" "nat_instance" {
  count = var.should_create_nat_instance ? 1 : 0

  name         = join("-", [var.prefix, "nat"])
  hostname     = join("-", [var.prefix, "nat"])
  zone         = var.nat_instance_external_subnet_id == null ? local.internal_subnet_zone : local.external_subnet_zone

  platform_id  = "standard-v2"
  resources {
    cores  = 2
    memory = 2
  }

  boot_disk {
    initialize_params {
      type = "network-hdd"
      image_id = data.yandex_compute_image.nat_image.0.image_id
      size = 13
    }
  }

  dynamic "network_interface" {
    for_each = var.nat_instance_external_subnet_id != null ? [var.nat_instance_external_subnet_id] : []
    content {
      subnet_id = network_interface.value
      ip_address = var.nat_instance_external_address
      nat       = false
    }
  }

  network_interface {
    subnet_id      = data.yandex_vpc_subnet.internal_subnet.id
    nat            = local.assign_external_ip_address
    nat_ip_address = local.assign_external_ip_address ? var.nat_instance_external_address : null
    ip_address     = local.nat_instance_internal_address_calculated
  }

  lifecycle {
    ignore_changes = [
      hostname,
      metadata,
      boot_disk[0].initialize_params[0].image_id,
      boot_disk[0].initialize_params[0].size
    ]
  }

  metadata = {
    user-data = local.user_data
  }
}
