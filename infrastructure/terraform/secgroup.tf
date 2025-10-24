resource "vkcs_networking_secgroup" "example" {
  name        = "example-secgroup"
  description = "Security group for SSH and HTTP access"
}

resource "vkcs_networking_secgroup_rule" "ssh_in" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 22
  port_range_max    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.example.id
}

resource "vkcs_networking_secgroup_rule" "http_in" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 80
  port_range_max    = 80
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.example.id
}

resource "vkcs_networking_secgroup_rule" "https_in" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 443
  port_range_max    = 443
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.example.id
}

resource "vkcs_networking_secgroup_rule" "frontend" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 3000
  port_range_max    = 3000
  security_group_id = vkcs_networking_secgroup.example.id
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "vkcs_networking_secgroup_rule" "backend" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 8080
  port_range_max    = 8080
  security_group_id = vkcs_networking_secgroup.example.id
  remote_ip_prefix  = "0.0.0.0/0"
}
