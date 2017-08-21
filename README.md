# matchbox [![Build Status](https://travis-ci.org/coreos/matchbox.svg?branch=master)](https://travis-ci.org/coreos/matchbox) [![GoDoc](https://godoc.org/github.com/coreos/matchbox?status.png)](https://godoc.org/github.com/coreos/matchbox) [![Docker Repository on Quay](https://quay.io/repository/coreos/matchbox/status "Docker Repository on Quay")](https://quay.io/repository/coreos/matchbox) [![IRC](https://img.shields.io/badge/irc-%23coreos-449FD8.svg)](https://botbot.me/freenode/coreos)

`matchbox` is a service that matches bare-metal machines (based on labels like MAC, UUID, etc.) to profiles that PXE boot and provision Container Linux clusters. Profiles specify the kernel/initrd, kernel arguments, iPXE config, GRUB config, [Container Linux Config][cl-config], or other configs a machine should use. Matchbox can be [installed](Documentation/deployment.md) as a binary, RPM, container image, or deployed on a Kubernetes cluster and it provides an authenticated gRPC API for clients like [Terraform][terraform].

* [Documentation][docs]
* [matchbox Service](Documentation/matchbox.md)
* [Profiles](Documentation/matchbox.md#profiles)
* [Groups](Documentation/matchbox.md#groups)
* Config Templates
    * [Container Linux Config][cl-config]
    * [Cloud-Config][cloud-config]
* [Configuration](Documentation/config.md)
* [HTTP API](Documentation/api.md) / [gRPC API](https://godoc.org/github.com/coreos/matchbox/matchbox/client)
* [Background: Machine Lifecycle](Documentation/machine-lifecycle.md)
* [Background: PXE Booting](Documentation/network-booting.md)

### Installation

* Installation
    * Installing on [Container Linux / other distros](Documentation/deployment.md)
    * Installing on [Kubernetes](Documentation/deployment.md#kubernetes)
    * Running with [rkt](Documentation/deployment.md#rkt) / [docker](Documentation/deployment.md#docker)
* [Network Setup](Documentation/network-setup.md)

### Tutorials

* [Getting Started](Documentation/getting-started.md) - provision Container Linux machines (beginner)
* Local QEMU/KVM
    * [matchbox with rkt](Documentation/getting-started-rkt.md)
    * [matchbox with Docker](Documentation/getting-started-docker.md)
* Clusters
  * [etcd3](Documentation/getting-started-rkt.md) - Install a 3-node etcd3 cluster
  * [Kubernetes](Documentation/bootkube.md) - Install a 3-node self-hosted Kubernetes v1.7.3 cluster
* Clusters (Terraform-based)
  * [etcd3](examples/terraform/etcd3-install/README.md) - Install a 3-node etcd3 cluster
  * [Kubernetes](examples/terraform/bootkube-install/README.md) - Install a 3-node self-hosted Kubernetes v1.7.3 cluster

## Contrib

* [dnsmasq](contrib/dnsmasq/README.md) - Run DHCP, TFTP, and DNS services with docker or rkt
* [squid](contrib/squid/README.md) - Run a transparent cache proxy
* [terraform-provider-matchbox](https://github.com/coreos/terraform-provider-matchbox) - Terraform provider plugin for Matchbox

[docs]: https://coreos.com/matchbox/docs/latest
[terraform]: https://github.com/coreos/terraform-provider-matchbox
[cl-config]: Documentation/container-linux-config.md
[cloud-config]: Documentation/cloud-config.md
