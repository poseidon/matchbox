
# Examples

These examples network boot and provision VMs into CoreOS clusters using `bootcfg`.

| Name       | Description | CoreOS Version | FS | Reference | 
|------------|-------------|----------------|----|-----------|
| pxe | CoreOS via iPXE | alpha/983.0.0 | RAM | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| grub | CoreOS via GRUB2 Netboot | alpha/983.0.0 | RAM | NA |
| pxe-disk | CoreOS via iPXE, with a root filesystem | alpha/983.0.0 | Disk | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| etcd, etcd-docker | Cluster with 3 etcd nodes, 2 proxies | alpha/983.0.0 | RAM | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) |
| etcd-install | Install a 3-node etcd cluster to disk | alpha/983.0.0 | Disk | [reference](https://coreos.com/os/docs/latest/installing-to-disk.html) |
| k8s, k8s-docker | Kubernetes cluster with 1 master and 2 workers, TLS-authentication | alpha/983.0.0 | Disk | [reference](https://github.com/coreos/coreos-kubernetes) |
| k8s-install | Install a Kubernetes cluster to disk (1 master) | alpha/983.0.0 | Disk | [reference](https://github.com/coreos/coreos-kubernetes) |

## Experimental

These CoreOS clusters are experimental and have **NOT** been hardened for production yet. They demonstrate Ignition and cloud-init provisioning of higher order clusters.

## Getting Started

Get started running the `bootcfg` on your Linux machine to boot clusters of libvirt PXE VMs.

* [Getting Started with rkt](../Documentation/getting-started-rkt.md)
* [Getting Started with Docker](../Documentation/getting-started-docker.md)

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

## Kubernetes

The Kubernetes examples create Kubernetes clusters with CoreOS hosts and TLS authentication.

### Assets

Download the CoreOS PXE image assets to `examples/assets/coreos`. These images are served to network boot machines by `bootcfg`.

    ./scripts/get-coreos alpha 983.0.0

**Note**: TLS assets are served to any machines which request them. This is unsuitable for production where machines and networks are untrusted. Read about our longer term security plans at [Distributed Trusted Computing](https://coreos.com/blog/coreos-trusted-computing.html).

Generate a root CA and Kubernetes TLS assets for components (`admin`, `apiserver`, `worker`).

    rm -rf examples/assets/tls
    # for Kubernetes on CNI metal0, i.e. rkt
    ./scripts/tls/k8s-certgen -d examples/assets/tls -s 172.15.0.21 -m IP.1=10.3.0.1,IP.2=172.15.0.21 -w IP.1=172.15.0.22,IP.2=172.15.0.23
    # for Kubernetes on docker0
    ./scripts/tls/k8s-certgen -d examples/assets/tls -s 172.17.0.21 -m IP.1=10.3.0.1,IP.2=172.17.0.21 -w IP.1=172.17.0.22,IP.2=172.17.0.23

See the [Cluster TLS OpenSSL Generation](https://coreos.com/kubernetes/docs/latest/openssl.html) document or [Kubernetes Step by Step](https://coreos.com/kubernetes/docs/latest/getting-started.html) for more details.

### Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your host. Use the provided kubeconfig's to access the Kubernetes cluster created on rkt `metal0` or `docker0`.

    cd /path/to/coreos-baremetal
    kubectl --kubeconfig=examples/assets/tls/kubeconfig get nodes

Get all pods.

    kubectl --kubeconfig=examples/assets/tls/kubeconfig get pods --all-namespaces

On my laptop, VMs download and network boot CoreOS in the first 45 seconds, the Kubernetes API becomes available after about 150 seconds, and add-on pods are scheduled after 3 minutes. On physical hosts and networks, OS and container image download times are a bit longer.

## Tectonic

Now sign up for [Tectonic Starter](https://tectonic.com/starter/) for free and deploy the [Tectonic Console](https://tectonic.com/enterprise/docs/latest/deployer/tectonic_console.html) with a few `kubectl` commands!

