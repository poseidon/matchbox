// Create common profiles
module "profiles" {
  source                  = "../profiles"
  matchbox_http_endpoint  = "${var.matchbox_http_endpoint}"
  container_linux_version = "${var.container_linux_version}"
  container_linux_channel = "${var.container_linux_channel}"
}
