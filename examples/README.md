
# Examples

These examples network boot and provision machines into CoreOS clusters using `bootcfg`. You can re-use their profiles to provision your own physical machines.

| Name       | Description | CoreOS Version | FS | Docs | 
|------------|-------------|----------------|----|-----------|
| pxe | CoreOS via iPXE | alpha/1053.2.0 | RAM | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| grub | CoreOS via GRUB2 Netboot | alpha/1053.2.0 | RAM | NA |
| pxe-disk | CoreOS via iPXE, with a root filesystem | alpha/1053.2.0 | Disk | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| etcd, etcd-docker | iPXE boot a 3 node etcd cluster and proxy | alpha/1053.2.0 | RAM | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) |
| etcd-install | Install a 3-node etcd cluster to disk | alpha/1053.2.0 | Disk | [reference](https://coreos.com/os/docs/latest/installing-to-disk.html) |
| k8s, k8s-docker | Kubernetes cluster with 1 master and 2 workers, TLS-authentication | alpha/1053.2.0 | Disk | [tutorial](../Documentation/kubernetes.md) |
| k8s-install | Install a Kubernetes cluster to disk (1 master) | alpha/1053.2.0 | Disk | [tutorial](../Documentation/kubernetes.md) |
| bootkube | iPXE boot a self-hosted Kubernetes cluster (with bootkube) | alpha/1053.2.0 | Disk | [tutorial](../Documentation/bootkube.md) |
| bootkube-install | Install a self-hosted Kubernetes cluster (with bootkube) | alpha/1053.2.0 | Disk | [tutorial](../Documentation/bootkube.md) |
| torus | CoreOS Torus distributed storage | alpha/1053.2.0 | Disk | [tutorial](../Documentation/torus.md) |

## Tutorials

Get started running `bootcfg` on your Linux machine to network boot and provision clusters of VMs or physical hardware.

* [bootcfg with rkt](../Documentation/getting-started-rkt.md)
* [bootcfg with Docker](../Documentation/getting-started-docker.md)
* [Static Kubernetes](../Documentation/kubernetes.md) v1.2.4
* [Self-hosted Kubernetes](../Documentation/bootkube.md) v1.3.0-beta.2
* [Torus Storage](..Documentation/torus.md)

## Experimental

These examples demonstrate booting and provisioning various (often experimental) CoreOS clusters. They have **NOT** been hardened for production yet. You should write or adapt Ignition configs to suit your needs and hardware.

## SSH Keys

Most examples allow `ssh_authorized_keys` to be added for the `core` user as machine group metadata.

    # /var/lib/bootcfg/groups/default.json
    {
        "name": "Example Machine Group",
        "profile": "pxe",
        "metadata": {
            "ssh_authorized_keys": ["ssh-rsa pub-key-goes-here"]
        }
    }

## Conditional Variables

### "pxe"

Some examples check the `pxe` variable to determine whether to create a `/dev/sda1` filesystem and partition for PXEing with `root=/dev/sda1` ("pxe":"true") or to write files to the existing filesystem on `/dev/disk/by-label/ROOT` ("pxe":"false").

### "skip_networkd"

Some examples (mainly Kubernetes examples) check the `skip_networkd` variable to determine whether to skip configuring networkd. When `true`, the default networkd config is used, which uses DCHP to setup networking. Use this if you've pre-configured static IP mappings for Kubernetes nodes in your DHCP server. Otherwise, `networkd_address`, `networkd_dns`, and `networkd_gateway` machine metadata are used to populate a networkd configuration on each host.
