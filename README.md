
# CoreOS on Baremetal

[![Build Status](https://travis-ci.org/coreos/coreos-baremetal.svg?branch=master)](https://travis-ci.org/coreos/coreos-baremetal) [![GoDoc](https://godoc.org/github.com/coreos/coreos-baremetal?status.png)](https://godoc.org/github.com/coreos/coreos-baremetal) [![Docker Repository on Quay](https://quay.io/repository/coreos/bootcfg/status "Docker Repository on Quay")](https://quay.io/repository/coreos/bootcfg) [![IRC](https://img.shields.io/badge/irc-%23coreos-449FD8.svg)](https://botbot.me/freenode/coreos)

Guides and a service for network booting and provisioning CoreOS clusters on virtual or physical hardware.

## Guides

* [Network Setup](Documentation/network-setup.md)
* [Machine Lifecycle](Documentation/machine-lifecycle.md)
* [Background: PXE Booting](Documentation/network-booting.md)

## bootcfg

`bootcfg` is an HTTP and gRPC service that renders signed [Ignition configs](https://coreos.com/ignition/docs/latest/what-is-ignition.html), [cloud-configs](https://coreos.com/os/docs/latest/cloud-config.html), network boot configs, and metadata to machines to create CoreOS clusters. Groups match machines based on labels (e.g. MAC, UUID, stage, region) and use named Profiles for provisioning. Network boot endpoints provide PXE, iPXE, and GRUB. `bootcfg` can be deployed as a binary, as an [appc](https://github.com/appc/spec) container with [rkt](https://coreos.com/rkt/docs/latest/), or as a Docker container.

* [bootcfg Service](Documentation/bootcfg.md)
* [Profiles](Documentation/bootcfg.md#profiles)
* [Groups](Documentation/bootcfg.md#groups-and-metadata)
* Config Templates
    * [Ignition](Documentation/ignition.md)
    * [Cloud-Config](Documentation/cloud-config.md)
* Tutorials (libvirt)
    * [bootcfg with rkt](Documentation/getting-started-rkt.md)
    * [bootcfg with Docker](Documentation/getting-started-docker.md)
* [Configuration](Documentation/config.md)
* [HTTP API](Documentation/api.md)
* [gRPC API](https://godoc.org/github.com/coreos/coreos-baremetal/bootcfg/client)
* Backends
    * [FileStore](Documentation/bootcfg.md#data)
* Installation
    * [CoreOS / Linux distros](Documentation/deployment.md)
    * [rkt](Documentation/deployment.md#rkt) / [docker](Documentation/deployment.md#docker)
    * [Kubernetes](Documentation/deployment.md#kubernetes)
* [Troubleshooting](Documentation/troubleshooting.md)
* Going Further
    * [gRPC API Usage](Documentation/config.md#grpc-api)
    * [Metadata Endpoint](Documentation/api.md#metadata)
    * OpenPGP [Signing](Documentation/api.md#openpgp-signatures)
    * [Tectonic Installer](https://tectonic.com/blog/tectonic-1-3-release.html) (graphical)

### Examples

The [examples](examples) network boot and provision CoreOS clusters. Network boot [libvirt](scripts/README.md#libvirt) VMs to try the examples on your Linux laptop.

* Multi-node [Kubernetes cluster](Documentation/kubernetes.md) with TLS
* Multi-node [self-hosted Kubernetes cluster](Documentation/bootkube.md)
* Multi-node etcd cluster
* Multi-node [Torus](Documentation/torus.md) distributed storage cluster
* Network boot or Install to Disk
* Multi-stage CoreOS installs
* [GRUB Netboot](Documentation/grub.md) CoreOS
* iPXE Boot CoreOS with a root fs
* iPXE Boot CoreOS
* Lab [examples](https://github.com/dghubble/metal)
