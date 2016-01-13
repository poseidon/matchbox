
# CoreOS on Baremetal

[![Build Status](https://travis-ci.org/coreos/coreos-baremetal.svg?branch=master)](https://travis-ci.org/coreos/coreos-baremetal) [![GoDoc](http://godoc.org/github.com/coreos/coreos-baremetal?status.png)](http://godoc.org/github.com/coreos/coreos-baremetal) [![Docker Repository on Quay](https://quay.io/repository/coreos/bootcfg/status "Docker Repository on Quay")](https://quay.io/repository/coreos/bootcfg)

CoreOS on Baremetal contains guides for network booting and configuring CoreOS clusters on virtual or physical hardware.

## Guides

* [Network Booting](Documentation/network-booting.md)
* [Config Service](Documentation/bootcfg.md)
* [Libvirt Guide](Documentation/virtual-hardware.md)
* [Baremetal Guide](Documentation/physical-hardware.md)

## Config Service

The config service provides network boot (PXE, iPXE, Pixiecore), [Ignition](https://coreos.com/ignition/docs/latest/what-is-ignition.html), and [Cloud-Init](https://github.com/coreos/coreos-cloudinit) configs to machines based on hardware attributes (e.g. UUID, MAC, hostname) or free-form tag matchers.

* [API](Documentation/api.md)
* [Flags](Documentation/config.md)

## Examples

Get started with the declarative [examples](examples) which network boot different CoreOS clusters. Use the [libvirt script](scripts/libvirt) to quickly setup a network of virtual hardware on your Linux box.

* Single Node etcd Cluster
* Multi Node etcd Cluster
* Kubernetes Cluster (1 master, 1 worker, 1 etcd)
