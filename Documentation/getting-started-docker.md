# Getting started with Docker

In this tutorial, we'll run `matchbox` on your Linux machine with Docker to network boot and provision a cluster of QEMU/KVM Container Linux machines locally. You'll be able to create Kubernetes clusters, etcd3 clusters, and test network setups.

*Note*: To provision physical machines, see [network setup](network-setup.md) and [deployment](deployment.md).

## Requirements

Install the package dependencies and start the Docker daemon.

```sh
$ # Fedora
$ sudo dnf install docker virt-install virt-manager
$ sudo systemctl start docker

$ # Debian/Ubuntu
$ # check Docker's docs to install Docker 1.8+ on Debian/Ubuntu
$ sudo apt-get install virt-manager virtinst qemu-kvm
```

Clone the [matchbox](https://github.com/coreos/matchbox) source which contains the examples and scripts.

```sh
$ git clone https://github.com/coreos/matchbox.git
$ cd matchbox
```

Download CoreOS Container Linux image assets referenced by the `etcd3` [example](../examples) to `examples/assets`.

```sh
$ ./scripts/get-coreos stable 1465.7.0 ./examples/assets
```

For development convenience, add `/etc/hosts` entries for nodes so they may be referenced by name.

```sh
# /etc/hosts
...
172.17.0.21 node1.example.com
172.17.0.22 node2.example.com
172.17.0.23 node3.example.com
```

## Containers

Run the `matchbox` and `dnsmasq` services on the `docker0` bridge. `dnsmasq` will run DHCP, DNS and TFTP services to create a suitable network boot environment. `matchbox` will serve configs to machines as they PXE boot.

The `devnet` convenience script can start these services and accepts the name of any example cluster in [examples](../examples).

```sh
$ export CONTAINER_RUNTIME=docker
$ sudo -E ./scripts/devnet create etcd3
```

Inspect the logs.

```
$ sudo -E ./scripts/devnet status
```

Take a look at the [etcd3 groups](../examples/groups/etcd3) to get an idea of how machines are mapped to Profiles. Explore some endpoints exposed by the service, say for QEMU/KVM node1.

* iPXE [http://127.0.0.1:8080/ipxe?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/ipxe?mac=52:54:00:a1:9c:ae)
* Ignition [http://127.0.0.1:8080/ignition?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/ignition?mac=52:54:00:a1:9c:ae)
* Metadata [http://127.0.0.1:8080/metadata?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/metadata?mac=52:54:00:a1:9c:ae)

### Manual

If you prefer to start the containers yourself, instead of using `devnet`,

```sh
$ sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/matchbox:Z -v $PWD/examples/groups/etcd3:/var/lib/matchbox/groups:Z quay.io/coreos/matchbox:latest -address=0.0.0.0:8080 -log-level=debug
$ sudo docker run --name dnsmasq --cap-add=NET_ADMIN -v $PWD/contrib/dnsmasq/docker0.conf:/etc/dnsmasq.conf:Z quay.io/coreos/dnsmasq -d
```

## Client VMs

Create QEMU/KVM VMs which have known hardware attributes. The nodes will be attached to the `docker0` bridge, where Docker containers run.

```sh
$ sudo ./scripts/libvirt create-docker
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

Clean up the containers and VM machines.

```sh
$ sudo -E ./scripts/devnet destroy
$ sudo ./scripts/libvirt destroy
```

## Going further

Learn more about [matchbox](matchbox.md) or explore the other [example](../examples) clusters. Try the [k8s example](bootkube.md) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl`.
