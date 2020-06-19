# Uses the ubuntu 'netboot' install option where ubuntu downloads needed packages durring install.
# You will need to manually deploy preseed.cfg to /var/lib/matchbox/assets/ubuntu/bionic/preseed.cfg
resource "matchbox_profile" "ubuntu-bionic-netboot-install" {
  name = "ubuntu-bionic-install"
  kernel = "http://archive.ubuntu.com/ubuntu/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/ubuntu-installer/amd64/linux"
  initrd = [
    "http://archive.ubuntu.com/ubuntu/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/ubuntu-installer/amd64/initrd.gz"
  ]
  # preseed.cfg is based on the preseed.cfg examples in chef/bento github repo
  # Hostnames should be assigned from DHCP https://askubuntu.com/a/675620/172065
  args = [
    "console-setup/ask_detect=false",
    "console-setup/layoutcode=us",
    "debconf/frontend=noninteractive",
    "debian-installer=en_US.UTF-8",
    "fb=false",
    "initrd=initrd.gz",
    "kbd-chooser/method=us",
    "keyboard-configuration/layout=USA",
    "keyboard-configuration/variant=USA",
    "locale=en_US.UTF-8",
    "netcfg/get_domain=unassigned-domain",
    "netcfg/get_hostname=unassigned-hostname",
    "preseed/url=${var.matchbox_http_endpoint }/generic?mac=$${mac:hexhyp}"
  ]
    generic_config = "${file("./generic/ubuntu/bionic/preseed.cfg")}"
}

# This uses cached assets generated with `scripts/get-ubuntu bionic`
resource "matchbox_profile" "ubuntu-bionic-asset-install" {
  name = "ubuntu-bionic-asset-install"
  kernel = "${var.matchbox_http_endpoint}/assets/ubuntu/bionic/linux"
  initrd = [
    "${var.matchbox_http_endpoint}/assets/ubuntu/bionic/initrd.gz"
  ]
  # preseed.cfg is based on the preseed.cfg examples in chef/bento github repo
  # Hostnames should be assigned from DHCP https://askubuntu.com/a/675620/172065
  args = [
    "console-setup/ask_detect=false",
    "console-setup/layoutcode=us",
    "debconf/frontend=noninteractive",
    "debian-installer=en_US.UTF-8",
    "fb=false",
    "initrd=initrd.gz",
    "kbd-chooser/method=us",
    "keyboard-configuration/layout=USA",
    "keyboard-configuration/variant=USA",
    "locale=en_US.UTF-8",
    "netcfg/get_domain=unassigned-domain",
    "netcfg/get_hostname=unassigned-hostname",
    "preseed/url=${var.matchbox_http_endpoint }/generic?mac=$${mac:hexhyp}"
  ]
    generic_config = "${file("./generic/ubuntu/bionic/preseed.cfg")}"

}