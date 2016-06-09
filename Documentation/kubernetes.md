
# Kubernetes

The Kubernetes examples provision a 3 node v1.2.4 Kubernetes cluster with one master, two workers, and TLS authentication. A 3 node etcd cluster is run on the hosts for Kubernetes and to coordinate CoreOS auto-updates (if installed to disk).

## Requirements

Ensure that you've gone through the [bootcfg with rkt](getting-started-rkt.md) guide and understand the basics. In particular, you should be able to:

* Use rkt or Docker to start `bootcfg`
* Create a network boot environment with `coreos/dnsmasq`
* Create the example libvirt client VMs

## Examples

The [examples](../examples) statically assign IP addresses to libvirt client VMs created by `scripts/libvirt`. VMs are setup on the `metal0` CNI bridge for rkt or the `docker0` bridge for Docker. You can use the same examples for real hardware, but you'll need to update the MAC/IP addresses.

* [k8s](../examples/groups/k8s) - iPXE boot a Kubernetes cluster (use rkt)
* [k8s-docker](../examples/groups/k8s-docker) - iPXE boot a Kubernetes cluster on `docker0` (use docker)
* [k8s-install](../examples/groups/k8s-install) - Install a Kubernetes cluster to disk (use rkt)

### Assets

Download the CoreOS image assets referenced in the target [profile](../examples/profiles).

    ./scripts/get-coreos alpha 1053.2.0

Generate a root CA and Kubernetes TLS assets for components (`admin`, `apiserver`, `worker`).

    rm -rf examples/assets/tls
    # for Kubernetes on CNI metal0, i.e. rkt
    ./scripts/tls/k8s-certgen -d examples/assets/tls -s 172.15.0.21 -m IP.1=10.3.0.1,IP.2=172.15.0.21 -w IP.1=172.15.0.22,IP.2=172.15.0.23
    # for Kubernetes on docker0
    ./scripts/tls/k8s-certgen -d examples/assets/tls -s 172.17.0.21 -m IP.1=10.3.0.1,IP.2=172.17.0.21 -w IP.1=172.17.0.22,IP.2=172.17.0.23

**Note**: TLS assets are served to any machines which request them. This is unsuitable for production where machines and networks are untrusted. Read about our longer term security plans at [Distributed Trusted Computing](https://coreos.com/blog/coreos-trusted-computing.html). See the [Cluster TLS OpenSSL Generation](https://coreos.com/kubernetes/docs/latest/openssl.html) document or [Kubernetes Step by Step](https://coreos.com/kubernetes/docs/latest/getting-started.html) for more details.

Optionally add your SSH public key to each machine group definition [as shown](../examples/README.md#ssh-keys).

## Containers

Use rkt or docker to start `bootcfg` with the desired example machine groups. Create a network boot environment with `coreos/dnsmasq` and create VMs with `scripts/libvirt` to power-on your machines. Client machines should boot and provision themselves.

Revisit [bootcfg with rkt](getting-started-rkt.md) or [bootcfg with Docker](getting-started-docker.md) for help.

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster created on rkt `metal0` or `docker0`.

    cd /path/to/coreos-baremetal
    kubectl --kubeconfig=examples/assets/tls/kubeconfig get nodes
    NAME          STATUS                     AGE
    172.15.0.21   Ready,SchedulingDisabled   6m
    172.15.0.22   Ready                      5m
    172.15.0.23   Ready                      6m

Get all pods.

    kubectl --kubeconfig=examples/assets/tls/kubeconfig get pods --all-namespaces
    NAMESPACE     NAME                                  READY     STATUS    RESTARTS   AGE
    kube-system   heapster-v1.0.2-808903792-yyyvw       2/2       Running   0          6m
    kube-system   kube-apiserver-172.15.0.21            1/1       Running   0          5m
    kube-system   kube-controller-manager-172.15.0.21   1/1       Running   0          5m
    kube-system   kube-dns-v11-uaajz                    4/4       Running   0          6m
    kube-system   kube-proxy-172.15.0.21                1/1       Running   0          6m
    kube-system   kube-proxy-172.15.0.22                1/1       Running   0          6m
    kube-system   kube-proxy-172.15.0.23                1/1       Running   0          6m
    kube-system   kube-scheduler-172.15.0.21            1/1       Running   0          6m

Machines should download and network boot CoreOS in about a minute. It can take 2-3 minutes for the Kubernetes API to become available and for add-on pods to be scheduled.

## Tectonic

Now sign up for [Tectonic Starter](https://tectonic.com/starter/) for free and deploy the [Tectonic Console](https://tectonic.com/enterprise/docs/latest/deployer/tectonic_console.html) with a few `kubectl` commands!

<img src='img/tectonic-console.png' class="img-center" alt="Tectonic Console"/>

