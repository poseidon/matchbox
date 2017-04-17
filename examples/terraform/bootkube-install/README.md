# Self-hosted Kubernetes

The self-hosted Kubernetes example provisions a 3 node "self-hosted" Kubernetes v1.6.1 cluster. On-host kubelets wait for an apiserver to become reachable, then yield to kubelet pods scheduled via daemonset. [bootkube](https://github.com/kubernetes-incubator/bootkube) is run on any controller to bootstrap a temporary apiserver which schedules control plane components as pods before exiting. An etcd cluster backs Kubernetes and coordinates CoreOS auto-updates (enabled for disk installs).

## Requirements

* Create a PXE network boot environment (e.g. with `coreos/dnsmasq`)
* Run a `matchbox` service with the gRPC API enabled
* 3 machines with known DNS names and MAC addresses for this example
* Matchbox provider credentials: a `client.crt`, `client.key`, and `ca.crt`.

Install [bootkube](https://github.com/kubernetes-incubator/bootkube/releases) v0.4.0 and add it somewhere on your PATH.

```sh
bootkube version
Version v0.4.0
```

Use the `bootkube` tool to render Kubernetes manifests and credentials into an `--asset-dir`. Later, `bootkube` will schedule these manifests during bootstrapping and the credentials will be used to access your cluster.

```sh
bootkube render --asset-dir=assets --api-servers=https://node1.example.com:443 --api-server-alt-names=DNS=node1.example.com
```

## Infrastructure

Plan and apply terraform configurations. Create `bootkube-controller`, `bootkube-worker`, and `install-reboot` profiles and Container Linux configs. Create matcher groups for `node1.example.com`, `node2.example.com`, and `node3.example.com`.

```
cd examples/bootkube-install
terraform plan
terraform apply
```

Power on each machine and wait for it to PXE boot, install CoreOS to disk, and provision itself.

## Bootstrap

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

Optionally watch the Kubernetes control plane bootstrapping with the bootkube temporary api-server. You will see quite a bit of output.

```
$ ssh core@node1.example.com 'journalctl -f -u bootkube'
[  299.241291] bootkube[5]:     Pod Status:     kube-api-checkpoint     Running
[  299.241618] bootkube[5]:     Pod Status:          kube-apiserver     Running
[  299.241804] bootkube[5]:     Pod Status:          kube-scheduler     Running
[  299.241993] bootkube[5]:     Pod Status: kube-controller-manager     Running
[  299.311743] bootkube[5]: All self-hosted control plane components successfully started
```

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster. Verify that the cluster is accessible and that the kubelet, apiserver, scheduler, and controller-manager are running as pods.

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

Try deleting pods to see that the cluster is resilient to failures and machine restarts (CoreOS auto-updates).
