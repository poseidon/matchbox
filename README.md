
# CoreOS on Baremetal

[![Build Status](https://travis-ci.org/coreos/coreos-baremetal.svg?branch=master)](https://travis-ci.org/coreos/coreos-baremetal) [![GoDoc](http://godoc.org/github.com/coreos/coreos-baremetal?status.png)](http://godoc.org/github.com/coreos/coreos-baremetal) [![Docker Repository on Quay](https://quay.io/repository/coreos/bootcfg/status "Docker Repository on Quay")](https://quay.io/repository/coreos/bootcfg)

CoreOS on Baremetal contains guides for network booting and configuring CoreOS clusters on virtual or physical hardware.

## Guides

* [Network Booting](Documentation/network-booting.md)
* [Config Service](Documentation/bootcfg.md)
* [Libvirt Guide](Documentation/virtual-hardware.md)
* [Baremetal Guide](Documentation/physical-hardware.md)

## Config Service

The config service renders signed [Ignition](https://coreos.com/ignition/docs/latest/what-is-ignition.html) configs, [Cloud-Init](https://github.com/coreos/coreos-cloudinit) configs, and metadata to machines based on hardware attributes (e.g. UUID, MAC) or arbitrary tags (e.g. os=installed, region=us-central). Network boot endpoints provide PXE, iPXE, and Pixiecore support.

* [Getting Started with rkt](Documentation/getting-started-rkt.md)
* [Getting Started with Docker](Documentation/getting-started-docker.md)
* [API](Documentation/api.md)
* [Flags](Documentation/config.md)

### Examples

Boot machines into CoreOS clusters of higher-order systems like Kubernetes or etcd according to the declarative [examples](examples). Use the [libvirt script](scripts/libvirt) to quickly setup a network of virtual hardware on your Linux box.

* Multi Node etcd Cluster
* Kubernetes Cluster (1 master, 1 worker, 1 etcd)
