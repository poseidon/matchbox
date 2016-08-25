
# Kubernetes

The Kubernetes example provisions a 3 node Kubernetes v1.3.4 cluster with one controller, two workers, and TLS authentication. An etcd cluster backs Kubernetes and coordinates CoreOS auto-updates (enabled for disk installs).

## Requirements

Ensure that you've gone through the [bootcfg with rkt](getting-started-rkt.md) or [bootcfg with docker](getting-started-docker.md) guide and understand the basics. In particular, you should be able to:

* Use rkt or Docker to start `bootcfg`
* Create a network boot environment with `coreos/dnsmasq`
* Create the example libvirt client VMs

## Examples

The [examples](../examples) statically assign IP addresses to libvirt client VMs created by `scripts/libvirt`. VMs are setup on the `metal0` CNI bridge for rkt or the `docker0` bridge for Docker. The examples can be used for physical machines if you update the MAC/IP addresses. See [network setup](network-setup.md) and [deployment](deployment.md).

* [k8s](../examples/groups/k8s) - iPXE boot a Kubernetes cluster
* [k8s-install](../examples/groups/k8s-install) - Install a Kubernetes cluster to disk
* [Lab examples](https://github.com/dghubble/metal) - Lab hardware examples

### Assets

Download the CoreOS image assets referenced in the target [profile](../examples/profiles).

    ./scripts/get-coreos alpha 1153.0.0 ./examples/assets

Add your SSH public key to each machine group definition [as shown](../examples/README.md#ssh-keys).

Generate a root CA and Kubernetes TLS assets for components (`admin`, `apiserver`, `worker`).

    rm -rf examples/assets/tls
    # for Kubernetes on CNI metal0 (for rkt)
    ./scripts/tls/k8s-certgen -d examples/assets/tls -s 172.15.0.21 -m IP.1=10.3.0.1,IP.2=172.15.0.21,DNS.1=node1.example.com -w DNS.1=node2.example.com,DNS.2=node3.example.com
    # for Kubernetes on docker0 (for docker)
    ./scripts/tls/k8s-certgen -d examples/assets/tls -s 172.17.0.21 -m IP.1=10.3.0.1,IP.2=172.17.0.21,DNS.1=node1.example.com -w DNS.1=node2.example.com,DNS.2=node3.example.com

**Note**: TLS assets are served to any machines which request them, which requires a trusted network. Alternately, provisioning may be tweaked to require TLS assets be securely copied to each host. Read about our longer term security plans at [Distributed Trusted Computing](https://coreos.com/blog/coreos-trusted-computing.html).

## Containers

Use rkt or docker to start `bootcfg` and mount the desired example resources. Create a network boot environment and power-on your machines. Revisit [bootcfg with rkt](getting-started-rkt.md) or [bootcfg with Docker](getting-started-docker.md) for help.

Client machines should boot and provision themselves. Local client VMs should network boot CoreOS in about a 1 minute and the Kubernetes API should be available after 3-4 minutes (each node downloads a ~160MB Hyperkube). If you chose `k8s-install`, notice that machines install CoreOS and then reboot (in libvirt, you must hit "power" again). Time to network boot and provision Kubernetes clusters on physical hardware depends on a number of factors (POST duration, boot device iteration, network speed, etc.).

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster created on rkt `metal0` or `docker0`.

    $ cd /path/to/coreos-baremetal
    $ kubectl --kubeconfig=examples/assets/tls/kubeconfig get nodes
    NAME                STATUS    AGE
    node1.example.com   Ready     43s
    node2.example.com   Ready     38s
    node3.example.com   Ready     37s

Get all pods.

    $ kubectl --kubeconfig=examples/assets/tls/kubeconfig get pods --all-namespaces
    NAMESPACE     NAME                                  READY     STATUS    RESTARTS   AGE
    kube-system   heapster-v1.1.0-3647315203-tes6g      2/2       Running   0          14m
    kube-system   kube-apiserver-172.15.0.21            1/1       Running   0          14m
    kube-system   kube-controller-manager-172.15.0.21   1/1       Running   0          14m
    kube-system   kube-dns-v15-nfbz4                    3/3       Running   0          14m
    kube-system   kube-proxy-172.15.0.21                1/1       Running   0          14m
    kube-system   kube-proxy-172.15.0.22                1/1       Running   0          14m
    kube-system   kube-proxy-172.15.0.23                1/1       Running   0          14m
    kube-system   kube-scheduler-172.15.0.21            1/1       Running   0          13m
    kube-system   kubernetes-dashboard-v1.1.0-m1gyy     1/1       Running   0          14m

## Kubernetes Dashboard

Access the Kubernetes Dashboard with `kubeconfig` credentials by port forwarding to the dashboard pod.

    $ kubectl --kubeconfig=examples/assets/tls/kubeconfig port-forward kubernetes-dashboard-v1.1.0-SOME-ID 9090 --namespace=kube-system
    Forwarding from 127.0.0.1:9090 -> 9090

Then visit [http://127.0.0.1:9090](http://127.0.0.1:9090/).

<img src='img/kubernetes-dashboard.png' class="img-center" alt="Kubernetes Dashboard"/>

## Tectonic

Sign up for [Tectonic Starter](https://tectonic.com/starter/) for free and deploy the [Tectonic Console](https://tectonic.com/enterprise/docs/latest/deployer/tectonic_console.html) with a few `kubectl` commands!

<img src='img/tectonic-console.png' class="img-center" alt="Tectonic Console"/>

