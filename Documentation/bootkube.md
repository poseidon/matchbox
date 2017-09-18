# Kubernetes

The Kubernetes example provisions a 3 node Kubernetes v1.7.5 cluster. [bootkube](https://github.com/kubernetes-incubator/bootkube) is run once on a controller node to bootstrap Kubernetes control plane components as pods before exiting. An etcd3 cluster across controllers is used to back Kubernetes.

## Requirements

Ensure that you've gone through the [matchbox with rkt](getting-started-rkt.md) or [matchbox with docker](getting-started-docker.md) guide and understand the basics. In particular, you should be able to:

* Use rkt or Docker to start `matchbox`
* Create a network boot environment with `coreos/dnsmasq`
* Create the example libvirt client VMs
* `/etc/hosts` entries for `node[1-3].example.com`

Install [bootkube](https://github.com/kubernetes-incubator/bootkube/releases) v0.6.2 and add it on your $PATH.

```sh
$ bootkube version
Version: v0.6.2
```

## Examples

The [examples](../examples) statically assign IP addresses to libvirt client VMs created by `scripts/libvirt`. The examples can be used for physical machines if you update the MAC addresses. See [network setup](network-setup.md) and [deployment](deployment.md).

* [bootkube](../examples/groups/bootkube) - iPXE boot a self-hosted Kubernetes cluster
* [bootkube-install](../examples/groups/bootkube-install) - Install a self-hosted Kubernetes cluster

## Assets

Download the CoreOS Container Linux image assets referenced in the target [profile](../examples/profiles).

```sh
$ ./scripts/get-coreos stable 1465.7.0 ./examples/assets
```

Add your SSH public key to each machine group definition [as shown](../examples/README.md#ssh-keys).

```json
{
    "profile": "bootkube-worker",
    "metadata": {
        "ssh_authorized_keys": ["ssh-rsa pub-key-goes-here"]
    }
}
```

Use the `bootkube` tool to render Kubernetes manifests and credentials into an `--asset-dir`. Later, `bootkube` will schedule these manifests during bootstrapping and the credentials will be used to access your cluster.

```sh
bootkube render --asset-dir=assets --api-servers=https://node1.example.com:443 --api-server-alt-names=DNS=node1.example.com --etcd-servers=https://node1.example.com:2379
```

## Containers

Use rkt or docker to start `matchbox` and mount the desired example resources. Create a network boot environment and power-on your machines. Revisit [matchbox with rkt](getting-started-rkt.md) or [matchbox with Docker](getting-started-docker.md) for help.

Client machines should boot and provision themselves. Local client VMs should network boot Container Linux and become available via SSH in about 1 minute. If you chose `bootkube-install`, notice that machines install Container Linux and then reboot (in libvirt, you must hit "power" again). Time to network boot and provision physical hardware depends on a number of factors (POST duration, boot device iteration, network speed, etc.).

## bootkube

We're ready to use bootkube to create a temporary control plane and bootstrap a self-hosted Kubernetes cluster.

Secure copy the etcd TLS assets to `/etc/ssl/etcd/*` on **every controller** node.

```sh
for node in 'node1'; do
    scp -r assets/tls/etcd-* assets/tls/etcd core@$node.example.com:/home/core/
    ssh core@$node.example.com 'sudo mkdir -p /etc/ssl/etcd && sudo mv etcd-* etcd /etc/ssl/etcd/ && sudo chown -R etcd:etcd /etc/ssl/etcd && sudo chmod -R 500 /etc/ssl/etcd/'
done
```

Secure copy the `kubeconfig` to `/etc/kubernetes/kubeconfig` on **every node** to path activate the `kubelet.service`.

```sh
for node in 'node1' 'node2' 'node3'; do
    scp assets/auth/kubeconfig core@$node.example.com:/home/core/kubeconfig
    ssh core@$node.example.com 'sudo mv kubeconfig /etc/kubernetes/kubeconfig'
done
```

Secure copy the `bootkube` generated assets to **any controller** node and run `bootkube-start` (takes ~10 minutes).

```sh
scp -r assets core@node1.example.com:/home/core
ssh core@node1.example.com 'sudo mv assets /opt/bootkube/assets && sudo systemctl start bootkube'
```

Watch the Kubernetes control plane bootstrapping with the bootkube temporary api-server. You will see quite a bit of output.

```sh
$ ssh core@node1.example.com 'journalctl -f -u bootkube'
[  299.241291] bootkube[5]:     Pod Status:     kube-api-checkpoint     Running
[  299.241618] bootkube[5]:     Pod Status:          kube-apiserver     Running
[  299.241804] bootkube[5]:     Pod Status:          kube-scheduler     Running
[  299.241993] bootkube[5]:     Pod Status: kube-controller-manager     Running
[  299.311743] bootkube[5]: All self-hosted control plane components successfully started
```

[Verify](#verify) the Kubernetes cluster is accessible once complete. Then install **important** cluster [addons](cluster-addons.md). You may cleanup the `bootkube` assets on the node, but you should keep the copy on your laptop. It contains a `kubeconfig` used to access the cluster.

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster. Verify that the cluster is accessible and that the apiserver, scheduler, and controller-manager are running as pods.

```sh
$ export KUBECONFIG=assets/auth/kubeconfig
$ kubectl get nodes
NAME                STATUS    AGE       VERSION
node1.example.com   Ready     11m       v1.7.5+coreos.0
node2.example.com   Ready     11m       v1.7.5+coreos.0
node3.example.com   Ready     11m       v1.7.5+coreos.0

$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                       READY     STATUS    RESTARTS   AGE
kube-system   kube-apiserver-zd1k3                       1/1       Running   0          7m
kube-system   kube-controller-manager-762207937-2ztxb    1/1       Running   0          7m
kube-system   kube-controller-manager-762207937-vf6bk    1/1       Running   1          7m
kube-system   kube-dns-2431531914-qc752                  3/3       Running   0          7m
kube-system   kube-flannel-180mz                         2/2       Running   1          7m
kube-system   kube-flannel-jjr0x                         2/2       Running   0          7m
kube-system   kube-flannel-mlr9w                         2/2       Running   0          7m
kube-system   kube-proxy-0jlq7                           1/1       Running   0          7m
kube-system   kube-proxy-k4mjl                           1/1       Running   0          7m
kube-system   kube-proxy-l4xrd                           1/1       Running   0          7m
kube-system   kube-scheduler-1873228005-5d2mk            1/1       Running   0          7m
kube-system   kube-scheduler-1873228005-s4w27            1/1       Running   0          7m
kube-system   pod-checkpointer-hb960                     1/1       Running   0          7m
kube-system   pod-checkpointer-hb960-node1.example.com   1/1       Running   0          6m
```

## Addons

Install **important** cluster [addons](cluster-addons.md).

## Going further

[Learn](bootkube-upgrades.md) to upgrade a self-hosted Kubernetes cluster.
