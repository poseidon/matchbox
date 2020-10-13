# Matchbox

Matchbox is a service that matches bare-metal machines to profiles that PXE boot and provision clusters. Machines are matched by labels like MAC or UUID during PXE and profiles specify a kernel/initrd, iPXE config, and Ignition config.

## Features

* Chainload via iPXE and match hardware labels
* Provision Fedora CoreOS or Flatcar Linux (powered by [Ignition](https://github.com/coreos/ignition))
* Authenticated gRPC API for clients (e.g. Terraform)

## Installation

Matchbox can be installed from a binary or a container image.

* Install Matchbox on [Kubernetes](deployment.md#kubernetes), on a [Linux](deployment.md) host, or as a [container](deployment.md#docker)
* Setup a PXE-enabled [network](network-setup.md)

## Tutorials

[Getting started](getting-started.md) provisioning machines with Fedora CoreOS or Flatcar Linux.

* [Local QEMU/KVM](getting-started-docker.md)
    * Fedora CoreOS (live PXE or PXE install to disk)
    * Flatcar Linux (live PXE or PXE install to disk)
* Clusters
    * [etcd3](getting-started-docker.md) - Install a 3-node etcd3 cluster
    * [etcd3](https://github.com/poseidon/matchbox/tree/master/examples/terraform/etcd3-install) - Install a 3-node etcd3 cluster (terraform-based)

## Related

* [dnsmasq](https://github.com/poseidon/matchbox/tree/master/contrib/dnsmasq) - container image to run DHCP, TFTP, and DNS services
* [terraform-provider-matchbox](https://github.com/poseidon/terraform-provider-matchbox) - Terraform provider plugin for Matchbox
* [Typhoon](https://typhoon.psdn.io/) - minimal and free Kubernetes distribution, supporting bare-metal
