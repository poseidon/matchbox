
# CoreOS on Baremetal

[![Build Status](https://travis-ci.org/coreos/coreos-baremetal.svg?branch=master)](https://travis-ci.org/coreos/coreos-baremetal) [![GoDoc](https://godoc.org/github.com/coreos/coreos-baremetal?status.png)](https://godoc.org/github.com/coreos/coreos-baremetal) [![Docker Repository on Quay](https://quay.io/repository/coreos/bootcfg/status "Docker Repository on Quay")](https://quay.io/repository/coreos/bootcfg)

CoreOS on Baremetal contains guides for network booting and configuring CoreOS clusters on virtual or physical hardware.

## Guides

* [Network Booting](Documentation/network-booting.md)

## bootcfg

`bootcfg` is a service that renders signed [Ignition configs](https://coreos.com/ignition/docs/latest/what-is-ignition.html), [cloud-configs](https://coreos.com/os/docs/latest/cloud-config.html), network boot configs, and metadata to machines based on hardware attributes (e.g. UUID, MAC) or tags (e.g. os=installed, region=us-central) to create CoreOS clusters. Network boot endpoints provide PXE, iPXE, and Pixiecore support. `bootcfg` can run as an [application container](https://github.com/appc/spec) with [rkt](https://coreos.com/rkt/docs/latest/), as a Docker container, or as a binary.

* [Getting Started with rkt](Documentation/getting-started-rkt.md)
* [Getting Started with Docker](Documentation/getting-started-docker.md)
* [bootcfg](Documentation/bootcfg.md)
    * [Groups](Documentation/bootcfg.md#groups-and-metadata)
    * [Specs](Documentation/bootcfg.md#spec)
    * [Ignition](Documentation/ignition.md)
    * [Cloud-Config](Documentation/cloud-config.md)
* [OpenPGP Signing](Documentation/openpgp.md)
* [Flags](Documentation/config.md)
* [API](Documentation/api.md)
* [Troubleshooting](Documentation/troubleshooting.md)
* [Hacking](Documentation/dev/develop.md)

### Examples

Use the [examples](examples) to boot machines into CoreOS clusters of higher-order systems, like Kubernetes. Quickly setup a network of virtual hardware on your Linux box for testing with [libvirt](scripts/README.md#libvirt).

* Multi-node Kubernetes cluster with TLS
* Multi-node etcd cluster
* Install CoreOS to disk and provision with Ignition
* GRUB Netboot CoreOS
* PXE Boot CoreOS with a root fs
* PXE Boot CoreOS
