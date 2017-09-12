# Getting started with rkt

In this tutorial, we'll run `matchbox` on your Linux machine with `rkt` and `CNI` to network boot and provision a cluster of QEMU/KVM Container Linux machines locally. You'll be able to create Kubernetes clustes, etcd3 clusters, and test network setups.

*Note*: To provision physical machines, see [network setup](network-setup.md) and [deployment](deployment.md).

## Requirements

Install [rkt](https://coreos.com/rkt/docs/latest/distributions.html) 1.12.0 or higher ([example script](https://github.com/dghubble/phoenix/blob/master/fedora/sources.sh)) and setup rkt [privilege separation](https://coreos.com/rkt/docs/latest/trying-out-rkt.html).

Next, install the package dependencies.

```sh
# Fedora
$ sudo dnf install virt-install virt-manager

# Debian/Ubuntu
$ sudo apt-get install virt-manager virtinst qemu-kvm systemd-container
```

**Note**: rkt does not yet integrate with SELinux on Fedora. As a workaround, temporarily set enforcement to permissive if you are comfortable (`sudo setenforce Permissive`). Check the rkt [distribution notes](https://github.com/coreos/rkt/blob/master/Documentation/distributions.md) or see the tracking [issue](https://github.com/coreos/rkt/issues/1727).

Clone the [matchbox](https://github.com/coreos/matchbox) source which contains the examples and scripts.

```sh
$ git clone https://github.com/coreos/matchbox.git
$ cd matchbox
```

Download CoreOS Container Linux image assets referenced by the `etcd3` [example](../examples) to `examples/assets`.

```sh
$ ./scripts/get-coreos stable 1465.7.0 ./examples/assets
```

## Network

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

```sh
$ sudo firewall-cmd --add-interface=metal0 --zone=trusted
$ sudo firewall-cmd --add-interface=metal0 --zone=trusted --permanent
```

For development convenience, you may wish to add `/etc/hosts` entries for nodes to refer to them by name.

```
# /etc/hosts
...
172.18.0.21 node1.example.com
172.18.0.22 node2.example.com
172.18.0.23 node3.example.com
```

## Containers

Run the `matchbox` and `dnsmasq` services on the `metal0` bridge. `dnsmasq` will run DHCP, DNS, and TFTP services to create a suitable network boot environment. `matchbox` will serve configs to machinesas they PXE boot.

The `devnet` convenience script can rkt run these services in systemd transient units and accepts the name of any example cluster in [examples](../examples).

```sh
$ export CONTAINER_RUNTIME=rkt
$ sudo -E ./scripts/devnet create etcd3
```

Inspect the journal logs.

```
$ sudo -E ./scripts/devnet status
$ journalctl -f -u dev-matchbox
$ journalctl -f -u dev-dnsmasq
```

Take a look at the [etcd3 groups](../examples/groups/etcd3) to get an idea of how machines are mapped to Profiles. Explore some endpoints exposed by the service, say for QEMU/KVM node1.

* iPXE [http://172.18.0.2:8080/ipxe?mac=52:54:00:a1:9c:ae](http://172.18.0.2:8080/ipxe?mac=52:54:00:a1:9c:ae)
* Ignition [http://172.18.0.2:8080/ignition?mac=52:54:00:a1:9c:ae](http://172.18.0.2:8080/ignition?mac=52:54:00:a1:9c:ae)
* Metadata [http://172.18.0.2:8080/metadata?mac=52:54:00:a1:9c:ae](http://172.18.0.2:8080/metadata?mac=52:54:00:a1:9c:ae)

### Manual

If you prefer to start the containers yourself, instead of using `devnet`,

```sh
sudo rkt run --net=metal0:IP=172.18.0.2 \
  --mount volume=data,target=/var/lib/matchbox \
  --volume data,kind=host,source=$PWD/examples \
  --mount volume=groups,target=/var/lib/matchbox/groups \
  --volume groups,kind=host,source=$PWD/examples/groups/etcd3 \
  quay.io/coreos/matchbox:v0.6.1 -- -address=0.0.0.0:8080 -log-level=debug
```
```sh
sudo rkt run --net=metal0:IP=172.18.0.3 \
  --dns=host \
  --mount volume=config,target=/etc/dnsmasq.conf \
  --volume config,kind=host,source=$PWD/contrib/dnsmasq/metal0.conf \
  quay.io/coreos/dnsmasq:v0.4.1 \
  --caps-retain=CAP_NET_ADMIN,CAP_NET_BIND_SERVICE,CAP_SETGID,CAP_SETUID,CAP_NET_RAW
```

If you get an error about the IP assignment, stop old pods and run garbage collection.

```sh
$ sudo rkt gc --grace-period=0
```

## Client VMs

Create QEMU/KVM VMs which have known hardware attributes. The nodes will be attached to the `metal0` bridge, where your pods run.

```sh
$ sudo ./scripts/libvirt create
```

You can connect to the serial console of any node. If you provisioned nodes with an SSH key, you can SSH after bring-up.

```sh
$ sudo virsh console node1
```

You can also use `virt-manager` to watch the console.

```sh
$ sudo virt-manager
```

Use the wrapper script to act on all nodes.

```sh
$ sudo ./scripts/libvirt [start|reboot|shutdown|poweroff|destroy]
```

## Verify

The VMs should network boot and provision themselves into a three node etcd3 cluster, with other nodes behaving as etcd3 gateways.

The example profile added autologin so you can verify that etcd3 works between nodes.

```sh
$ systemctl status etcd-member
$ etcdctl set /message hello
$ etcdctl get /message
```

## Clean up

Clean up the systemd units running `matchbox` and `dnsmasq`.

```sh
$ sudo -E ./scripts/devnet destroy
```

Clean up VM machines.

```sh
$ sudo ./scripts/libvirt destroy
```

Press ^] three times to stop any rkt pod.

## Going further

Learn more about [matchbox](matchbox.md) or explore the other [example](../examples) clusters. Try the [k8s example](bootkube.md) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl`.
