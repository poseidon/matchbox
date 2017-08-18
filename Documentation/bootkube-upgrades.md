# Upgrading self-hosted Kubernetes

CoreOS Kubernetes clusters "self-host" the apiserver, scheduler, controller-manager, flannel, kube-dns, and kube-proxy as Kubernetes pods, like ordinary applications (except with taint tolerations). This allows upgrades to be performed in-place using (mostly) `kubectl` as an alternative to re-provisioning.

Let's upgrade a Kubernetes v1.6.6 cluster to v1.6.7 as an example.

## Stability

This guide shows how to attempt a in-place upgrade of a Kubernetes cluster setup via the [examples](../examples). It does not provide exact diffs, migrations between breaking changes, the stability of a fresh re-provision, or any guarantees. Evaluate whether in-place updates are appropriate for your Kubernetes cluster and be prepared to perform a fresh re-provision if something goes wrong, especially between Kubernetes minor releases (e.g. 1.6 to 1.7).

Matchbox Kubernetes examples provide a vanilla Kubernetes cluster with only free (as in freedom and cost) software components. If you require currated updates, migrations, or guarantees for production, consider [Tectonic](https://coreos.com/tectonic/) by CoreOS.

**Note: Tectonic users should NOT manually upgrade. Follow the [Tectonic docs](https://coreos.com/tectonic/docs/latest/admin/upgrade.html)**

## Inspect

Show the control plane daemonsets and deployments which will need to be updated.

```sh
$ kubectl get daemonsets -n=kube-system
NAME                             DESIRED   CURRENT   READY     UP-TO-DATE   AVAILABLE   NODE-SELECTOR                     AGE
kube-apiserver                   1         1         1         1            1           node-role.kubernetes.io/master=   21d
kube-etcd-network-checkpointer   1         1         1         1            1           node-role.kubernetes.io/master=   21d
kube-flannel                     4         4         4         4            4           <none>                            21d
kube-proxy                       4         4         4         4            4           <none>                            21d
pod-checkpointer                 1         1         1         1            1           node-role.kubernetes.io/master=   21d

$ kubectl get deployments -n=kube-system
kube-controller-manager           2         2         2            2           21d
kube-dns                          1         1         1            1           21d
kube-scheduler                    2         2         2            2           21d
```

Check the current Kubernetes version.

```sh
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"6", GitVersion:"v1.6.2", GitCommit:"477efc3cbe6a7effca06bd1452fa356e2201e1ee", GitTreeState:"clean", BuildDate:"2017-04-19T20:33:11Z", GoVersion:"go1.7.5", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"6", GitVersion:"v1.6.6+coreos.1", GitCommit:"42a5c8b99c994a51d9ceaed5d0254f177e97d419", GitTreeState:"clean", BuildDate:"2017-06-21T01:10:07Z", GoVersion:"go1.7.6", Compiler:"gc", Platform:"linux/amd64"}
```

```sh
$ kubectl get nodes
NAME                               STATUS    AGE       VERSION
node1.example.com                  Ready     21d       v1.6.6+coreos.1
node2.example.com                  Ready     21d       v1.6.6+coreos.1
node3.example.com                  Ready     21d       v1.6.6+coreos.1
node4.example.com                  Ready     21d       v1.6.6+coreos.1
```

## Strategy

Update control plane components with `kubectl`. Then update the `kubelet` systemd unit on each host.

Prepare the changes to the Kubernetes manifests by generating assets for a target Kubernetes cluster (e.g. bootkube `v0.5.0` produces Kubernetes 1.6.6 and bootkube `v0.5.1` produces Kubernetes 1.6.7). Choose the tool used during creation of the cluster:

* [kubernetes-incubator/bootkube](https://github.com/kubernetes-incubator/bootkube) - install the `bootkube` binary for the target version and render assets
* [poseidon/bootkube-terraform](https://github.com/poseidon/bootkube-terraform) - checkout the tag for the target version and `terraform apply` to render assets

Diff the generated assets against the assets used when originally creating the cluster. In simple cases, you may only need to bump the hyperkube image. In more complex cases, some manifests may have new flags or configuration.

## Control Plane

### kube-apiserver

Edit the `kube-apiserver` daemonset to rolling update the apiserver.

```sh
$ kubectl edit daemonset kube-apiserver -n=kube-system
```

If you only have one apiserver, the cluster may be momentarily unavailable.

### kube-scheduler

Edit the `kube-scheduler` deployment to rolling update the scheduler.

```sh
$ kubectl edit deployments kube-scheduler -n=kube-system
```

### kube-controller-manager

Edit the `kube-controller-manager` deployment to rolling update the controller manager.

```sh
$ kubectl edit deployments kube-controller-manager -n=kube-system
```

### kube-proxy

Edit the `kube-proxy` daemonset to rolling update the proxy.

```sh
$ kubectl edit daemonset kube-proxy -n=kube-system
```

### Others

If there are changes between the prior version and target version manifests, update the `kube-dns` deployment, `kube-flannel` daemonset, or `pod-checkpointer` daemonset.

### Verify

Verify the control plane components updated.

```sh
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"6", GitVersion:"v1.6.2", GitCommit:"477efc3cbe6a7effca06bd1452fa356e2201e1ee", GitTreeState:"clean", BuildDate:"2017-04-19T20:33:11Z", GoVersion:"go1.7.5", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"6", GitVersion:"v1.6.7+coreos.0", GitCommit:"c8c505ee26ac3ab4d1dff506c46bc5538bc66733", GitTreeState:"clean", BuildDate:"2017-07-06T17:38:33Z", GoVersion:"go1.7.6", Compiler:"gc", Platform:"linux/amd64"}
```

```sh
$ kubectl get nodes
NAME                               STATUS    AGE       VERSION
node1.example.com                  Ready     21d       v1.6.7+coreos.0
node2.example.com                  Ready     21d       v1.6.7+coreos.0
node3.example.com                  Ready     21d       v1.6.7+coreos.0
node4.example.com                  Ready     21d       v1.6.7+coreos.0
```

## kubelet

SSH to each node and update `/etc/kubernetes/kubelet.env`. Restart the `kubelet.service`.

```sh
ssh core@node1.example.com
sudo vim /etc/kubernetes/kubelet.env
sudo systemctl restart kubelet
```

### Verify

Verify the kubelet and kube-proxy of each node updated.

```sh
$ kubectl get nodes -o yaml | grep 'kubeletVersion\|kubeProxyVersion'
      kubeProxyVersion: v1.6.7+coreos.0
      kubeletVersion: v1.6.7+coreos.0
      kubeProxyVersion: v1.6.7+coreos.0
      kubeletVersion: v1.6.7+coreos.0
      kubeProxyVersion: v1.6.7+coreos.0
      kubeletVersion: v1.6.7+coreos.0
      kubeProxyVersion: v1.6.7+coreos.0
      kubeletVersion: v1.6.7+coreos.0
```

Kubernetes control plane components have been successfully updated!
