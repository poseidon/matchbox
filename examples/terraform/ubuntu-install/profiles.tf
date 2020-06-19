# Uses the ubuntu 'netboot' install option where ubuntu downloads needed packages durring install.
# You will need to manually deploy preseed.cfg to /var/lib/matchbox/assets/ubuntu/bionic/preseed.cfg
resource "matchbox_profile" "ubuntu-18.04-netboot-install" {
  name = "ubuntu-18.04-install"
  kernel = "http://archive.ubuntu.com/ubuntu/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/ubuntu-installer/amd64/linux"
  initrd = [
    "http://archive.ubuntu.com/ubuntu/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/ubuntu-installer/amd64/initrd.gz"
  ]
  # Note the configuration is intentionally split between args = [], and generic/preseed.cfg for demonstration purposes.
  # You could choose to use one or the other or a mixture of both.
  # preseed.cfg is based on the preseed.cfg examples in chef/bento github repo
  args = [
    "initrd=initrd.gz",
    "locale=en_US.UTF-8",
    "keyboard-configuration/layoutcode=us",
    "hostname=foobar",
    "preseed/url=http://matchbox.example.com:8080/assets/ubuntu/bionic/preseed.cfg"
  ]
}

# This uses cached assets generated with `scripts/get-ubuntu bionic`
# You will need to manually deploy preseed.cfg to /var/lib/matchbox/assets/ubuntu/bionic/preseed.cfg
resource "matchbox_profile" "ubuntu-18.04-asset-install" {
  name = "ubuntu-18.04-asset-install"
  kernel = "${var.matchbox_http_endpoint}/asset/ubuntu/bionic/linux"
  initrd = [
    "${var.matchbox_http_endpoint}/asset/ubuntu/bionic/initrd.gz"
  ]
  args = [
    "initrd=initrd.gz",
    "preseed/url=http://matchbox.example.com:8080/assets/ubuntu/bionic/preseed.cfg"
  ]
}