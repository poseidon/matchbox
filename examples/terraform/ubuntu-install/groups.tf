resource "matchbox_group" "vmware" {
  name = "vmware"
  profile = "${matchbox_profile.ubuntu-18.04-netboot-install.name}"

  selector = {
    mac = "00:50:56:29:54:97"
  }
  metadata = {
    foo = "bar"
  }
}
