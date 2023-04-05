// Create a flatcar-install profile
resource "matchbox_profile" "flatcar-install" {
  name   = "flatcar-install"
  kernel = "/assets/flatcar/3227.2.0/flatcar_production_pxe.vmlinuz"
  initrd = [
    "/assets/flatcar/3227.2.0/flatcar_production_pxe_image.cpio.gz",
  ]

  args = [
    "initrd=flatcar_production_pxe_image.cpio.gz",
    "flatcar.config.url=${var.matchbox_http_endpoint}/ignition?uuid=${uuid}&mac=${mac:hexhyp}",
    "flatcar.first_boot=yes",
  ]

  raw_ignition = data.ct_config.install.rendered
}

data "ct_config" "install" {
  content = templatefile("butane/flatcar-install.yaml", {
    matchbox_http_endpoint = var.matchbox_http_endpoint
    ssh_authorized_key     = var.ssh_authorized_key
  })
  strict = true
}


// Profile to set an SSH authorized key on first boot from disk
resource "matchbox_profile" "worker" {
  name         = "worker"
  raw_ignition = data.ct_config.worker.rendered
}

data "ct_config" "worker" {
  content = templatefile("butane/flatcar.yaml", {
    ssh_authorized_key = var.ssh_authorized_key
  })
  strict = true
}
