
# Self-Hosted Kubernetes

The self-hosted Kubernetes example provisions a 3 node "self-hosted" Kubernetes v1.4.7 cluster. On-host kubelets wait for an apiserver to become reachable, then yield to kubelet pods scheduled via daemonset. [bootkube](https://github.com/kubernetes-incubator/bootkube) is run on any controller to bootstrap a temporary apiserver which schedules control plane components as pods before exiting. An etcd cluster backs Kubernetes and coordinates CoreOS auto-updates (enabled for disk installs).

## Requirements

Ensure that you've gone through the [bootcfg with rkt](getting-started-rkt.md) or [bootcfg with docker](getting-started-docker.md) guide and understand the basics. In particular, you should be able to:

* Use rkt or Docker to start `bootcfg`
* Create a network boot environment with `coreos/dnsmasq`
* Create the example libvirt client VMs
* `/etc/hosts` entries for `node[1-3].example.com` (or pass custom names to `k8s-certgen`)

Build and install the [fork of bootkube](https://github.com/dghubble/bootkube), which supports DNS names.

    $ bootkube version
    Version: 4a4ae6e78a59258b528f7de57db8d5cf5786a8a5

## Examples

The [examples](../examples) statically assign IP addresses to libvirt client VMs created by `scripts/libvirt`. The examples can be used for physical machines if you update the MAC addresses. See [network setup](network-setup.md) and [deployment](deployment.md).

* [bootkube](../examples/groups/bootkube) - iPXE boot a self-hosted Kubernetes cluster
* [bootkube-install](../examples/groups/bootkube-install) - Install a self-hosted Kubernetes cluster

### Assets

Download the CoreOS image assets referenced in the target [profile](../examples/profiles).

    ./scripts/get-coreos stable 1185.3.0 ./examples/assets

Add your SSH public key to each machine group definition [as shown](../examples/README.md#ssh-keys).

    {
        "profile": "bootkube-worker",
        "metadata": {
            "ssh_authorized_keys": ["ssh-rsa pub-key-goes-here"]
        }
    }

Use the `bootkube` tool to render Kubernetes manifests and credentials into an `--asset-dir`. Later, `bootkube` will schedule these manifests during bootstrapping and the credentials will be used to access your cluster.

    bootkube render --asset-dir=assets --api-servers=https://node1.example.com:443 --api-server-alt-names=DNS=node1.example.com

## Containers

Use rkt or docker to start `bootcfg` and mount the desired example resources. Create a network boot environment and power-on your machines. Revisit [bootcfg with rkt](getting-started-rkt.md) or [bootcfg with Docker](getting-started-docker.md) for help.

Client machines should boot and provision themselves. Local client VMs should network boot CoreOS and become available via SSH in about 1 minute. If you chose `bootkube-install`, notice that machines install CoreOS and then reboot (in libvirt, you must hit "power" again). Time to network boot and provision physical hardware depends on a number of factors (POST duration, boot device iteration, network speed, etc.).

## bootkube

We're ready to use bootkube to create a temporary control plane and bootstrap a self-hosted Kubernetes cluster.

Secure copy the `kubeconfig` to `/etc/kubernetes/kubeconfig` on **every** node which will path activate the `kubelet.service`.

    for node in 'node1' 'node2' 'node3'; do
        scp assets/auth/kubeconfig core@$node.example.com:/home/core/kubeconfig
        ssh core@$node.example.com 'sudo mv kubeconfig /etc/kubernetes/kubeconfig'
    done

Secure copy the `bootkube` generated assets to any controller node and run `bootkube-start`.

    scp -r assets core@node1.example.com:/home/core/assets
    ssh core@node1.example.com 'sudo systemctl start bootkube'

Optionally watch the Kubernetes control plane bootstrapping with the bootkube temporary api-server. You will see quite a bit of output.

    $ ssh core@node1.example.com 'journalctl -f -u bootkube'
    [  299.241291] bootkube[5]:     Pod Status:     kube-api-checkpoint     Running
    [  299.241618] bootkube[5]:     Pod Status:          kube-apiserver     Running
    [  299.241804] bootkube[5]:     Pod Status:          kube-scheduler     Running
    [  299.241993] bootkube[5]:     Pod Status: kube-controller-manager     Running
    [  299.311743] bootkube[5]: All self-hosted control plane components successfully started

You may cleanup the `bootkube` assets on the node, but you should keep the copy on your laptop. It contains a `kubeconfig` used to access the cluster.

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
    kube-system   checkpoint-installer-cpjrm                 1/1       Running   0          2m
    kube-system   kube-apiserver-rvjes                       1/1       Running   0          2m
    kube-system   kube-controller-manager-3900529476-5n9xb   1/1       Running   0          4m
    kube-system   kube-controller-manager-3900529476-rq6p8   1/1       Running   0          4m
    kube-system   kube-dns-4101612645-oeu5g                  4/4       Running   0          4m
    kube-system   kube-proxy-f5kb8                           1/1       Running   0          2m
    kube-system   kube-proxy-jkg4z                           1/1       Running   0          2m
    kube-system   kube-proxy-rrmuv                           1/1       Running   0          2m
    kube-system   kube-scheduler-1084603659-qaqbk            1/1       Running   0          3m
    kube-system   kube-scheduler-1084603659-rztvw            1/1       Running   0          4m
    kube-system   pod-checkpointer-node1.example.com         1/1       Running   0          1m

Try deleting pods to see that the cluster is resilient to failures and machine restarts (CoreOS auto-updates).

## Going Further

[Learn](bootkube-upgrades.md) to upgrade a self-hosted Kubernetes cluster.