// Fedora CoreOS profile
resource "matchbox_profile" "fedora-coreos-install" {
  name   = "worker"
  kernel = "https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-kernel-x86_64"
  initrd = [
    "https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-initramfs.x86_64.img",
    "https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-rootfs.x86_64.img"
  ]

  args = [
    "rd.neednet=1",
    "coreos.inst.install_dev=/dev/sda",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "coreos.inst.image_url=https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-metal.x86_64.raw.xz",
    "console=tty0",
    "console=ttyS0",
  ]

  raw_ignition = data.ct_config.worker-ignition.rendered
}

data "ct_config" "worker-ignition" {
  content = data.template_file.worker-config.rendered
  strict  = true
}

data "template_file" "worker-config" {
  template = file("fcc/fedora-coreos.yaml")
  vars = {
    ssh_authorized_key = var.ssh_authorized_key
  }
}


