
# CoreOS on Baremetal

Guides and a service for network booting and provisioning CoreOS clusters on virtual or physical hardware.

## Guides

* [Network Setup](network-setup.md)
* [Machine Lifecycle](machine-lifecycle.md)
* [Background: PXE Booting](network-booting.md)

## matchbox

`matchbox` is an HTTP and gRPC service that renders signed [Ignition configs](https://coreos.com/ignition/docs/latest/what-is-ignition.html), [cloud-configs](https://coreos.com/os/docs/latest/cloud-config.html), network boot configs, and metadata to machines to create CoreOS clusters. Groups match machines based on labels (e.g. MAC, UUID, stage, region) and use named Profiles for provisioning. Network boot endpoints provide PXE, iPXE, and GRUB. `matchbox` can be deployed as a binary, as an [appc](https://github.com/appc/spec) container with [rkt](https://coreos.com/rkt/docs/latest/), or as a Docker container.

* [matchbox Service](matchbox.md)
* [Profiles](matchbox.md#profiles)
* [Groups](matchbox.md#groups)
* Machine Configs
    * [Ignition](ignition.md)
    * [Cloud-Config](cloud-config.md)
* Tutorials (QEMU/KVM)
    * [matchbox with rkt](getting-started-rkt.md)
    * [matchbox with Docker](getting-started-docker.md)
* [Configuration](config.md)
* [HTTP API](api.md)
* [gRPC API](https://godoc.org/github.com/coreos/matchbox/matchbox/client)
* Installation
    * [CoreOS / Linux distros](deployment.md)
    * [rkt](deployment.md#rkt) / [docker](deployment.md#docker)
    * [Kubernetes](deployment.md#kubernetes)
* Clients
    * bootcmd CLI (POC)
    * Tectonic Installer ([guide](https://tectonic.com/enterprise/docs/latest/deployer/platform-baremetal.html), [blog](https://tectonic.com/blog/tectonic-1-3-release.html))
* Backends
    * [FileStore](matchbox.md#data)
* [Troubleshooting](troubleshooting.md)
* Going Further
    * [gRPC API Usage](config.md#grpc-api)
    * [Metadata Endpoint](api.md#metadata)
    * OpenPGP [Signing](api.md#openpgp-signatures)
    * [GRUB](grub.md)

### Examples

The [examples](https://github.com/coreos/matchbox/tree/master/examples) network boot and provision CoreOS clusters. Network boot QEMU/KVM VMs to try the examples on your Linux laptop.

* Multi-node [Kubernetes cluster](kubernetes.md)
* Multi-node [rktnetes](rktnetes.md) cluster (i.e. Kubernetes with rkt as the container runtime)
* Multi-node [self-hosted](bootkube.md) Kubernetes cluster
* [Upgrading](bootkube-upgrades.md) self-hosted Kubernetes clusters
* Multi-node etcd2 or etcd3 cluster
* Network boot or Install to Disk
* Multi-stage CoreOS installs
* [GRUB Netboot](grub.md) CoreOS
* iPXE Boot CoreOS with a root fs
* iPXE Boot CoreOS
* Lab [examples](https://github.com/dghubble/metal)
