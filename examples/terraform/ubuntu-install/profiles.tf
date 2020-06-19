resource "matchbox_profile" "ubuntu-install" {
  name = "ubuntu-20.04-install"
  kernel = "http://archive.ubuntu.com/ubuntu/dists/focal/main/installer-amd64/current/legacy-images/hd-media/vmlinuz"
  initrd = [
    "http://archive.ubuntu.com/ubuntu/dists/focal/main/installer-amd64/current/legacy-images/hd-media/initrd.gz"
  ]
  args = [
    // "initrd=pxe/ubuntu/ubuntu-20.04-x86_64.img",
    "console-setup/ask_detect=false",
    "locale=en_US.UTF-8",
    "keyboard-configuration/layoutcode=us",
    "hostname=unassigned",
    "debian-installer=en_US.UTF-8",
    "kbd-chooser/method=us",
    "keyboard-configuration/layout=USA",
    "keyboard-configuration/variant=USA",
    "preseed/url=http://localhost:8080/assets/preseed.cfg"
  ]
  generic_config = "${file("./generic/preseed.cfg")}"
}