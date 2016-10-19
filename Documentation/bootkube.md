
# Self-Hosted Kubernetes

The self-hosted Kubernetes example provisions a 3 node "self-hosted" Kubernetes v1.4.1 cluster. On-host kubelets wait for an apiserver to become reachable, then yield to kubelet pods scheduled via daemonset. [bootkube](https://github.com/kubernetes-incubator/bootkube) is run on any controller to bootstrap a temporary apiserver which schedules control plane components as pods before exiting. An etcd cluster backs Kubernetes and coordinates CoreOS auto-updates (enabled for disk installs).

## Requirements

Ensure that you've gone through the [bootcfg with rkt](getting-started-rkt.md) or [bootcfg with docker](getting-started-docker.md) guide and understand the basics. In particular, you should be able to:

* Use rkt or Docker to start `bootcfg`
* Create a network boot environment with `coreos/dnsmasq`
* Create the example libvirt client VMs

Build and install the [fork of bootkube](https://github.com/dghubble/bootkube), which supports DNS names.

## Examples

The [examples](../examples) statically assign IP addresses to libvirt client VMs created by `scripts/libvirt`. The examples can be used for physical machines if you update the MAC addresses. See [network setup](network-setup.md) and [deployment](deployment.md).

* [bootkube](../examples/groups/bootkube) - iPXE boot a self-hosted Kubernetes cluster
* [bootkube-install](../examples/groups/bootkube-install) - Install a self-hosted Kubernetes cluster

### Assets

Download the CoreOS image assets referenced in the target [profile](../examples/profiles).

    ./scripts/get-coreos beta 1185.1.0 ./examples/assets

Add your SSH public key to each machine group definition [as shown](../examples/README.md#ssh-keys).

    {
        "profile": "bootkube-worker",
        "metadata": {
            "ssh_authorized_keys": ["ssh-rsa pub-key-goes-here"]
        }
    }

Use the `bootkube` tool to render Kubernetes manifests and credentials into an `--asset-dir`. Later, `bootkube` will schedule these manifests during bootstrapping and the credentials will be used to access your cluster.

    # If running with docker, use 172.17.0.21 instead of 172.15.0.21
    bootkube render --asset-dir=assets --api-servers=https://172.15.0.21:443 --etcd-servers=http://node1.example.com:2379 --api-server-alt-names=DNS=node1.example.com,IP=172.15.0.21

## Containers

Use rkt or docker to start `bootcfg` and mount the desired example resources. Create a network boot environment and power-on your machines. Revisit [bootcfg with rkt](getting-started-rkt.md) or [bootcfg with Docker](getting-started-docker.md) for help.

Client machines should boot and provision themselves. Local client VMs should network boot CoreOS and become available via SSH in about 1 minute. If you chose `bootkube-install`, notice that machines install CoreOS and then reboot (in libvirt, you must hit "power" again). Time to network boot and provision physical hardware depends on a number of factors (POST duration, boot device iteration, network speed, etc.).

## bootkube

We're ready to use [bootkube](https://github.com/kubernetes-incubator/bootkube) to create a temporary control plane and bootstrap a self-hosted Kubernetes cluster.

Secure copy the `kubeconfig` to `/etc/kubernetes/kubeconfig` on **every** node (i.e. 172.15.0.21-23 for metal0 or 172.17.0.21-23 for docker0).

    for node in '172.15.0.21' '172.15.0.22' '172.15.0.23'; do
        scp assets/auth/kubeconfig core@$node:/home/core/kubeconfig
        ssh core@$node 'sudo mv kubeconfig /etc/kubernetes/kubeconfig'
    done

Secure copy the `bootkube` generated assets to any controller node and run `bootkube-start`.

    scp -r assets core@172.15.0.21:/home/core/assets
    ssh core@172.15.0.21 'sudo ./bootkube-start'

Watch the temporary control plane logs until the scheduled kubelet takes over in place of the on-host kubelet.

    [  299.241291] bootkube[5]:     Pod Status:     kube-api-checkpoint     Running
    [  299.241618] bootkube[5]:     Pod Status:          kube-apiserver     Running
    [  299.241804] bootkube[5]:     Pod Status:          kube-scheduler     Running
    [  299.241993] bootkube[5]:     Pod Status: kube-controller-manager     Running
    [  299.311743] bootkube[5]: All self-hosted control plane components successfully started

You may cleanup the `bootkube` assets on the node, but you should keep the copy on your laptop. It contains a `kubeconfig` and may need to be re-used if the last apiserver were to fail and bootstrapping were needed.

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster. Verify that the cluster is accessible and that the kubelet, apiserver, scheduler, and controller-manager are running as pods.

    $ KUBECONFIG=assets/auth/kubeconfig
    $ kubectl get nodes
    NAME                STATUS    AGE
    node1.example.com   Ready     3m
    node2.example.com   Ready     3m
    node3.example.com   Ready     3m

    $ kubectl get pods --all-namespaces
    NAMESPACE     NAME                                       READY     STATUS    RESTARTS   AGE
    kube-system   kube-api-checkpoint-node1.example.com      1/1       Running   0          4m
    kube-system   kube-apiserver-iffsz                       2/2       Running   0          5m
    kube-system   kube-controller-manager-1148212084-1zx9g   1/1       Running   0          6m
    kube-system   kube-dns-v20-3531996453-r18ht              3/3       Running   0          5m
    kube-system   kube-proxy-36jj8                           1/1       Running   0          5m
    kube-system   kube-proxy-fdt2t                           1/1       Running   0          6m
    kube-system   kube-proxy-sttgn                           1/1       Running   0          5m
    kube-system   kube-scheduler-1921762579-z6jn6            1/1       Running   0          6m
    kube-system   kubelet-1ibsf                              1/1       Running   0          6m
    kube-system   kubelet-65h6j                              1/1       Running   0          5m
    kube-system   kubelet-d1qql                              1/1       Running   0          5m

Try deleting pods to see that the cluster is resilient to failures and machine restarts (CoreOS auto-updates).

## Going Further

[Learn](bootkube-upgrades.md) to upgrade a self-hosted Kubernetes cluster.