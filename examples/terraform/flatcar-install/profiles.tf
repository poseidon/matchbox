// Create a flatcar-install profile
resource "matchbox_profile" "flatcar-install" {
  name   = "flatcar-install"
  kernel = "http://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe.vmlinuz"
  initrd = [
    "http://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe_image.cpio.gz",
  ]

  args = [
    "initrd=flatcar_production_pxe_image.cpio.gz",
    "flatcar.config.url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "flatcar.first_boot=yes",
    "console=tty0",
    "console=ttyS0",
  ]

  container_linux_config = file("./clc/flatcar-install.yaml")
}

// Profile to set an SSH authorized key on first boot from disk
resource "matchbox_profile" "worker" {
  name                   = "worker"
  container_linux_config = file("./clc/flatcar.yaml")
}
