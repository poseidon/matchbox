# Getting Started with rkt

Get started with `bootcfg` on your Linux machine with rkt, CNI, and appc. 

In this tutorial, we'll run `bootcfg` to boot and provision a cluster of four VM machines on a CNI bridge (`metal0`). You'll be able to boot etcd clusters, Kubernetes clusters, and more, while testing different network setups.

## Requirements

Install [rkt](https://github.com/coreos/rkt/releases) and [acbuild](https://github.com/appc/acbuild/releases) from the latest releases. For rkt, see the setup and privilege separation [docs](https://coreos.com/rkt/docs/latest/trying-out-rkt.html). For acbuild:

    tar xzvf acbuild.tar.gz
    sudo ln -s /path/to/acbuild /usr/local/bin/acbuild

Install package dependencies.

    # Fedora
    sudo dnf install virt-install virt-manager

    # Debian/Ubuntu
    sudo apt-get install virt-manager virtinst qemu-kvm systemd-container

**Note**: Currently, rkt does not integrate with SELinux on Fedora. As a workaround, temporarily set enforcement to permissive if you are comfortable (`sudo setenforce Permissive`). Check the rkt [distribution notes](https://github.com/coreos/rkt/blob/master/Documentation/distributions.md) or track the [issue](https://github.com/coreos/rkt/issues/1727).

Clone the [coreos-baremetal](https://github.com/coreos/coreos-baremetal) source which contains the examples and scripts.

    git clone https://github.com/coreos/coreos-baremetal.git
    cd coreos-baremetal

Download the CoreOS PXE image assets to `assets/coreos`. The examples instruct machines to load these from `bootcfg`.

    ./scripts/get-coreos
    ./scripts/get-coreos channel version

Define the `metal0` virtual bridge with [CNI](https://github.com/appc/cni).

```bash
sudo mkdir -p /etc/rkt/net.d
sudo bash -c 'cat > /etc/rkt/net.d/20-metal.conf << EOF
{
  "name": "metal0",
  "type": "bridge",
  "bridge": "metal0",
  "isGateway": true,
  "ipMasq": true,
  "ipam": {
    "type": "host-local",
    "subnet": "172.15.0.0/16",
    "routes" : [ { "dst" : "0.0.0.0/0" } ]
   }
}
EOF'
```

On Fedora, add the `metal0` interface to the trusted zone in your firewall configuration.

    sudo firewall-cmd --add-interface=metal0 --zone=trusted

## Application Container

Run `bootcfg` on the `metal0` network, with a known IP we'll have DNS point to.

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
    sudo virt-manager

## Network

Since the virtual network has no network boot services, use the `dnsmasq` ACI to create an iPXE network boot environment which runs DHCP, DNS, and TFTP. The `dnsmasq` container can help test different network setups.

Build the `dnsmasq.aci` ACI.

    cd contrib/dnsmasq
    ./get-tftp-files
    sudo ./build-aci

Run `dnsmasq.aci` to create a DHCP and TFTP server pointing to config server.

    sudo rkt --insecure-options=image run dnsmasq.aci --net=metal0:IP=172.15.0.3 -- -d -q --dhcp-range=172.15.0.50,172.15.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.15.0.1 --address=/bootcfg.foo/172.15.0.2

In this case, dnsmasq runs a DHCP server allocating IPs to VMs between 172.15.0.50 and 172.15.0.99, resolves bootcfg.foo to 172.15.0.2 (the IP where `bootcfg` runs), and points iPXE clients to `http://bootcfg.foo:8080/boot.ipxe`.

## Verify

Reboot the VM machines and use `virt-manager` to watch the console.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt start

At this point, the VMs will PXE boot and use Ignition (preferred over cloud config) to set up a three node etcd cluster, with other nodes behaving as etcd proxies.

The example spec added autologin so you can check that etcd works between nodes.

    systemctl status etcd2
    etcdctl set /message hello
    etcdctl get /message

Press ^] three times to stop a rkt pod. Clean up the VM machines.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt destroy
    sudo ./scripts/libvirt delete-disks

## Going Further

Explore the [examples](../examples). Try the `k8s-rkt.yaml` [example](../examples/README.md#kubernetes) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl`.

Learn more about [bootcfg](bootcfg.md), enable [OpenPGP signing](openpgp.md), or adapt an example for your own [physical hardware](physical-hardware.md) and network.
