resource "matchbox_group" "foobar" {
  name = "foobar"
  profile = "${matchbox_profile.ubuntu-bionic-netboot-install.name}"

  selector = {
    mac = "00:50:56:29:54:97"
  }
  metadata = {
    fullname = "vagrant"
    password = "vagrant"
    username = "vagrant"
  }
}
