resource "matchbox_profile" "ubuntu-18.04-install" {
  name = "ubuntu-18.04-install"
  kernel = "http://archive.ubuntu.com/ubuntu/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/ubuntu-installer/amd64/linux"
  initrd = [
    "http://archive.ubuntu.com/ubuntu/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/ubuntu-installer/amd64/initrd.gz"
  ]
  args = [
    "initrd=initrd.gz",
    "locale=en_US.UTF-8",
    "keyboard-configuration/layoutcode=us",
    "hostname=foobar",
    "preseed/url=http://matchbox.example.com:8080/assets/preseed.cfg"
  ]
  generic_config = "${file("./generic/preseed.cfg")}"
}