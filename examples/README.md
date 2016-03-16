
# Examples

These examples show declarative configurations for network booting libvirt VMs into CoreOS clusters (Kubernetes, etcd) using `bootcfg`.

| Name       | Description | CoreOS Version | FS | Reference | 
|------------|-------------|----------------|----|-----------|
| pxe | CoreOS via iPXE | alpha/962.0.0 | RAM | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| grub | CoreOS via GRUB2 Netboot | beta/899.6.0 | RAM | NA |
| pxe-disk | CoreOS via iPXE, with a root filesystem | alpha/962.0.0 | Disk | [reference](https://coreos.com/os/docs/latest/booting-with-ipxe.html) |
| coreos-install | 2-stage Ignition: Install CoreOS, provision etcd cluster | alpha/962.0.0 | Disk | [reference](https://coreos.com/os/docs/latest/installing-to-disk.html) |
| etcd-rkt, etcd-docker | Cluster with 3 etcd nodes, 2 proxies | beta/899.6.0 | RAM | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) |
| k8s-rkt, k8s-docker | Kubernetes cluster with 1 master, 1 worker, 1 dedicated etcd node, TLS-authentication | beta/899.6.0 | RAM | [reference](https://github.com/coreos/coreos-kubernetes) |

## Experimental

These CoreOS clusters are experimental and have **NOT** been hardened for production yet. They demonstrate Ignition and cloud-init provisioning of higher order clusters.

## Getting Started

Get started running the `bootcfg` on your Linux machine to boot clusters of libvirt PXE VMs.

* [Getting Started with rkt](../Documentation/getting-started-rkt.md)
* [Getting Started with Docker](../Documentation/getting-started-docker.md)

## SSH Keys

Most example profiles configure machines with a `core` user and `ssh_authorized_keys`. Add your own key(s) as machine metadata.

    ---
    api_version: v1alpha1
    groups:
      - name: default
        profile: pxe
        metadata:
          ssh_authorized_keys:
            - "ssh-rsa pub-key-goes-here"

## Kubernetes

The Kubernetes cluster examples create a TLS-authenticated Kubernetes cluster with 1 master node, 1 worker node, and 1 etcd node, running without a disk.

You'll need to download the CoreOS Beta image, which ships with the kubelet, and generate TLS assets.

### TLS Assets

**Note**: TLS assets are served to any machines which request them. This is unsuitable for production where machines and networks are untrusted. Read about our longer term security plans at [Distributed Trusted Computing](https://coreos.com/blog/coreos-trusted-computing.html).

Generate a root CA and Kubernetes TLS assets for each component (`admin`, `apiserver`, `worker`).

    cd coreos-baremetal
    # for Kubernetes on CNI metal0, i.e. rkt
    ./scripts/tls/gen-rkt-k8s-secrets
    # for Kubernetes on docker0
    ./scripts/tls/gen-docker-k8s-secrets

Alternately, you can add your own CA certificate, entity certificates, and entity private keys to `assets/tls`.

    * ca.pem
    * apiserver.pem
    * apiserver-key.pem
    * worker.pem
    * worker-key.pem
    * admin.pem
    * admin-key.pem

See the [Cluster TLS OpenSSL Generation](https://coreos.com/kubernetes/docs/latest/openssl.html) document or [Kubernetes Step by Step](https://coreos.com/kubernetes/docs/latest/getting-started.html) for more details.

### Verify

Install the `kubectl` CLI on your host. Use the provided kubeconfig's to access the Kubernetes cluster created on rkt `metal0` or `docker0`.

    cd /path/to/coreos-baremetal
    # for kubernetes on CNI metal0, i.e. rkt
    kubectl --kubeconfig=examples/kubecfg-rkt get nodes
    # for kubernetes on docker0
    kubectl --kubeconfig=examples/kubecfg-docker get nodes

Get all pods.

    kubectl --kubeconfig=examples/kubecfg-rkt get pods --all-namespaces

On my laptop, it takes about 1 minute from boot until the Kubernetes API comes up. Then it takes another 1-2 minutes for all components including DNS to be pulled and started.

