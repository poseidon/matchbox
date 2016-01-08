
# Examples

Examples contains `bootcfg` data directories showcasing different network-bootable baremetal clusters. These examples work with libvirt VMs setup with known hardware attributes by `scripts/libvirt`.

## Clusters

| Name       | Description | Type | Docs          |
|------------|-------------|------|---------------|
| etcd-small | Cluster with 1 etcd node, 4 proxies | cloud config | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) |
| etcd-large | Cluster with 3 etcd nodes, 2 proxies | cloud config | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) |

## Usage

First, read about the [bootcfg](../Documentation/bootcfg.md) config service. Then setup the config service and network boot environment using the [libvirt guide](../Documentation/virtual-hardware.md).

Create 5 libvirt VM nodes on the `docker0` bridge.

    ./scripts/libvirt create       # create node1 ... node5

Generally, you just need to run the `coreos/bootcfg` container and the `coreos/dnsmasq` container (which runs network boot services).

For `coreos/bootcfg`, mount one of the example directories with `-v $PWD/examples/some-example:/data:Z` and provide the `-data-path` if the mount location differs from `/data`.

Reboot to see nodes boot with PXE and create a cluster.

    ./scripts/libvirt reboot

Clean up by powering off and destroying the VMs.

    ./scripts/libvirt shutdown        # graceful
    ./scripts/libvirt poweroff        # non-graceful
    ./scripts/libvirt destroy
