output "coreos-install" {
  value = "${matchbox_profile.coreos-install.name}"
}

output "etcd3" {
  value = "${matchbox_profile.etcd3.name}"
}

output "etcd3-gateway" {
  value = "${matchbox_profile.etcd3-gateway.name}"
}

output "bootkube-controller" {
  value = "${matchbox_profile.bootkube-controller.name}"
}

output "bootkube-worker" {
  value = "${matchbox_profile.bootkube-worker.name}"
}
