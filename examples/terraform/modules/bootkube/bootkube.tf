# Self-hosted Kubernetes assets (kubeconfig, manifests)
module "bootkube" {
  source = "git::https://github.com/dghubble/bootkube-terraform.git?ref=209da6d09b1cadad655eee56d63ff0dc750c5bda"

  cluster_name                  = "${var.cluster_name}"
  api_servers                   = ["${var.k8s_domain_name}"]
  etcd_servers                  = ["http://127.0.0.1:2379"]
  asset_dir                     = "${var.asset_dir}"
  pod_cidr                      = "${var.pod_cidr}"
  service_cidr                  = "${var.service_cidr}"
  kube_apiserver_service_ip     = "${var.k8s_apiserver_service_ip}"
  kube_dns_service_ip           = "${var.k8s_dns_service_ip}"
  kube_etcd_service_ip          = "${var.k8s_etcd_service_ip}"
  kube_bootstrap_etcd_service_ip = "${var.k8s_bootstrap_etcd_service_ip}"
  experimental_self_hosted_etcd = "${var.experimental_self_hosted_etcd}"
}
