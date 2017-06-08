// Default matcher group for machines
resource "matchbox_group" "default" {
  name    = "default"
  profile = "${matchbox_profile.coreos-install.name}"

  # no selector means all machines can be matched
  metadata {
    ignition_endpoint  = "${var.matchbox_http_endpoint}/ignition"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

// Match machines which have CoreOS Container Linux installed
resource "matchbox_group" "node1" {
  name    = "node1"
  profile = "${matchbox_profile.simple.name}"

  selector {
    os = "installed"
  }

  metadata {
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
