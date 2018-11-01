// Kubernetes cluster
module "cluster" {
  source = "git::https://github.com/poseidon/typhoon//bare-metal/container-linux/kubernetes?ref=v1.10.3"

  providers = {
    local = "local.default"
    null = "null.default"
    template = "template.default"
    tls = "tls.default"
  }

  # bare-metal
  cluster_name            = "${var.cluster_name}"
  matchbox_http_endpoint  = "${var.matchbox_http_endpoint}"
  os_channel              = "${var.os_channel}"
  os_version              = "${var.os_version}"

  # configuration
  k8s_domain_name    = "${var.k8s_domain_name}"
  ssh_authorized_key = "${var.ssh_authorized_key}"
  asset_dir          = "${var.asset_dir}"

  # machines
  controller_names   = "${var.controller_names}"
  controller_macs    = "${var.controller_macs}"
  controller_domains = "${var.controller_domains}"
  worker_names       = "${var.worker_names}"
  worker_macs        = "${var.worker_macs}"
  worker_domains     = "${var.worker_domains}"

  # optional
  networking          = "${var.networking}"
  cached_install      = "${var.cached_install}"
  install_disk        = "${var.install_disk}"
  container_linux_oem = "${var.container_linux_oem}"
  kernel_args         = "${var.kernel_args}"
}
