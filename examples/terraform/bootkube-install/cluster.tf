// Kubernetes cluster
module "cluster" {
  source = "git::https://github.com/poseidon/typhoon//bare-metal/container-linux/kubernetes?ref=v1.13.2"

  providers = {
    local    = "local.default"
    null     = "null.default"
    template = "template.default"
    tls      = "tls.default"
  }

  # bare-metal
  cluster_name           = "example"
  matchbox_http_endpoint = "${var.matchbox_http_endpoint}"
  os_channel             = "coreos-stable"
  os_version             = "1967.3.0"

  # configuration
  k8s_domain_name    = "cluster.example.com"
  ssh_authorized_key = "${var.ssh_authorized_key}"
  asset_dir          = "assets"
  cached_install     = "true"

  # machines
  controller_names   = ["node1"]
  controller_macs    = ["52:54:00:a1:9c:ae"]
  controller_domains = ["node1.example.com"]

  worker_names = [
    "node2",
    "node3",
  ]

  worker_macs = [
    "52:54:00:b2:2f:86",
    "52:54:00:c3:61:77",
  ]

  worker_domains = [
    "node2.example.com",
    "node3.example.com",
  ]
}
