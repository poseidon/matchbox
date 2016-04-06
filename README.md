
# CoreOS on Baremetal

[![Build Status](https://travis-ci.org/coreos/coreos-baremetal.svg?branch=master)](https://travis-ci.org/coreos/coreos-baremetal) [![GoDoc](https://godoc.org/github.com/coreos/coreos-baremetal?status.png)](https://godoc.org/github.com/coreos/coreos-baremetal) [![Docker Repository on Quay](https://quay.io/repository/coreos/bootcfg/status "Docker Repository on Quay")](https://quay.io/repository/coreos/bootcfg)

CoreOS on Baremetal contains guides for network booting and configuring CoreOS clusters on virtual or physical hardware.

## Guides

* [Network Booting](Documentation/network-booting.md)
* [Machine Lifecycle](Documentation/machine-lifecycle.md)

## bootcfg

`bootcfg` is a HTTP and gRPC service that renders signed [Ignition configs](https://coreos.com/ignition/docs/latest/what-is-ignition.html), [cloud-configs](https://coreos.com/os/docs/latest/cloud-config.html), network boot configs, and metadata to machines based on attribute labels (e.g. UUID, MAC, stage, region) to create CoreOS clusters. Network boot endpoints provide PXE, iPXE, GRUB, and Pixiecore support. `bootcfg` can run as an [ACI](https://github.com/appc/spec) with [rkt](https://coreos.com/rkt/docs/latest/), as a Docker container, or as a binary.

* [Getting Started with rkt](Documentation/getting-started-rkt.md)
* [Getting Started with Docker](Documentation/getting-started-docker.md)
* [bootcfg Service](Documentation/bootcfg.md)
    * [Groups](Documentation/bootcfg.md#groups-and-metadata)
    * [Profiles](Documentation/bootcfg.md#profiles)
    * [Ignition](Documentation/ignition.md)
    * [Cloud-Config](Documentation/cloud-config.md)
* [Flags](Documentation/config.md)
* [API](Documentation/api.md)
* [Deployment](Documentation/deployment.md)
    * [systemd](Documentation/deployment.md#systemd)
* [Troubleshooting](Documentation/troubleshooting.md)
* Going Further
    * [OpenPGP Signing](Documentation/openpgp.md)
    * [Development](Documentation/dev/develop.md)

### Examples

Check the [examples](examples) to find Profiles for booting and provisioning machines into higher-order CoreOS clusters. Network boot [libvirt](scripts/README.md#libvirt) VMs to try the examples on your Linux laptop.

* Multi-node Kubernetes cluster with TLS (network booted or installed to disk)
* Multi-node etcd cluster (network booted or installed to disk)
* Multi-stage CoreOS installs
* GRUB Netboot CoreOS
* iPXE Boot CoreOS with a root fs
* iPXE Boot CoreOS
