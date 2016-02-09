
# CoreOS on Baremetal

[![Build Status](https://travis-ci.org/coreos/coreos-baremetal.svg?branch=master)](https://travis-ci.org/coreos/coreos-baremetal) [![GoDoc](https://godoc.org/github.com/coreos/coreos-baremetal?status.png)](https://godoc.org/github.com/coreos/coreos-baremetal) [![Docker Repository on Quay](https://quay.io/repository/coreos/bootcfg/status "Docker Repository on Quay")](https://quay.io/repository/coreos/bootcfg)

CoreOS on Baremetal contains guides for network booting and configuring CoreOS clusters on virtual or physical hardware.

## Guides

* [Network Booting](Documentation/network-booting.md)
* [Config Service](Documentation/bootcfg.md)
* [Libvirt Guide](Documentation/virtual-hardware.md)
* [Baremetal Guide](Documentation/physical-hardware.md)

## Config Service

`bootcfg` is a service that renders signed [Ignition configs](https://coreos.com/ignition/docs/latest/what-is-ignition.html), [Cloud configs](https://github.com/coreos/coreos-cloudinit), network boot configs, and metadata to machines based on hardware attributes (e.g. UUID, MAC) or tags (e.g. os=installed, region=us-central) to create CoreOS clusters. Network boot endpoints provide PXE, iPXE, and Pixiecore support. `bootcfg` can run as an [application container](https://github.com/appc/spec) with [rkt](https://coreos.com/rkt/docs/latest/), as a Docker container, or as a binary.

* [Getting Started with rkt](Documentation/getting-started-rkt.md)
* [Getting Started with Docker](Documentation/getting-started-docker.md)
* [OpenPGP Signing](Documentation/openpgp.md)
* [Flags](Documentation/config.md)
* [API](Documentation/api.md)

### Examples

Boot machines into CoreOS clusters of higher-order systems like Kubernetes or etcd according to the declarative [examples](examples). Use the [libvirt script](scripts/libvirt) to quickly setup a network of virtual hardware on your Linux box.

* Multi Node etcd Cluster
* Kubernetes Cluster (1 master, 1 worker, 1 etcd)
