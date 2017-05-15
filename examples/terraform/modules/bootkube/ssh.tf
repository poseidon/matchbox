# Secure copy kubeconfig to all nodes to activate kubelet.service
resource "null_resource" "copy-kubeconfig" {
  count = "${length(var.controller_names) + length(var.worker_names)}"

  connection {
    type    = "ssh"
    host    = "${element(concat(var.controller_domains, var.worker_domains), count.index)}"
    user    = "core"
    timeout = "60m"
  }

  provisioner "file" {
    content     = "${module.bootkube.kubeconfig}"
    destination = "$HOME/kubeconfig"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mv /home/core/kubeconfig /etc/kubernetes/kubeconfig",
    ]
  }
}

# Secure copy bootkube assets to ONE controller and start bootkube to perform
# one-time self-hosted cluster bootstrapping.
resource "null_resource" "bootkube-start" {
  # Without depends_on, this remote-exec may start before the kubeconfig copy.
  # Terraform only does one task at a time, so it would try to bootstrap
  # Kubernetes and Tectonic while no Kubelets are running. Ensure all nodes
  # receive a kubeconfig before proceeding with bootkube and tectonic.
  depends_on = ["null_resource.copy-kubeconfig"]

  connection {
    type    = "ssh"
    host    = "${element(var.controller_domains, 0)}"
    user    = "core"
    timeout = "60m"
  }

  provisioner "file" {
    source      = "${var.asset_dir}"
    destination = "$HOME/assets"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mv /home/core/assets /opt/bootkube",
      "sudo systemctl start bootkube",
    ]
  }
}
