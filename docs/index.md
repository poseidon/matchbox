# Matchbox

Matchbox is a service that matches bare-metal machines to profiles that PXE boot and provision clusters. Machines are matched by labels like MAC or UUID during PXE and profiles specify a kernel/initrd, iPXE config, and Ignition config.

## Features

* Chainload via iPXE and match hardware labels
* Provision Fedora CoreOS or Flatcar Linux (powered by [Ignition](https://github.com/coreos/ignition))
* Authenticated gRPC API for clients (e.g. Terraform)

## Installation

Matchbox can be installed from a binary or a container image.

* Install Matchbox as a [binary](deployment.md#matchbox-binary), as a [container image](deployment.md#container-image), or on [Kubernetes](deployment.md#kubernetes)
* Setup a PXE-enabled [network](network-setup.md)

## Tutorials

Start provisioning machines with Fedora CoreOS or Flatcar Linux.

* [Terraform Usage](getting-started.md)
    * Fedora CoreOS (live PXE or PXE install to disk)
    * Flatcar Linux (live PXE or PXE install to disk)
* [Local QEMU/KVM](getting-started-docker.md)
    * Fedora CoreOS (live PXE or PXE install to disk)
    * Flatcar Linux (live PXE or PXE install to disk)

## Related

* [dnsmasq](https://github.com/poseidon/matchbox/tree/master/contrib/dnsmasq) - container image to run DHCP, TFTP, and DNS services
* [terraform-provider-matchbox](https://github.com/poseidon/terraform-provider-matchbox) - Terraform provider plugin for Matchbox
* [Typhoon](https://typhoon.psdn.io/) - minimal and free Kubernetes distribution, supporting bare-metal
