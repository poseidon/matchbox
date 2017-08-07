# etcd3

The `etcd3-install` example shows how to use matchbox to network boot and provision 3-node etcd3 cluster on bare-metal in an automated way.

## Requirements

Follow the getting started [tutorial](../../../Documentation/getting-started.md) to learn about matchbox and set up an environment that meets the requirements:

* Matchbox v0.6+ [installation](../../../Documentation/deployment.md) with gRPC API enabled
* Matchbox provider credentials `client.crt`, `client.key`, and `ca.crt`
* PXE [network boot](../../../Documentation/network-setup.md) environment
* Terraform v0.9+ and [terraform-provider-matchbox](https://github.com/coreos/terraform-provider-matchbox) installed locally on your system
* 3 machines with known DNS names and MAC addresses

If you prefer to provision QEMU/KVM VMs on your local Linux machine, set up the matchbox [development environment](../../../Documentation/getting-started-rkt.md).

```sh
sudo ./scripts/devnet create
```

## Usage

Clone the [matchbox](https://github.com/coreos/matchbox) project and take a look at the cluster examples.

```sh
$ git clone https://github.com/coreos/matchbox.git
$ cd matchbox/examples/terraform/etcd3-install
```

Copy the `terraform.tfvars.example` file to `terraform.tfvars`. Ensure `provider.tf` references your matchbox credentials.

```hcl
matchbox_http_endpoint = "http://matchbox.example.com:8080"
matchbox_rpc_endpoint = "matchbox.example.com:8081"
ssh_authorized_key = "ADD ME"
```

Configs in `etcd3-install` configure the matchbox provider, define profiles (e.g. `cached-container-linux-install`, `etcd3`), and define 3 groups which match machines by MAC address to a profile. These resources declare that the machines should PXE boot, install Container Linux to disk, and provision themselves into peers in a 3-node etcd3 cluster.

Note: The `cached-container-linux-install` profile will PXE boot and install Container Linux from matchbox [assets](https://github.com/coreos/matchbox/blob/master/Documentation/api.md#assets). If you have not populated the assets cache, use the `container-linux-install` profile to use public images (slower).

### Optional

You may set certain optional variables to override defaults.

```hcl
# install_disk = "/dev/sda"
# container_linux_oem = ""
```

## Apply

Fetch the [profiles](../README.md#modules) Terraform [module](https://www.terraform.io/docs/modules/index.html) which let's you use common machine profiles maintained in the matchbox repo (like `etcd3`).

```sh
$ terraform get
```

Plan and apply to create the resoures on Matchbox.

```sh
$ terraform plan
Plan: 10 to add, 0 to change, 0 to destroy.
$ terraform apply
Apply complete! Resources: 10 added, 0 changed, 0 destroyed.
```

## Machines

Power on each machine (with PXE boot device on next boot). Machines should network boot, install Container Linux to disk, reboot, and provision themselves as a 3-node etcd3 cluster. 

```sh
$ ipmitool -H node1.example.com -U USER -P PASS chassis bootdev pxe
$ ipmitool -H node1.example.com -U USER -P PASS power on
```

For local QEMU/KVM development, create the QEMU/KVM VMs.

```sh
$ sudo ./scripts/libvirt create
$ sudo ./scripts/libvirt [start|reboot|shutdown|poweroff|destroy]
```

## Verify

Verify each node is running etcd3 (i.e. etcd-member.service).

```sh
$ ssh core@node1.example.com
$ systemctl status etcd-member
```

Verify that etcd3 peers are healthy and communicating.

```sh
$ etcdctl cluster-health
$ etcdctl set /message hello
$ etcdctl get /message
```

## Going Further

Learn more about [matchbox](../../../Documentation/matchbox.md) or explore the other [example](../) clusters.
