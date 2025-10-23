data "vkcs_networking_network" "extnet" {
  external = true
}

resource "vkcs_networking_network" "example" {
  name = "example-network"
}

resource "vkcs_networking_subnet" "example" {
  name            = "example-subnet"
  network_id      = vkcs_networking_network.example.id
  cidr            = "192.168.199.0/24"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

resource "vkcs_networking_router" "example" {
  name                = "example-router"
  admin_state_up      = true
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_router_interface" "example" {
  router_id = vkcs_networking_router.example.id
  subnet_id = vkcs_networking_subnet.example.id
}
