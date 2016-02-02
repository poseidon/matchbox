
# Examples

Examples contains Config Service data directories showcasing different network-bootable baremetal clusters. These examples work with the libvirt VMs created by `scripts/libvirt` and the [Libvirt Guide](../Documentation/virtual-hardware.md).

| Name       | Description |  Docs          |
|------------|-------------|----------------|
| etcd-small | Cluster with 1 etcd node, 4 proxies | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) |
| etcd-large | Cluster with 3 etcd nodes, 2 proxies | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) |
| kubernetes | Kubernetes cluster with 1 master, 1 worker, 1 dedicated etcd node | [reference](https://github.com/coreos/coreos-kubernetes) |

## Experimental

These examples are experimental and have **NOT** been hardened for production. They are designed to demonstrate booting and configuring CoreOS clusters, especially locally.

## Virtual Hardware

Get started on your Linux machine by creating a network of virtual hardware. Install `libvirt` and `virt-manager` and clone the source.

    # Fedora/RHEL
    dnf install virt-manager

Create 5 libvirt VM nodes, which should be enough for any of the examples. The `scripts/libvirt` script will create 5 VM nodes with known hardware attributes, on the `docker0` bridge network.

    # clone the source
    git clone https://github.com/coreos/coreos-baremetal.git
    # create 5 nodes
    ./scripts/libvirt create

The nodes can be conveniently managed together.

    ./scripts/libvirt reboot
    ./scripts/libvirt start
    ./scripts/libvirt shutdown        # graceful
    ./scripts/libvirt poweroff        # non-graceful
    ./scripts/libvirt destroy

## Physical Hardware

You can use these examples to provision experimental physical clusters. You'll have to edit the IPs in the examples which reference the config service and cluster nodes to suit your network and update the `config.yaml` to match the attribtues of your hardware. For Kubernetes, be sure to generate TLS assets with the new set of IPs. Read the [Baremetal Guide](../Documentation/physical-hardware.md) for details.

## Config Service

The Config service matches machines to boot configurations, ignition configs, and cloud configs. It optionally serves OS images and other assets.

Let's run the config service on the virtual network.

    docker pull quay.io/coreos/bootcfg:latest

Run the command for the example you wish to use.

**etcd-small Cluster**

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/examples/etcd-small:/data:Z -v $PWD/assets:/assets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

**etcd-large Cluster**

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/examples/etcd-large:/data:Z -v $PWD/assets:/assets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

**Kubernetes Cluster**

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/examples/kubernetes:/data:Z -v $PWD/assets:/assets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

The mounted data directory (e.g. `-v $PWD/examples/etcd-small:/data:Z`) depends on the example you wish to run.

## Assets

The examples require the CoreOS stable PXE kernel and initrd images be served by the Config service. Run `get-coreos` to download those images to `assets`.

    ./scripts/get-coreos

The following examples require a few additional assets be downloaded or generated. Acquire those assets before continuing.

* [Kubernetes Cluster](kubernetes)

## Network Environment

Run an iPXE setup with DHCP and TFTP on the virtual network on your machine similar to what would be present on a real network. This allocates IP addresses to VM hosts, points PXE booting clients to the config service, and chainloads iPXE.

    sudo docker run --rm --cap-add=NET_ADMIN quay.io/coreos/dnsmasq -d -q --dhcp-range=172.17.0.43,172.17.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.17.0.1 --address=/bootcfg.foo/172.17.0.2

You may need to update your firewall to allow DHCP and TFTP services. If your config service container has a different IP address or subnet, the IPs in examples will require some adjustments to match your local setup.

## Boot

Reboot the nodes to PXE boot them into your new cluster!

    ./scripts/libvirt reboot
    # if nodes are in a non-booted state
    ./scripts/libvirt poweroff
    ./scripts/libvirt start

The examples use autologin for debugging and checking that nodes were setup correctly depending on the example. If something goes wrong, see [troubleshooting](../Documentation/troubleshooting.md).

If everything works, congratulations! Stay tuned for developments.

## Further Reading

See the [libvirt guide](../Documentation/virtual-hardware.md) or [baremetal guide](../Documentation/physical-hardware.md) for more information.

