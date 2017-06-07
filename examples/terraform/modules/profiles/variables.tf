variable "matchbox_http_endpoint" {
  type        = "string"
  description = "Matchbox HTTP read-only endpoint (e.g. http://matchbox.example.com:8080)"
}

variable "container_linux_version" {
  type        = "string"
  description = "Container Linux version of the kernel/initrd to PXE or the image to install"
}

variable "container_linux_channel" {
  type        = "string"
  description = "Container Linux channel corresponding to the container_linux_version"
}

# optional

variable "install_disk" {
  type        = "string"
  default     = "/dev/sda"
  description = "Disk device to which the install profiles should install Container Linux (e.g. /dev/sda)"
}

variable "container_linux_oem" {
  type        = "string"
  default     = ""
  description = "Specify an OEM image id to use as base for the installation (e.g. ami, vmware_raw, xen) or leave blank for the default image"
}
