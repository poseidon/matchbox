# matchbox [![Build Status](https://github.com/poseidon/matchbox/workflows/test/badge.svg)](https://github.com/poseidon/matchbox/actions?query=workflow%3Atest+branch%3Amaster) [![GoDoc](https://godoc.org/github.com/poseidon/matchbox?status.svg)](https://godoc.org/github.com/poseidon/matchbox) [![Quay](https://img.shields.io/badge/container-quay-green)](https://quay.io/repository/poseidon/matchbox)

`matchbox` is a service that matches bare-metal machines to profiles that PXE boot and provision clusters. Machines are matched by labels like MAC or UUID during PXE and profiles specify a kernel/initrd, iPXE config, and Container Linux or Fedora CoreOS config.

## Features

* Chainload via iPXE and match hardware labels
* Provision Fedora CoreOS and Flatcar Linux (powered by [Ignition](https://github.com/coreos/ignition))
* Authenticated gRPC API for clients (e.g. Terraform)

## Documentation

* [Docs](https://matchbox.psdn.io/)
* [Configuration](docs/config.md)
* [HTTP API](docs/api-http.md) / [gRPC API](docs/grpc-api.md)

## Installation

Matchbox can be installed from a binary or a container image.

* Install Matchbox on [Kubernetes](docs/deployment.md#kubernetes), on a [Linux](docs/deployment.md) host, or as a [container](docs/deployment.md#docker)
* Setup a PXE-enabled [network](docs/network-setup.md)

## Tutorials

[Getting started](docs/getting-started.md) provisioning machines with Flatcar Linux.

* Local QEMU/KVM
    * [matchbox with Docker](docs/getting-started-docker.md)
* Clusters
    * [etcd3](docs/getting-started-docker.md) - Install a 3-node etcd3 cluster
    * [etcd3](https://github.com/poseidon/matchbox/tree/master/examples/terraform/etcd3-install) - Install a 3-node etcd3 cluster (terraform-based)

## Contrib

* [dnsmasq](contrib/dnsmasq/README.md) - Run DHCP, TFTP, and DNS services as a container
* [terraform-provider-matchbox](https://github.com/poseidon/terraform-provider-matchbox) - Terraform provider plugin for Matchbox
