
# Examples

These examples network boot and provision machines into CoreOS clusters using `bootcfg`. You can re-use their profiles to provision your own physical machines.

| Name       | Description | CoreOS Version | FS | Docs | 
|------------|-------------|----------------|----|-----------|
| pxe | CoreOS via iPXE | alpha/1153.0.0 | RAM | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| grub | CoreOS via GRUB2 Netboot | alpha/1153.0.0 | RAM | NA |
| pxe-disk | CoreOS via iPXE, with a root filesystem | alpha/1153.0.0 | Disk | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| etcd | iPXE boot a 3 node etcd cluster and proxy | alpha/1153.0.0 | RAM | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) |
| etcd-install | Install a 3-node etcd cluster to disk | alpha/1153.0.0 | Disk | [reference](https://coreos.com/os/docs/latest/installing-to-disk.html) |
| etcd3 | Install a 3-node etcd3 cluster | alpha/1153.0.0 | RAM | None |
| etcd3-install | Install a 3-node etcd3 cluster to disk | alpha/1153.0.0 | Disk | None |
| k8s | Kubernetes cluster with 1 master, 2 workers, and TLS-authentication | alpha/1153.0.0 | Disk | [tutorial](../Documentation/kubernetes.md) |
| k8s-install | Kubernetes cluster, installed to disk | alpha/1153.0.0 | Disk | [tutorial](../Documentation/kubernetes.md) |
| rktnetes | Kubernetes cluster with rkt container runtime, 1 master, workers, TLS auth (experimental) | alpha/1153.0.0 | Disk | None |
| rktnetes-install | Kubernetes cluster with rkt container runtime, installed to disk (experimental) | alpha/1153.0.0 | Disk | None |
| bootkube | iPXE boot a self-hosted Kubernetes cluster (with bootkube) | alpha/1153.0.0 | Disk | [tutorial](../Documentation/bootkube.md) |
| bootkube-install | Install a self-hosted Kubernetes cluster (with bootkube) | alpha/1153.0.0 | Disk | [tutorial](../Documentation/bootkube.md) |
| torus | Torus distributed storage | alpha/1153.0.0 | Disk | [tutorial](../Documentation/torus.md) |

## Tutorials

Get started running `bootcfg` on your Linux machine to network boot and provision clusters of VMs or physical hardware.

* Getting Started
	* [bootcfg with rkt](../Documentation/getting-started-rkt.md)
	* [bootcfg with Docker](../Documentation/getting-started-docker.md)
* [Kubernetes (static manifests)](../Documentation/kubernetes.md)
* [Kubernetes (self-hosted)](../Documentation/bootkube.md)
* [Torus Storage](../Documentation/torus.md)
* [Lab Examples](https://github.com/dghubble/metal)

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
