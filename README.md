
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
* [OpenPGP Signing](Documentation/openpgp.md)
* [Flags](Documentation/config.md)
* [API](Documentation/api.md)
* [Deployment](Documentation/deployment.md)
    * [systemd](Documentation/deployment.md#systemd)
* [Troubleshooting](Documentation/troubleshooting.md)
* [Hacking](Documentation/dev/develop.md)

### Examples

Check the [examples](examples) to find Profiles for booting machines into higher-order CoreOS clusters. Quickly network boot some [libvirt](scripts/README.md#libvirt) VMs to try the flow on your Linux machine.

* Multi-node Kubernetes cluster with TLS
* Multi-node etcd cluster
* Multi-stage CoreOS install to disk and provision with Ignition
* GRUB Netboot CoreOS
* PXE Boot CoreOS with a root fs
* PXE Boot CoreOS
