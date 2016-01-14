
# Kubernetes

This example provisions a Kubernetes cluster with 1 master node, 1 worker node, and a dedicated etcd node. Each node uses a static IP address on the local network.

## Assets

Download the required CoreOS Beta image assets.

    ./scripts/get-coreos beta 877.1.0

Next, add or generate a root CA and Kubernetes TLS assets for each component.

### TLS Assets

Use the `generate-tls` script to generate throw-away TLS assets. The script will generate a root CA and `admin`, `apiserver`, and `worker` certificates in `assets/tls`.

    ./examples/kubernetes/scripts/generate-tls

Alternately, if you have existing Public Key Infrastructure, add your CA certificate, entity certificates, and entity private keys to `assets/tls`.

    * ca.pem
    * apiserver.pem
    * apiserver-key.pem
    * worker.pem
    * worker-key.pem
    * admin.pem
    * admin-key.pem

See the [Cluster TLS OpenSSL Generation](https://coreos.com/kubernetes/docs/latest/openssl.html) document or [Kubernetes Step by Step](https://coreos.com/kubernetes/docs/latest/getting-started.html) for more details.

Return the the general examples [README](../README).

## Usage

Install `kubectl` on your host and use the `examples/kubernetes/kubeconfig` file which references the top level `assets/tls`.

    cd /path/to/coreos-baremetal
    kubectl --kubeconfig=examples/kubernetes/kubeconfig get nodes

Watch pod events.

    kubectl --kubeconfig=examples/kubernetes/kubeconfig get pods --all-namespaces -w

Get all pods.

    kubectl --kubeconfig=examples/kubernetes/kubeconfig get pods --all-namespaces

On my laptop, it takes about 1 minute from boot until the Kubernetes API comes up. Then it takes another 1-2 minutes for all components including DNS to be pulled and started.