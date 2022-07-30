// Default matcher group for machines
resource "matchbox_group" "default" {
  name    = "default"
  profile = matchbox_profile.flatcar-install.name
}

// Match install stage Flatcar Linux machines
resource "matchbox_group" "stage-1" {
  name    = "worker"
  profile = matchbox_profile.worker.name

  selector = {
    os = "installed"
  }
}
