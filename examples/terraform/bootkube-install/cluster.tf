// Self-hosted Kubernetes cluster
module "cluster" {
  source = "../modules/bootkube"

  matchbox_http_endpoint = "${var.matchbox_http_endpoint}"
  ssh_authorized_key     = "${var.ssh_authorized_key}"

  cluster_name            = "${var.cluster_name}"
  container_linux_channel = "${var.container_linux_channel}"
  container_linux_version = "${var.container_linux_version}"

  # Machines
  controller_names   = "${var.controller_names}"
  controller_macs    = "${var.controller_macs}"
  controller_domains = "${var.controller_domains}"
  worker_names       = "${var.worker_names}"
  worker_macs        = "${var.worker_macs}"
  worker_domains     = "${var.worker_domains}"

  # bootkube assets
  k8s_domain_name = "${var.k8s_domain_name}"
  asset_dir       = "${var.asset_dir}"

  # Optional
  cached_install                = "${var.cached_install}"
  install_disk                  = "${var.install_disk}"
  container_linux_oem           = "${var.container_linux_oem}"
  experimental_self_hosted_etcd = "${var.experimental_self_hosted_etcd}"
}
