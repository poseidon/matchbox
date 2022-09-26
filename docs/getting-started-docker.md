# Getting started with Docker

In this tutorial, we'll run `matchbox` on a Linux machine with Docker to network boot and provision local QEMU/KVM machines as Fedora CoreOS or Flatcar Linux machines. You'll be able to test network setups and Ignition provisioning.

!!! note
    To provision physical machines, see [network setup](network-setup.md) and [deployment](deployment.md).

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

Clone the [matchbox](https://github.com/poseidon/matchbox) source which contains the examples and scripts.

```sh
$ git clone https://github.com/poseidon/matchbox.git
$ cd matchbox
```

Download Fedora CoreOS or Flatcar Linux image assets to `examples/assets`.

```sh
$ ./scripts/get-fedora-coreos stable 36.20220906.3.2 ./examples/assets
$ ./scripts/get-flatcar stable 3227.2.0 ./examples/assets
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

The `devnet` convenience script can start these services and accepts the name of any example in [examples](https://github.com/poseidon/matchbox/tree/master/examples).

```sh
$ sudo ./scripts/devnet create fedora-coreos
```

Inspect the logs.

```
$ sudo ./scripts/devnet status
```

Inspect the examples and Matchbox endpoints to see how machines (e.g. node1 with MAC `52:54:00:a1:9c:ae`) are mapped to Profiles, and therefore iPXE and Ignition configs.

* iPXE [http://127.0.0.1:8080/ipxe?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/ipxe?mac=52:54:00:a1:9c:ae)
* Ignition [http://127.0.0.1:8080/ignition?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/ignition?mac=52:54:00:a1:9c:ae)
* Metadata [http://127.0.0.1:8080/metadata?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/metadata?mac=52:54:00:a1:9c:ae)

### Manual

If you prefer to start the containers yourself, instead of using `devnet`,

```sh
$ sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/matchbox:Z -v $PWD/examples/groups/fedora-coreos:/var/lib/matchbox/groups:Z quay.io/poseidon/matchbox:latest -address=0.0.0.0:8080 -log-level=debug
$ sudo docker run --name dnsmasq --cap-add=NET_ADMIN -v $PWD/contrib/dnsmasq/docker0.conf:/etc/dnsmasq.conf:Z quay.io/poseidon/dnsmasq -d
```

## Client VMs

Create QEMU/KVM VMs which have known hardware attributes. The nodes will be attached to the `docker0` bridge, where Docker containers run.

```sh
$ sudo ./scripts/libvirt create
```

If you provisioned nodes with an SSH key, you can SSH after bring-up.

```sh
$ ssh core@node1.example.com
```

If you set a `console=ttyS0` kernel arg, you can connect to the serial console of any node (ctrl+] to exit).

```
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

The VMs should network boot and provision themselves as declared.

```
cat /etc/os-release
```

## Clean up

Clean up the containers and VM machines.

```sh
$ sudo ./scripts/devnet destroy
$ sudo ./scripts/libvirt destroy
```

## Going Further

Learn more about [matchbox](matchbox.md) or explore the other [examples](https://github.com/poseidon/matchbox/tree/master/examples).

Try different examples and Ignition declarations:

* Declare an SSH authorized public key (see examples README)
* Declare a systemd unit
* Declare file or directory content

