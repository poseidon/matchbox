// Default matcher group for machines
resource "matchbox_group" "default" {
  name    = "default"
  profile = matchbox_profile.fedora-coreos-install.name
}

