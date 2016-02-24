
# Examples

These examples show declarative configurations for network booting libvirt VMs into CoreOS clusters (Kubernetes, etcd) using `bootcfg`.

| Name       | Description |  Reference     | CoreOS Version |
|------------|-------------|----------------|----------------|
| etcd | Cluster with 3 etcd nodes, 2 proxies | [reference](https://coreos.com/os/docs/latest/cluster-architectures.html) | beta/899.6.0 |
| Kubernetes | Kubernetes cluster with 1 master, 1 worker, 1 dedicated etcd node | [reference](https://github.com/coreos/coreos-kubernetes) | beta/899.6.0 |
| Disk install w etcd | 2-stage Ignition: Install CoreOS, provision etcd cluster | [reference](https://coreos.com/os/docs/latest/installing-to-disk.html) | alpha/962.0.0,935.0.0 |
| alpha-pxe | PXE CoreOS alpha node, uses configured SSH keys | | alpha/962.0.0 |

## Experimental

These CoreOS clusters are experimental and have **NOT** been hardened for production yet. They demonstrate Ignition (initrd) and cloud-init provisioning of higher order clusters.

## Getting Started

Get started running the `bootcfg` on your Linux machine to boot clusters of libvirt PXE VMs.

* [Getting Started with rkt](../Documentation/getting-started-rkt.md)
* [Getting Started with Docker](../Documentation/getting-started-docker.md)

## Physical Hardware

Run `bootcfg` to boot and configure physical machines (for testing). Update the network values in the `*.yaml` config to match your hardware and network. Generate TLS assets if required for the example (e.g. Kubernetes).

Continue to the [Physical Hardware Guide](../Documentation/physical-hardware.md) for details.

## Examples

See the Getting Started with [rkt](getting-started-rkt.md) or [Docker](getting-started-docker.md) for a walk-through.

### rkt

etcd cluster with 3 nodes on `metal0`, other nodes act as proxies.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/assets --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/data --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg -- -address=0.0.0.0:8080 -log-level=debug -config /data/etcd-rkt.yaml

Kubernetes cluster with one master, one worker, and one dedicated etcd on `metal0`.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/assets --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/data --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg -- -address=0.0.0.0:8080 -log-level=debug -config /data/k8s-rkt.yaml

### Docker

etcd cluster with 3 nodes on `docker0`, other nodes act as proxies.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/data:Z -v $PWD/assets:/assets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -config /data/etcd-docker.yaml

Kubernetes cluster with one master, one worker, and one dedicated etcd on `docker0`.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/data:Z -v $PWD/assets:/assets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -config /data/k8s-docker.yaml

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

