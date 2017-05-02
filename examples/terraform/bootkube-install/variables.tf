variable "matchbox_http_endpoint" {
  type = "string"
  description = "Matchbox HTTP read-only endpoint (e.g. http://matchbox.example.com:8080)"
}

variable "matchbox_rpc_endpoint" {
  type = "string"
  description = "Matchbox gRPC API endpoint, without the protocol (e.g. matchbox.example.com:8081)"
}

variable "ssh_authorized_key" {
  type = "string"
  description = "SSH public key to set as an authorized_key on machines"
}

variable "k8s_dns_service_ip" {
  type = "string"
  default = "10.3.0.10"
  description = "Cluster DNS servce IP address passed via the Kubelet --cluster-dns flag"
}
