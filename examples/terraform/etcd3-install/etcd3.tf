// Create popular profiles (convenience module)
module "profiles" {
  source                  = "../modules/profiles"
  matchbox_http_endpoint  = "${var.matchbox_http_endpoint}"
  container_linux_version = "1520.8.0"
  container_linux_channel = "stable"
  install_disk            = "${var.install_disk}"
  container_linux_oem     = "${var.container_linux_oem}"
}

// Install Container Linux to disk before provisioning
resource "matchbox_group" "default" {
  name    = "default"
  profile = "${module.profiles.cached-container-linux-install}"

  // No selector, matches all nodes

  metadata {
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

// Create matcher groups for 3 machines

resource "matchbox_group" "node1" {
  name    = "node1"
  profile = "${module.profiles.etcd3}"

  selector {
    mac = "52:54:00:a1:9c:ae"
    os  = "installed"
  }

  metadata {
    domain_name          = "node1.example.com"
    etcd_name            = "node1"
    etcd_initial_cluster = "node1=http://node1.example.com:2380,node2=http://node2.example.com:2380,node3=http://node3.example.com:2380"
    ssh_authorized_key   = "${var.ssh_authorized_key}"
  }
}

resource "matchbox_group" "node2" {
  name    = "node2"
  profile = "${module.profiles.etcd3}"

  selector {
    mac = "52:54:00:b2:2f:86"
    os  = "installed"
  }

  metadata {
    domain_name          = "node2.example.com"
    etcd_name            = "node2"
    etcd_initial_cluster = "node1=http://node1.example.com:2380,node2=http://node2.example.com:2380,node3=http://node3.example.com:2380"
    ssh_authorized_key   = "${var.ssh_authorized_key}"
  }
}

resource "matchbox_group" "node3" {
  name    = "node3"
  profile = "${module.profiles.etcd3}"

  selector {
    mac = "52:54:00:c3:61:77"
    os  = "installed"
  }

  metadata {
    domain_name          = "node3.example.com"
    etcd_name            = "node3"
    etcd_initial_cluster = "node1=http://node1.example.com:2380,node2=http://node2.example.com:2380,node3=http://node3.example.com:2380"
    ssh_authorized_key   = "${var.ssh_authorized_key}"
  }
}
