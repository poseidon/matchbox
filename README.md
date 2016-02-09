
# CoreOS on Baremetal

[![Build Status](https://travis-ci.org/coreos/coreos-baremetal.svg?branch=master)](https://travis-ci.org/coreos/coreos-baremetal) [![GoDoc](https://godoc.org/github.com/coreos/coreos-baremetal?status.png)](https://godoc.org/github.com/coreos/coreos-baremetal) [![Docker Repository on Quay](https://quay.io/repository/coreos/bootcfg/status "Docker Repository on Quay")](https://quay.io/repository/coreos/bootcfg)

CoreOS on Baremetal contains guides for network booting and configuring CoreOS clusters on virtual or physical hardware.

## Guides

* [Network Booting](Documentation/network-booting.md)
* [Libvirt Guide](Documentation/virtual-hardware.md)
* [Baremetal Guide](Documentation/physical-hardware.md)

## bootcfg

`bootcfg` is a service that renders signed [Ignition configs](https://coreos.com/ignition/docs/latest/what-is-ignition.html), [cloud-configs](https://coreos.com/os/docs/latest/cloud-config.html), network boot configs, and metadata to machines based on hardware attributes (e.g. UUID, MAC) or tags (e.g. os=installed, region=us-central) to create CoreOS clusters. Network boot endpoints provide PXE, iPXE, and Pixiecore support. `bootcfg` can run as an [application container](https://github.com/appc/spec) with [rkt](https://coreos.com/rkt/docs/latest/), as a Docker container, or as a binary.

* [Getting Started with rkt](Documentation/getting-started-rkt.md)
* [Getting Started with Docker](Documentation/getting-started-docker.md)
* [bootcfg](Documentation/bootcfg.md)
    * [Ignition](Documentation/ignition.md)
    * [Cloud-Config](Documentation/cloud-config.md)
    * [Groups](Documentation/bootcfg.md#groups-and-metadata)
* [OpenPGP Signing](Documentation/openpgp.md)
* [Flags](Documentation/config.md)
* [API](Documentation/api.md)

### Examples

Use the [examples](examples) to boot machines into CoreOS clusters of higher-order systems, like Kubernetes. Quickly setup a network of virtual hardware on your Linux box for testing with the [libvirt script](scripts/libvirt).

* TLS-auth Kubernetes cluster (1 master, 1 worker, 1 etcd)
* Multi Node etcd cluster
* Install CoreOS to disk with followup Ignition stages
