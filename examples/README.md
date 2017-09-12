# Examples

Matchbox automates network booting and provisioning of clusters. These examples show how to use matchbox on-premise or locally with [QEMU/KVM](scripts/README.md#libvirt).

## Terraform Examples

These examples use [Terraform](https://www.terraform.io/intro/) as a client to Matchbox.

| Name                          | Description                   |
|-------------------------------|-------------------------------|
| [simple-install](terraform/simple-install) | Install Container Linux with an SSH key |
| [etcd3-install](terraform/etcd3-install) | Install a 3-node etcd3 cluster |
| [bootkube-install](terraform/bootkube-install) | Install a 3-node Kubernetes v1.7.5 cluster |

### Customization

You are encouraged to look through the examples and Terraform modules. Implement your own profiles or package them as modules to meet your needs. We've just provided a starting point. Learn more about [matchbox](../Documentation/matchbox.md) and [Container Linux configs](../Documentation/container-linux-config.md).

## Manual Examples

These examples mount raw Matchbox objects into a Matchbox server's `/var/lib/matchbox/` directory.

| Name       | Description | CoreOS Container Linux Version | FS | Docs | 
|------------|-------------|----------------|----|-----------|
| simple | CoreOS Container Linux with autologin, using iPXE | stable/1465.7.0 | RAM | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| simple-install | CoreOS Container Linux Install, using iPXE | stable/1465.7.0 | RAM | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| grub | CoreOS Container Linux via GRUB2 Netboot | stable/1465.7.0 | RAM | NA |
| etcd3 | PXE boot a 3-node etcd3 cluster with proxies | stable/1465.7.0 | RAM | None |
| etcd3-install | Install a 3-node etcd3 cluster to disk | stable/1465.7.0 | Disk | None |
| bootkube | PXE boot a 3-node Kubernetes v1.7.5 cluster | stable/1465.7.0 | Disk | [tutorial](../Documentation/bootkube.md) |
| bootkube-install | Install a 3-node Kubernetes v1.7.5 cluster | stable/1465.7.0 | Disk | [tutorial](../Documentation/bootkube.md) |

### Customization

#### Autologin

Example profiles pass the `coreos.autologin` kernel argument. This skips the password prompt for development and troubleshooting and should be removed **before production**.

## SSH Keys

Example groups allow `ssh_authorized_keys` to be added for the `core` user as metadata. You might also include this directly in your Ignition.

    # /var/lib/matchbox/groups/default.json
    {
        "name": "Example Machine Group",
        "profile": "pxe",
        "metadata": {
            "ssh_authorized_keys": ["ssh-rsa pub-key-goes-here"]
        }
    }

#### Conditional Variables

**"pxe"**

Some examples check the `pxe` variable to determine whether to create a `/dev/sda1` filesystem and partition for PXEing with `root=/dev/sda1` ("pxe":"true") or to write files to the existing filesystem on `/dev/disk/by-label/ROOT` ("pxe":"false").
