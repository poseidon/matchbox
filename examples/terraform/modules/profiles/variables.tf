variable "matchbox_http_endpoint" {
  type = "string"
  description = "Matchbox HTTP read-only endpoint (e.g. http://matchbox.example.com:8080)"
}

variable "coreos_version" {
  type = "string"
  description = "CoreOS kernel/initrd version to PXE boot. Must be present in matchbox assets."
}
