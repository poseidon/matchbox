# Matchbox

Matchbox is a service that matches bare-metal machines to profiles that PXE boot and provision clusters. Machines are matched by labels like MAC or UUID during PXE and profiles specify a kernel/initrd, iPXE config, and Container Linux or Fedora CoreOS config.

## Features

* Chainload via iPXE and match hardware labels
* Provision Container Linux and Fedora CoreOS (powered by [Ignition](https://github.com/coreos/ignition))
* Authenticated gRPC API for clients (e.g. Terraform)

## Installation

Matchbox can be installed from a binary or a container image.

* Install Matchbox on [Kubernetes](deployment.md#kubernetes), on a [Linux](deployment.md) host, or as a [container](deployment.md#docker)
* Setup a PXE-enabled [network](network-setup.md)

## Tutorials

[Getting started](getting-started.md) provisioning machines with Container Linux.

* Local QEMU/KVM
    * [matchbox with Docker](getting-started-docker.md)
* Clusters
    * [etcd3](getting-started-docker.md) - Install a 3-node etcd3 cluster
    * [etcd3](https://github.com/poseidon/matchbox/tree/master/examples/terraform/etcd3-install) - Install a 3-node etcd3 cluster (terraform-based)

## Related

* [dnsmasq](https://github.com/poseidon/matchbox/tree/master/contrib/dnsmasq) - container image to run DHCP, TFTP, and DNS services
* [terraform-provider-matchbox](https://github.com/poseidon/terraform-provider-matchbox) - Terraform provider plugin for Matchbox
* [Typhoon](https://typhoon.psdn.io/) - minimal and free Kubernetes distribution, supporting bare-metal
