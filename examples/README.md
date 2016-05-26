
# Examples

These examples network boot and provision VMs into CoreOS clusters using `bootcfg`.

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

## Tutorials

Get started running `bootcfg` on your Linux machine to network boot and provision clusters of VMs or physical hardware.

* [bootcfg with rkt](../Documentation/getting-started-rkt.md)
* [bootcfg with Docker](../Documentation/getting-started-docker.md)
* [Kubernetes v1.2.4](../Documentation/kubernetes.md)
* [Self-hosted Kubernetes](../Documentation/bootkube.md) (experimental)

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
