
# Getting Started with rkt

In this tutorial, we'll run `matchbox` on your Linux machine with `rkt` and `CNI` to network boot and provision a cluster of QEMU/KVM CoreOS machines locally. You'll be able to create Kubernetes clustes, etcd clusters, and test network setups.

*Note*: To provision physical machines, see [network setup](network-setup.md) and [deployment](deployment.md).

## Requirements

Install [rkt](https://coreos.com/rkt/docs/latest/distributions.html) 1.8 or higher ([example script](https://github.com/dghubble/phoenix/blob/master/fedora/sources.sh)) and setup rkt [privilege separation](https://coreos.com/rkt/docs/latest/trying-out-rkt.html).

Next, install the package dependencies.

    # Fedora
    sudo dnf install virt-install virt-manager

    # Debian/Ubuntu
    sudo apt-get install virt-manager virtinst qemu-kvm systemd-container

**Note**: rkt does not yet integrate with SELinux on Fedora. As a workaround, temporarily set enforcement to permissive if you are comfortable (`sudo setenforce Permissive`). Check the rkt [distribution notes](https://github.com/coreos/rkt/blob/master/Documentation/distributions.md) or see the tracking [issue](https://github.com/coreos/rkt/issues/1727).

Clone the [matchbox](https://github.com/coreos/matchbox) source which contains the examples and scripts.

    git clone https://github.com/coreos/matchbox.git
    cd matchbox

Download CoreOS image assets referenced by the `etcd` [example](../examples) to `examples/assets`.

    ./scripts/get-coreos stable 1185.3.0 ./examples/assets

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
    "subnet": "172.18.0.0/24",
    "routes" : [ { "dst" : "0.0.0.0/0" } ]
   }
}
EOF'
```

On Fedora, add the `metal0` interface to the trusted zone in your firewall configuration.

    sudo firewall-cmd --add-interface=metal0 --zone=trusted

After a recent update, you may see a warning that NetworkManager controls the interface. Work-around this using the firewall-config GUI to add `metal0` to the trusted zone.

For development convenience, add `/etc/hosts` entries for nodes so they may be referenced by name as you would in production.

    # /etc/hosts
    ...
    172.18.0.21 node1.example.com
    172.18.0.22 node2.example.com
    172.18.0.23 node3.example.com

Trust the needed ACIs.

## Containers

Run the `matchbox` and `dnsmasq` services on the `metal0` bridge. `dnsmasq` will run DHCP, DNS, and TFTP services to create a suitable network boot environment. `matchbox` will serve provisioning configs to machines on the network which attempt to PXE boot.

Trust the needed ACIs.

    sudo rkt trust --prefix quay.io/coreos/matchbox
    sudo rkt trust --prefix quay.io/coreos/alpine-sh
    sudo rkt trust --prefix coreos.com/dnsmasq

The `devnet` wrapper script can quickly rkt run `matchbox` and `dnsmasq` in systemd transient units. Create can take the name of any example cluster in [examples](../examples).

    sudo ./scripts/devnet create etcd

Inspect the journal logs or check the status of the systemd services.

    # quick status
    sudo ./scripts/devnet status
    # tail logs
    journalctl -f -u dev-matchbox
    journalctl -f -u dev-dnsmasq

Take a look at the [etcd groups](../examples/groups/etcd) to get an idea of how machines are mapped to Profiles. Explore some endpoints exposed by the service, say for QEMU/KVM node1.

* iPXE [http://172.18.0.2:8080/ipxe?mac=52:54:00:a1:9c:ae](http://172.18.0.2:8080/ipxe?mac=52:54:00:a1:9c:ae)
* Ignition [http://172.18.0.2:8080/ignition?mac=52:54:00:a1:9c:ae](http://172.18.0.2:8080/ignition?mac=52:54:00:a1:9c:ae)
* Metadata [http://172.18.0.2:8080/metadata?mac=52:54:00:a1:9c:ae](http://172.18.0.2:8080/metadata?mac=52:54:00:a1:9c:ae)

### Manual

If you prefer to start the containers yourself, instead of using `devnet`:

    # matchbox with etcd example
    sudo rkt run --net=metal0:IP=172.18.0.2 --mount volume=data,target=/var/lib/matchbox --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/matchbox/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/matchbox:latest -- -address=0.0.0.0:8080 -log-level=debug
    # dnsmasq
    sudo rkt run coreos.com/dnsmasq:v0.3.0 --net=metal0:IP=172.18.0.3 --mount volume=config,target=/etc/dnsmasq.conf --volume config,kind=host,source=$PWD/contrib/dnsmasq/metal0.conf

If you get an error about the IP assignment, stop old pods and run garbage collection.

    sudo rkt gc --grace-period=0

## Client VMs

Create QEMU/KVM VMs which have known hardware attributes. The nodes will be attached to the `metal0` bridge, where your pods run.

    sudo ./scripts/libvirt create

You can connect to the serial console of any node. If you provisioned nodes with an SSH key, you can SSH after bring-up.

    sudo virsh console node1

You can also use `virt-manager` to watch the console.

    sudo virt-manager

Use the wrapper script to act on all nodes.

    sudo ./scripts/libvirt [start|reboot|shutdown|poweroff|destroy]

## Verify

The VMs should network boot and provision themselves into a three node etcd cluster, with other nodes behaving as etcd proxies.

The example profile added autologin so you can verify that etcd works between nodes.

    systemctl status etcd2
    etcdctl set /message hello
    etcdctl get /message

## Cleanup

Clean up the systemd units running `matchbox` and `dnsmasq`.

    sudo ./scripts/devnet destroy

Clean up VM machines.

    sudo ./scripts/libvirt destroy

Press ^] three times to stop any rkt pod.

## Going Further

Learn more about [matchbox](matchbox.md) or explore the other [example](../examples) clusters. Try the [k8s example](kubernetes.md) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl`.

