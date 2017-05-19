# Self-hosted Kubernetes assets (kubeconfig, manifests)
module "bootkube" {
  source = "git::https://github.com/dghubble/bootkube-terraform.git?ref=368f09bcf189e1fe9d6185e1795e3c4bf8b1fab4"

  cluster_name                  = "${var.cluster_name}"
  api_servers                   = ["${var.k8s_domain_name}"]
  etcd_servers                  = ["http://127.0.0.1:2379"]
  asset_dir                     = "${var.asset_dir}"
  pod_cidr                      = "${var.pod_cidr}"
  service_cidr                  = "${var.service_cidr}"
  experimental_self_hosted_etcd = "${var.experimental_self_hosted_etcd}"
}
