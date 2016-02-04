# Getting Started with rkt

Get started with the Config service on your Linux machine with rkt, CNI, and appc. 

In this tutorial, we'll run the Config service (`bootcfg`) to boot and provision a cluster of four VM machines on a CNI bridge (`metal0`). You'll be able to boot etcd clusters, Kubernetes clusters, and more, while emulating different network setups.

## Requirements

**Note**: Currently, rkt and the Fedora/RHEL/CentOS SELinux policies aren't supported. See the [issue](https://github.com/coreos/rkt/issues/1727) tracking the work and policy changes. To test these examples on your laptop, set SELinux enforcement to permissive if you are comfortable (`sudo setenforce 0`). Enable it again when you are finished.

Install [rkt](https://github.com/coreos/rkt/releases), [acbuild](https://github.com/appc/acbuild), and package dependencies. These examples have been tested on Fedora 23.

    sudo dnf install virt-install virt-manager

Clone the [coreos-baremetal](https://github.com/coreos/coreos-baremetal) source which contains the examples and scripts.

    git clone https://github.com/coreos/coreos-baremetal.git
    cd coreos-baremetal

Download the CoreOS PXE image assets to `assets/coreos`. The examples instruct machines to load these from the Config server, though you could change this.

    ./scripts/get-coreos

Define the `metal0` virtual bridge with [CNI](https://github.com/appc/cni).

```bash
sudo bash -c 'cat > /etc/rkt/net.d/20-metal.conf << EOF{
  "name": "metal0",
  "type": "bridge",
  "bridge": "metal0",
  "isGateway": true,
  "ipMasq": true,
  "ipam": {
    "type": "host-local",
    "subnet": "172.15.0.0/16",
    "routes" : [ { "dst" : "172.15.0.0/16" } ]
   }
}
EOF'
```

## Application Container

Run the Config service (`bootcfg`) on the `metal0` network, with a known IP we'll use in later steps with DNS.

    sudo rkt --insecure-options=image fetch docker://quay.io/coreos/bootcfg

Currently, the insecure flag is needed since Docker images do not support signature verification. We'll ship an ACI soon to address this.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/assets --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/data --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg -- -address=0.0.0.0:8080 -log-level=debug -config /data/etcd-rkt.yaml

If you get an error about the IP assignment, garbage collect old pods.

    sudo rkt gc --grace-period=0

Take a look at [etcd-rkt.yaml](../examples/etcd-rkt.yaml) to get an idea of how machines are matched to specifications. Explore some endpoints exposed by the service.

* [node1's ipxe](http://172.15.0.2:8080/ipxe?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)
* [node1's Ignition](http://172.15.0.2:8080/ignition?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)
* [node1's Metadata](http://172.15.0.2:8080/metadata?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)

## Client VMs

Create four VM nodes which have known hardware attributes. The nodes will be attached to the `metal0` bridge where your pods run.

    sudo ./scripts/libvirt create-rkt

In your Firewall Configuration, add `metal0` as a trusted interface.

## Network

Since the virtual network has no network boot services, use the `dnsmasq` ACI to set up an example iPXE environment which runs DHCP, DNS, and TFTP. The `dnsmasq` container can help test different network setups.

Build the `dnsmasq.aci` ACI.

    cd contrib/dnsmasq
    sudo ./build-aci

Run `dnsmasq.aci` to create a DHCP and TFTP server pointing to config server.

    sudo rkt --insecure-options=image run dnsmasq.aci --net=metal0:IP=172.15.0.3 -- -d -q --dhcp-range=172.15.0.50,172.15.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.15.0.1 --address=/bootcfg.foo/172.15.0.2

In this case, dnsmasq runs a DHCP server allocating IPs to VMs between 172.15.0.50 and 172.15.0.99, resolves bootcfg.foo to 172.15.0.2 (the IP where `bootcfg` runs), and points iPXE clients to `http://bootcfg.foo:8080/boot.ipxe`.

## Verify

Reboot the VM machines and use `virt-manager` to watch the console.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt start

At this point, the VMs will PXE boot and use Ignition (preferred over cloud config) to set up a three node etcd cluster, with other nodes behaving as etcd proxies.

On VMs with autologin, check etcd2 works between different nodes.

    systemctl status etcd2
    etcdctl set /message hello
    etcdctl get /message

Press ^] three times to stop a rkt pod. Clean up the VM machines.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt destroy

## Going Further

Explore the [examples](../examples). Try the `k8s-rkt.yaml` [example](../examples/README.md#kubernetes) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl`.

Add a GPG key to sign all rendered configs.

Learn more about the [config service](bootcfg.md) or adapt an example for your own [physical hardware and network](physical-hardware.md).
