
# Self-Hosted Kubernetes

The self-hosted Kubernetes example provisions a 3 node Kubernetes v1.3.0 cluster with etcd, flannel, and a special "runonce" host Kublet. The CoreOS [bootkube](https://github.com/coreos/bootkube) tool is used to bootstrap kubelet, apiserver, scheduler, and controller-manager as pods, which can be managed via kubectl. `bootkube start` is run on any controller (master) to create a temporary control-plane and start Kubernetes components initially. An etcd cluster backs Kubernetes and coordinates CoreOS auto-updates (enabled for disk installs).

## Experimental

Self-hosted Kubernetes is under very active development by CoreOS.

## Requirements

Ensure that you've gone through the [bootcfg with rkt](getting-started-rkt.md) guide and understand the basics. In particular, you should be able to:

* Use rkt to start `bootcfg`
* Create a network boot environment with `coreos/dnsmasq`
* Create the example libvirt client VMs

Build and install [bootkube](https://github.com/coreos/bootkube/releases) v0.1.2.

## Examples

The [examples](../examples) statically assign IP addresses to libvirt client VMs created by `scripts/libvirt`. The examples can be used for physical machines if you update the MAC/IP addresses. See [network setup](network-setup.md) and [deployment](deployment.md).

* [bootkube](../examples/groups/bootkube) - iPXE boot a bootkube-ready cluster (use rkt)
* [bootkube-install](../examples/groups/bootkube-install) - Install a bootkube-ready cluster (use rkt)

### Assets

Download the CoreOS image assets referenced in the target [profile](../examples/profiles).

    ./scripts/get-coreos alpha 1109.1.0 ./examples/assets

Add your SSH public key to each machine group definition [as shown](../examples/README.md#ssh-keys).

    {
        "profile": "bootkube-worker",
        "metadata": {
            "ssh_authorized_keys": ["ssh-rsa pub-key-goes-here"]
        }
    }

Use the `bootkube` tool to render Kubernetes manifests and credentials into an `--asset-dir`. Later, `bootkube` will schedule these manifests during bootstrapping and the credentials will be used to access your cluster.

    bootkube render --asset-dir=assets --api-servers=https://172.15.0.21:443 --etcd-servers=http://172.15.0.21:2379 --api-server-alt-names=IP=172.15.0.21

## Containers

Run the latest `bootcfg` ACI with rkt and the `bootkube` example (or `bootkube-install`).

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/bootkube quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -log-level=debug

Create a network boot environment and power-on your machines. Revisit [bootcfg with rkt](getting-started-rkt.md) for help.

Client machines should boot and provision themselves. Local client VMs should network boot CoreOS and become available via SSH in about 1 minute. If you chose `bootkube-install`, notice that machines install CoreOS and then reboot (in libvirt, you must hit "power" again). Time to network boot and provision physical hardware depends on a number of factors (POST duration, boot device iteration, network speed, etc.).

## bootkube

We're ready to use [bootkube](https://github.com/coreos/bootkube) to create a temporary control plane and bootstrap a self-hosted Kubernetes cluster.

Secure copy the `kubeconfig` to `/etc/kuberentes/kubeconfig` on **every** node (i.e. repeat for 172.15.0.22, 172.15.0.23).

    scp assets/auth/kubeconfig core@172.15.0.21:/home/core/kubeconfig
    ssh core@172.15.0.21
    sudo mv kubeconfig /etc/kubernetes/kubeconfig

Secure copy the `bootkube` generated assets to any one of the master nodes.

    scp -r assets core@172.15.0.21:/home/core/assets

SSH to the chosen master node and bootstrap the cluster with `bootkube-start`.

    ssh core@172.15.0.21 'sudo ./bootkube-start'

Watch the temporary control plane logs until the scheduled kubelet takes over in place of the runonce host kubelet.

    I0425 12:38:23.746330   29538 status.go:87] Pod status kubelet: Running
    I0425 12:38:23.746361   29538 status.go:87] Pod status kube-apiserver: Running
    I0425 12:38:23.746370   29538 status.go:87] Pod status kube-scheduler: Running
    I0425 12:38:23.746378   29538 status.go:87] Pod status kube-controller-manager: Running

You may cleanup the `bootkube` assets on the node, but you should keep the copy on your laptop. They contain a `kubeconfig` and may need to be re-used if the last apiserver were to fail and bootstrapping were needed.

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster. Verify that the cluster is accessible and that the kubelet, apiserver, scheduler, and controller-manager are running as pods.

    $ kubectl --kubeconfig=assets/auth/kubeconfig get nodes
    NAME          STATUS    AGE
    172.15.0.21   Ready     3m
    172.15.0.22   Ready     3m
    172.15.0.23   Ready     3m

    $ kubectl --kubeconfig=assets/auth/kubeconfig get pods --all-namespaces
    kube-system   kube-api-checkpoint-172.15.0.21            1/1       Running   0          2m
    kube-system   kube-apiserver-wq4mh                       2/2       Running   0          2m
    kube-system   kube-controller-manager-2834499578-y9cnl   1/1       Running   0          2m
    kube-system   kube-dns-v11-2259792283-5tpld              4/4       Running   0          2m
    kube-system   kube-proxy-8zr1b                           1/1       Running   0          2m
    kube-system   kube-proxy-i9cgw                           1/1       Running   0          2m
    kube-system   kube-proxy-n6qg3                           1/1       Running   0          2m
    kube-system   kube-scheduler-4136156790-v9892            1/1       Running   0          2m
    kube-system   kubelet-9wilx                              1/1       Running   0          2m
    kube-system   kubelet-a6mmj                              1/1       Running   0          2m
    kube-system   kubelet-eomnb                              1/1       Running   0          2m

Try deleting pods to see that the cluster is resilient to failures and machine restarts (CoreOS auto-updates).
