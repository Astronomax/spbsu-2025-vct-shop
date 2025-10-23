data "vkcs_compute_flavor" "compute" {
  name = var.compute_flavor
}

data "vkcs_images_image" "compute" {
  visibility = "public"
  default    = true
  properties = {
    mcs_os_distro  = "ubuntu"
    mcs_os_version = "22.04"
  }
}

locals {
  user_data = base64encode(<<-EOF
#!/bin/bash
set -e
apt-get update
apt-get upgrade -y
apt-get install -y curl wget git
echo "Virtual machine initialized successfully" > /var/log/init.log
EOF
  )
}

resource "vkcs_compute_keypair" "default" {
  name       = "default"
  public_key = file("~/.ssh/id_rsa.pub")
}

resource "vkcs_compute_instance" "example" {
  name              = var.instance_name
  flavor_id         = data.vkcs_compute_flavor.compute.id
  key_pair          = vkcs_compute_keypair.default.name
  security_groups   = [vkcs_networking_secgroup.example.name]
  availability_zone = var.availability_zone_name
  user_data         = local.user_data

  block_device {
    uuid                  = data.vkcs_images_image.compute.id
    source_type           = "image"
    destination_type      = "volume"
    volume_type           = "ceph-ssd"
    volume_size           = 10
    boot_index            = 0
    delete_on_termination = true
  }

  network {
    uuid = vkcs_networking_network.example.id
  }

  depends_on = [
    vkcs_networking_network.example,
    vkcs_networking_subnet.example
  ]
}

resource "vkcs_networking_floatingip" "example" {
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_compute_floatingip_associate" "example" {
  floating_ip = vkcs_networking_floatingip.example.address
  instance_id = vkcs_compute_instance.example.id
}

output "instance_floating_ip" {
  description = "Floating IP address of the instance"
  value       = vkcs_networking_floatingip.example.address
}

output "instance_internal_ip" {
  description = "Internal IP address of the instance"
  value       = vkcs_compute_instance.example.network[0].fixed_ip_v4
}
