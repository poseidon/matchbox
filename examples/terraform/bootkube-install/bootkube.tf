// Create popular profiles (convenience module)
module "profiles" {
  source = "../modules/profiles"
  matchbox_http_endpoint = "${var.matchbox_http_endpoint}"
  container_linux_version = "1298.7.0"
  container_linux_channel = "stable"
}

// Install Container Linux to disk before provisioning
resource "matchbox_group" "default" {
  name = "default"
  profile = "${module.profiles.cached-container-linux-install}"
  // No selector, matches all nodes
  metadata {
    container_linux_channel = "stable"
    container_linux_version = "1298.7.0"
    ignition_endpoint = "${var.matchbox_http_endpoint}/ignition"
    baseurl = "${var.matchbox_http_endpoint}/assets/coreos"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

// Create a controller matcher group
resource "matchbox_group" "node1" {
  name = "node1"
  profile = "${module.profiles.bootkube-controller}"
  selector {
    mac = "52:54:00:a1:9c:ae"
    os = "installed"
  }
  metadata {
    domain_name = "node1.example.com"
    etcd_name = "node1"
    etcd_initial_cluster = "node1=http://node1.example.com:2380"
    etcd_on_host = "${var.experimental_self_hosted_etcd ? "false" : "true"}"
    k8s_dns_service_ip = "${var.k8s_dns_service_ip}"
    k8s_etcd_service_ip = "${var.k8s_etcd_service_ip}"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

// Create worker matcher groups

resource "matchbox_group" "node2" {
  name = "node2"
  profile = "${module.profiles.bootkube-worker}"
  selector {
    mac = "52:54:00:b2:2f:86"
    os = "installed"
  }
  metadata {
    domain_name = "node2.example.com"
    etcd_endpoints = "node1.example.com:2379"
    etcd_on_host = "${var.experimental_self_hosted_etcd ? "false" : "true"}"
    k8s_dns_service_ip = "${var.k8s_dns_service_ip}"
    k8s_etcd_service_ip = "${var.k8s_etcd_service_ip}"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

resource "matchbox_group" "node3" {
  name = "node3"
  profile = "${module.profiles.bootkube-worker}"
  selector {
    mac = "52:54:00:c3:61:77"
    os = "installed"
  }
  metadata {
    domain_name = "node3.example.com"
    etcd_endpoints = "node1.example.com:2379"
    etcd_on_host = "${var.experimental_self_hosted_etcd ? "false" : "true"}"
    k8s_dns_service_ip = "${var.k8s_dns_service_ip}"
    k8s_etcd_service_ip = "${var.k8s_etcd_service_ip}"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
