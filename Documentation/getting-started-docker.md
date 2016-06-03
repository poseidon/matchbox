
# Getting Started with Docker

In this tutorial, we'll run `bootcfg` on your Linux machine with Docker to network boot and provision a cluster of CoreOS machines locally. You'll be able to create Kubernetes clustes, etcd clusters, and test network setups.

If you're ready to try [rkt](https://coreos.com/rkt/docs/latest/), see [Getting Started with rkt](getting-started-rkt.md).

## Requirements

Install the package dependencies and start the Docker daemon.

    # Fedora
    sudo dnf install docker virt-install virt-manager
    sudo systemctl start docker

    # Debian/Ubuntu
    # check Docker's docs to install Docker 1.8+ on Debian/Ubuntu
    sudo apt-get install virt-manager virtinst qemu-kvm

Clone the [coreos-baremetal](https://github.com/coreos/coreos-baremetal) source which contains the examples and scripts.

    git clone https://github.com/coreos/coreos-baremetal.git
    cd coreos-baremetal

Download CoreOS image assets referenced by the `etcd-docker` [example](../examples) to `examples/assets`.

    ./scripts/get-coreos
    ./scripts/get-coreos channel version   # examples pin a required version

## Containers

#### Latest

Run the latest `bootcfg` Docker image from `quay.io/coreos/bootcfg` with the `etcd-docker` example. The container should receive the IP address 172.17.0.2 on the `docker0` bridge.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd-docker:/var/lib/bootcfg/groups:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

#### Release

Alternately, run the most recent tagged `bootcfg` [release](https://github.com/coreos/coreos-baremetal/releases).

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd-docker:/var/lib/bootcfg/groups:Z quay.io/coreos/bootcfg:v0.3.0 -address=0.0.0.0:8080 -log-level=debug

Take a look at the [etcd groups](../examples/groups/etcd-docker) to get an idea of how machines are mapped to Profiles. Explore some endpoints port mapped to localhost:8080.

* [node1's ipxe](http://127.0.0.1:8080/ipxe?mac=52:54:00:a1:9c:ae)
* [node1's Ignition](http://127.0.0.1:8080/ignition?mac=52:54:00:a1:9c:ae)
* [node1's Metadata](http://127.0.0.1:8080/metadata?mac=52:54:00:a1:9c:ae)

## Network

Since the virtual network has no network boot services, use the `dnsmasq` image to create an iPXE network boot environment which runs DHCP, DNS, and TFTP.

    sudo docker run --rm --cap-add=NET_ADMIN quay.io/coreos/dnsmasq -d -q --dhcp-range=172.17.0.43,172.17.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.17.0.1 --address=/bootcfg.foo/172.17.0.2

In this case, dnsmasq runs a DHCP server allocating IPs to VMs between 172.17.0.43 and 172.17.0.99, resolves `bootcfg.foo` to 172.17.0.2 (the IP where `bootcfg` runs), and points iPXE clients to `http://bootcfg.foo:8080/boot.ipxe`.

## Client VMs

Create VM nodes which have known hardware attributes. The nodes will be attached to the `docker0` bridge where Docker's containers run.

    sudo ./scripts/libvirt create-docker
    sudo virt-manager

You can use `virt-manager` to watch the console and reboot VM machines with

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt start

## Verify

The VMs should network boot and provision themselves into a three node etcd cluster, with other nodes behaving as etcd proxies.

The example profile added autologin so you can verify that etcd works between nodes.

    systemctl status etcd2
    etcdctl set /message hello
    etcdctl get /message
    fleetctl list-machines

Clean up the VM machines.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt destroy

## Going Further

Explore the [examples](../examples). Try the [k8s-docker example](kubernetes.md) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl` ([docs](../examples/README.md#kubernetes)).

Learn more about [bootcfg](bootcfg.md) or adapt an example for your own [physical hardware](physical-hardware.md) and network.
