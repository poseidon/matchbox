
# Getting Started with rkt

In this tutorial, we'll run `bootcfg` on your Linux machine with `rkt` and `CNI` to network boot and provision a cluster of CoreOS machines locally. You'll be able to create Kubernetes clustes, etcd clusters, and test network setups.

*Note*: To provision physical machines, see [network setup](network-setup.md) and [deployment](deployment.md).

## Requirements

Install [rkt](https://coreos.com/rkt/docs/latest/distributions.html) 1.8 or higher ([example script](https://github.com/dghubble/phoenix/blob/master/scripts/fedora/sources.sh)) and setup rkt [privilege separation](https://coreos.com/rkt/docs/latest/trying-out-rkt.html).

Next, install the package dependencies.

    # Fedora
    sudo dnf install virt-install virt-manager

    # Debian/Ubuntu
    sudo apt-get install virt-manager virtinst qemu-kvm systemd-container

**Note**: rkt does not yet integrate with SELinux on Fedora. As a workaround, temporarily set enforcement to permissive if you are comfortable (`sudo setenforce Permissive`). Check the rkt [distribution notes](https://github.com/coreos/rkt/blob/master/Documentation/distributions.md) or see the tracking [issue](https://github.com/coreos/rkt/issues/1727).

Clone the [coreos-baremetal](https://github.com/mikeynap/coreos-baremetal) source which contains the examples and scripts.

    git clone https://github.com/mikeynap/coreos-baremetal.git
    cd coreos-baremetal

Download CoreOS image assets referenced by the `etcd` [example](../examples) to `examples/assets`.

    ./scripts/get-coreos alpha 1053.2.0 ./examples/assets

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

Run the latest `bootcfg` ACI with rkt and the `etcd` example.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -log-level=debug

If you get an error about the IP assignment, stop old pods and run garbage collection.

    sudo rkt gc --grace-period=0

Take a look at the [etcd groups](../examples/groups/etcd) to get an idea of how machines are mapped to Profiles. Explore some endpoints exposed by the service.

* [node1's ipxe](http://172.15.0.2:8080/ipxe?mac=52:54:00:a1:9c:ae)
* [node1's Ignition](http://172.15.0.2:8080/ignition?mac=52:54:00:a1:9c:ae)
* [node1's Metadata](http://172.15.0.2:8080/metadata?mac=52:54:00:a1:9c:ae)

## Network

Since the virtual network has no network boot services, use the `dnsmasq` ACI to create an iPXE network boot environment which runs DHCP, DNS, and TFTP.

Trust the [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/).

    sudo rkt trust --prefix coreos.com/dnsmasq
    # gpg key fingerprint is: 18AD 5014 C99E F7E3 BA5F  6CE9 50BD D3E0 FC8A 365E

Run the `coreos.com/dnsmasq` ACI with rkt.

    sudo rkt run coreos.com/dnsmasq:v0.3.0 --net=metal0:IP=172.15.0.3 -- -d -q --dhcp-range=172.15.0.50,172.15.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.15.0.1 --address=/bootcfg.foo/172.15.0.2

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
    fleetctl list-machines

Press ^] three times to stop a rkt pod. Clean up the VM machines.

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt destroy

## Going Further

Learn more about [bootcfg](bootcfg.md) or explore the other [example](../examples) clusters. Try the [k8s example](kubernetes.md) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl` ([docs](../examples/README.md#kubernetes)).

