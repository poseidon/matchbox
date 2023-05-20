# matchbox
[![GoDoc](https://pkg.go.dev/badge/github.com/poseidon/matchbox.svg)](https://pkg.go.dev/github.com/poseidon/matchbox)
[![Quay](https://img.shields.io/badge/container-quay-green)](https://quay.io/repository/poseidon/matchbox)
[![Workflow](https://github.com/poseidon/matchbox/actions/workflows/build.yaml/badge.svg)](https://github.com/poseidon/matchbox/actions/workflows/build.yaml?query=branch%3Amain)
![Downloads](https://img.shields.io/github/downloads/poseidon/matchbox/total)
[![Sponsors](https://img.shields.io/github/sponsors/poseidon?logo=github)](https://github.com/sponsors/poseidon)
[![Mastodon](https://img.shields.io/badge/follow-news-6364ff?logo=mastodon)](https://fosstodon.org/@poseidon)

`matchbox` is a service that matches bare-metal machines to profiles that PXE boot and provision clusters. Machines are matched by labels like MAC or UUID during PXE and profiles specify a kernel/initrd, iPXE config, and Ignition config.

## Features

* Chainload via iPXE and match hardware labels
* Provision Fedora CoreOS or Flatcar Linux (powered by [Ignition](https://github.com/coreos/ignition))
* Authenticated gRPC API for clients (e.g. Terraform)

## Documentation

* [Docs](https://matchbox.psdn.io/)
* [Configuration](docs/config.md)
* [HTTP API](docs/api-http.md) / [gRPC API](docs/api-grpc.md)

## Installation

Matchbox can be installed from a binary or a container image.

* Install Matchbox as a [binary](docs/deployment.md#matchbox-binary), as a [container image](docs/deployment.md#container-image), or on [Kubernetes](docs/deployment.md#kubernetes)
* Setup a PXE-enabled [network](docs/network-setup.md)

## Tutorials

Start provisioning machines with Fedora CoreOS or Flatcar Linux.

* [Terraform Usage](docs/getting-started.md)
  * Fedora CoreOS (PXE install to disk)
  * Flatcar Linux (PXE install to disk)
* [Local QEMU/KVM](docs/getting-started-docker.md)
    * Fedora CoreOS (live PXE or PXE install to disk)
    * Flatcar Linux (live PXE or PXE install to disk)

## Contrib

* [dnsmasq](contrib/dnsmasq/README.md) - Run DHCP, TFTP, and DNS services as a container
* [terraform-provider-matchbox](https://github.com/poseidon/terraform-provider-matchbox) - Terraform provider plugin for Matchbox
