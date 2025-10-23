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
