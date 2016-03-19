# Getting Started with rkt

In this tutorial, we'll run `bootcfg` on your Linux machinem with `rkt` and `CNI`, to network boot and provision a cluster of CoreOS machines. You'll be able to create Kubernetes clustes, etcd clusters, or just install CoreOS and test network setups locally.

## Requirements

Install [rkt](https://github.com/coreos/rkt/releases) and [acbuild](https://github.com/appc/acbuild/releases) from the latest releases ([example script](https://github.com/dghubble/phoenix/blob/master/scripts/fedora/sources.sh)). Optionally setup rkt [privilege separation](https://coreos.com/rkt/docs/latest/trying-out-rkt.html).

Install package dependencies.

    # Fedora
    sudo dnf install virt-install virt-manager

    # Debian/Ubuntu
    sudo apt-get install virt-manager virtinst qemu-kvm systemd-container

**Note**: rkt does not yet integrate with SELinux on Fedora. As a workaround, temporarily set enforcement to permissive if you are comfortable (`sudo setenforce Permissive`). Check the rkt [distribution notes](https://github.com/coreos/rkt/blob/master/Documentation/distributions.md) or see the tracking [issue](https://github.com/coreos/rkt/issues/1727).

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

## Containers

#### Latest

Run the latest ACI with rkt.

    sudo rkt --insecure-options=image run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/var/bootcfg --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/etc/bootcfg --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -log-level=debug -config /etc/bootcfg/etcd-rkt.yaml

Note: The insecure flag is needed for this case, since [Quay.io](https://quay.io/repository/coreos/bootcfg) serves ACIs coverted from Docker images (docker2aci) and Docker images don't support signatures.

#### Release

Alternately, run a recent tagged and signed [release](https://github.com/coreos/coreos-baremetal/releases). Trust the [CoreOS App Signing Key](https://coreos.com/dist/pubkeys/app-signing-pubkey.gpg) for image signature verification.

    sudo rkt trust --prefix coreos.com/bootcfg
    # gpg key fingerprint is: 18AD 5014 C99E F7E3 BA5F  6CE9 50BD D3E0 FC8A 365E
    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/assets --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/data --volume data,kind=host,source=$PWD/examples coreos.com/bootcfg:v0.2.0 -- -address=0.0.0.0:8080 -log-level=debug -config /data/etcd-rkt.yaml

If you get an error about the IP assignment, garbage collect old pods.

    sudo rkt gc --grace-period=0
    ./scripts/rkt-gc-force

Take a look at [etcd-rkt.yaml](../examples/etcd-rkt.yaml) to get an idea of how machines are matched to profiles. Explore some endpoints exposed by the service.

* [node1's ipxe](http://172.15.0.2:8080/ipxe?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)
* [node1's Ignition](http://172.15.0.2:8080/ignition?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)
* [node1's Metadata](http://172.15.0.2:8080/metadata?uuid=16e7d8a7-bfa9-428b-9117-363341bb330b)

## Network

Since the virtual network has no network boot services, use the `dnsmasq` ACI to create an iPXE network boot environment which runs DHCP, DNS, and TFTP.

Trust the [CoreOS App Signing Key](https://coreos.com/dist/pubkeys/app-signing-pubkey.gpg).

    sudo rkt trust --prefix coreos.com/dnsmasq
    # gpg key fingerprint is: 18AD 5014 C99E F7E3 BA5F  6CE9 50BD D3E0 FC8A 365E

Run the `coreos.com/dnsmasq` ACI with rkt.

    sudo rkt run coreos.com/dnsmasq:v0.2.0 --net=metal0:IP=172.15.0.3 -- -d -q --dhcp-range=172.15.0.50,172.15.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.15.0.1 --address=/bootcfg.foo/172.15.0.2

In this case, dnsmasq runs a DHCP server allocating IPs to VMs between 172.15.0.50 and 172.15.0.99, resolves `bootcfg.foo` to 172.15.0.2 (the IP where `bootcfg` runs), and points iPXE clients to `http://bootcfg.foo:8080/boot.ipxe`.

## Client VMs

Create VM nodes which have known hardware attributes. The nodes will be attached to the `metal0` bridge where your pods run.

    sudo ./scripts/libvirt create-rkt
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

Press ^] three times to stop a rkt pod. Clean up the VM machines.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt destroy
    sudo ./scripts/libvirt delete-disks

## Going Further

Explore the [examples](../examples). Try the `k8s-rkt.yaml` [example](../examples/README.md#kubernetes) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl`.

Learn more about [bootcfg](bootcfg.md), enable [OpenPGP signing](openpgp.md), or adapt an example for your own [physical hardware](physical-hardware.md) and network.
