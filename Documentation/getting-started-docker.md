

# Getting Started with Docker

Get started with the Config service on your Linux machine with Docker. If you're ready to try [rkt](https://coreos.com/rkt/docs/latest/), see [Getting Started with rkt](getting-started-rkt.md).

In this tutorial, we'll run the Config service (`bootcfg`) to boot and provision a cluster of 5 VM machines on the `docker0` bridge. You'll be able to boot etcd clusters, Kubernetes clusters, and more, while emulating different network setups.

## Requirements

Install the dependencies. These examples have been tested on Fedora 23.

    sudo dnf install virt-install docker virt-manager
    sudo systemctl start docker

Clone the [coreos-baremetal](https://github.com/coreos/coreos-baremetal) source which contains the examples and scripts.

    git clone https://github.com/coreos/coreos-baremetal.git
    cd coreos-baremetal

Create 5 VM nodes which have known hardware attributes. The nodes will be attached to the `docker0` bridge where your containers run.

    sudo ./scripts/libvirt create-docker

Download the CoreOS PXE image assets to `assets/coreos`. The examples instruct machines to load these from the Config server, though you could change this.

    ./scripts/get-coreos

## Containers

Run the Config service (`bootcfg`). The `docker0` bridge should assign it the IP 172.17.0.2 (`sudo docker network inspect bridge`).

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/data:Z -v $PWD/assets:/assets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -config /data/etcd-docker.yaml

Take a look at [etcd-docker.yaml](../examples/etcd-docker.yaml) to get an idea of how machines are matched to specifications. Explore some endpoints port mapped to localhost:8080.

* [node1's ipxe](http://127.0.0.1:8080/ipxe?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)
* [node1's Ignition](http://127.0.0.1:8080/ignition?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)
* [node1's Metadata](http://127.0.0.1:8080/metadata?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)

Since the virtual network has no network boot services, use the `dnsmasq` container to set up an example iPXE environment which runs DHCP, DNS, and TFTP. The `dnsmasq` container can help test different network setups.

    sudo docker run --rm --cap-add=NET_ADMIN quay.io/coreos/dnsmasq -d -q --dhcp-range=172.17.0.43,172.17.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.17.0.1 --address=/bootcfg.foo/172.17.0.2

In this case, it runs a DHCP server allocating IPs to VMs between 172.17.0.43 and 172.17.0.99, resolves bootcfg.foo to 172.17.0.2 (the IP where `bootcfg` runs), and points iPXE clients to `http://bootcfg.foo:8080/boot.ipxe`.

## Verify

Reboot the VM machines and use `virt-manager` to watch the console.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt start

At this point, the VMs will PXE boot and use Ignition (preferred over cloud config) to set up a three node etcd cluster, with other nodes behaving as etcd proxies.

On VMs with autologin, check etcd2 works between different nodes.

    systemctl status etcd2
    etcdctl set /message hello
    etcdctl get /message

Clean up the VM machines.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt destroy

## Going Further

Explore the [examples](../examples). Try the `k8s-docker.yaml` config to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl`.

Add a GPG key to sign all rendered configs.

Learn more about the [config service](bootcfg.md) or adapt an example for your own [physical hardware and network](physical-hardware.md).
