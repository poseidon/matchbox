// Fedora CoreOS profile
resource "matchbox_profile" "fedora-coreos-install" {
  name   = "worker"
  kernel = "/assets/fedora-coreos/fedora-coreos-${var.os_version}-live-kernel-x86_64"
  initrd = [
    "--name main /assets/fedora-coreos/fedora-coreos-${var.os_version}-live-initramfs.x86_64.img"
  ]

  args = [
    "initrd=main",
    "coreos.live.rootfs_url=${var.matchbox_http_endpoint}/assets/fedora-coreos/fedora-coreos-${var.os_version}-live-rootfs.x86_64.img",
    "coreos.inst.install_dev=/dev/vda",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
  ]

  raw_ignition = data.ct_config.worker.rendered
}

data "ct_config" "worker" {
  content = templatefile("fcc/fedora-coreos.yaml", {
    ssh_authorized_key = var.ssh_authorized_key
  })
  strict = true
}
