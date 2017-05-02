# Self-hosted Kubernetes

The self-hosted Kubernetes example shows how to use matchbox to network boot and provision a 3 node "self-hosted" Kubernetes v1.6.2 cluster. [bootkube](https://github.com/kubernetes-incubator/bootkube) is run once on a controller node to bootstrap Kubernetes control plane components as pods before exiting. An etcd3 cluster across controllers is used to back Kubernetes and coordinate Container Linux auto-updates (enabled for disk installs).

## Requirements

Follow the getting started [tutorial](../../../Documentation/getting-started.md) to learn about matchbox and set up an environment that meets the requirements:

* Matchbox v0.6+ [installation](../../../Documentation/deployment.md) with gRPC API enabled
* Matchbox provider credentials `client.crt`, `client.key`, and `ca.crt`
* PXE [network boot](../../../Documentation/network-setup.md) environment
* Terraform v0.9+ and [terraform-provider-matchbox](https://github.com/coreos/terraform-provider-matchbox) installed locally on your system
* 3 machines with known DNS names and MAC addresses

If you prefer to provision QEMU/KVM VMs on your local Linux machine, set up the matchbox [development environment](../../../Documentation/getting-started-rkt.md).

```sh
sudo ./scripts/devnet create
```

## Usage

Clone the [matchbox](https://github.com/coreos/matchbox) project and take a look at the cluster examples.

```sh
$ git clone https://github.com/coreos/matchbox.git
$ cd matchbox/examples/terraform/bootkube-install
```

Copy the `terraform.tfvars.example` file to `terraform.tfvars`. Ensure `provider.tf` references your matchbox credentials.

```hcl
matchbox_http_endpoint = "http://matchbox.example.com:8080"
matchbox_rpc_endpoint = "matchbox.example.com:8081"
ssh_authorized_key = "ADD ME"
```

Configs in `bootkube-install` configure the matchbox provider, define profiles (e.g. `cached-container-linux-install`, `bootkube-controller`, `bootkube-worker`), and define 3 groups which match machines by MAC address to a profile. These resources declare that each machine should PXE boot and install Container Linux to disk. `node1` will provision itself as a controller, while `node2` and `noe3` provision themselves as workers.

Fetch the [profiles](../README.md#modules) Terraform [module](https://www.terraform.io/docs/modules/index.html) which let's you use common machine profiles maintained in the matchbox repo (like `bootkube`).

```sh
$ terraform get
```

Plan and apply to create the resources on Matchbox.

```sh
$ terraform plan
Plan: 10 to add, 0 to change, 0 to destroy.
$ terraform apply
Apply complete! Resources: 10 added, 0 changed, 0 destroyed.
```

Note: The `cached-container-linux-install` profile will PXE boot and install Container Linux from matchbox [assets](https://github.com/coreos/matchbox/blob/master/Documentation/api.md#assets). If you have not populated the assets cache, use the `container-linux-install` profile to use public images (slower).

## Machines

Power on each machine (with PXE boot device on next boot). Machines should network boot, install Container Linux to disk, reboot, and provision themselves as bootkube controllers or workers.

```sh
$ ipmitool -H node1.example.com -U USER -P PASS chassis bootdev pxe
$ ipmitool -H node1.example.com -U USER -P PASS power on
```

For local QEMU/KVM development, create the QEMU/KVM VMs.

```sh
$ sudo ./scripts/libvirt create
$ sudo ./scripts/libvirt [start|reboot|shutdown|poweroff|destroy]
```

## bootkube

*This section will soon be automated by terraform*

Install [bootkube](https://github.com/kubernetes-incubator/bootkube/releases) v0.4.2 and add it somewhere on your PATH.

```sh
bootkube version
Version v0.4.2
```

Use the `bootkube` tool to render Kubernetes manifests and credentials into an `--asset-dir`. Later, `bootkube` will schedule these manifests during bootstrapping and the credentials will be used to access your cluster.

```sh
bootkube render --asset-dir=assets --api-servers=https://node1.example.com:443 --api-server-alt-names=DNS=node1.example.com --etcd-servers=http://127.0.0.1:2379
```

Secure copy the kubeconfig to /etc/kubernetes/kubeconfig on every node which will path activate the `kubelet.service`.

```
for node in 'node1' 'node2' 'node3'; do
    scp assets/auth/kubeconfig core@$node.example.com:/home/core/kubeconfig
    ssh core@$node.example.com 'sudo mv kubeconfig /etc/kubernetes/kubeconfig'
done
```

Secure copy the bootkube generated assets to any controller node and run bootkube-start.

```
scp -r assets core@node1.example.com:/home/core
ssh core@node1.example.com 'sudo mv assets /opt/bootkube/assets && sudo systemctl start bootkube'
```

Optionally watch bootkube start the Kubernetes control plane.

```
$ ssh core@node1.example.com 'journalctl -f -u bootkube'
[  299.241291] bootkube[5]:     Pod Status:     kube-api-checkpoint     Running
[  299.241618] bootkube[5]:     Pod Status:          kube-apiserver     Running
[  299.241804] bootkube[5]:     Pod Status:          kube-scheduler     Running
[  299.241993] bootkube[5]:     Pod Status: kube-controller-manager     Running
[  299.311743] bootkube[5]: All self-hosted control plane components successfully started
```

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster. Verify that the cluster is accessible and that the apiserver, scheduler, and controller-manager are running as pods.

```sh
$ KUBECONFIG=assets/auth/kubeconfig
$ kubectl get nodes
NAME                STATUS    AGE
node1.example.com   Ready     3m
node2.example.com   Ready     3m
node3.example.com   Ready     3m

$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                       READY     STATUS    RESTARTS   AGE
kube-system   checkpoint-installer-p8g8r                 1/1       Running   1          13m
kube-system   kube-apiserver-s5gnx                       1/1       Running   1          41s
kube-system   kube-controller-manager-3438979800-jrlnd   1/1       Running   1          13m
kube-system   kube-controller-manager-3438979800-tkjx7   1/1       Running   1          13m
kube-system   kube-dns-4101612645-xt55f                  4/4       Running   4          13m
kube-system   kube-flannel-pl5c2                         2/2       Running   0          13m
kube-system   kube-flannel-r9t5r                         2/2       Running   3          13m
kube-system   kube-flannel-vfb0s                         2/2       Running   4          13m
kube-system   kube-proxy-cvhmj                           1/1       Running   0          13m
kube-system   kube-proxy-hf9mh                           1/1       Running   1          13m
kube-system   kube-proxy-kpl73                           1/1       Running   1          13m
kube-system   kube-scheduler-694795526-1l23b             1/1       Running   1          13m
kube-system   kube-scheduler-694795526-fks0b             1/1       Running   1          13m
kube-system   pod-checkpointer-node1.example.com         1/1       Running   2          10m
```

Try restarting machines or deleting pods to see that the cluster is resilient to failures.

## Going Further

Learn more about [matchbox](../../../Documentation/matchbox.md) or explore the other [example](../) clusters.
