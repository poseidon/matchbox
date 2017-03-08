
# Upgrading self-hosted Kubernetes

[Self-hosted](bootkube.md) Kubernetes clusters schedule Kubernetes components such as the apiserver, kubelet, scheduler, and controller-manager as pods like other applications (except with node selectors). This allows Kubernetes level operations to be performed to upgrade clusters in place, rather than by re-provisioning.

Let's upgrade a self-hosted Kubernetes v1.4.1 cluster to v1.4.3 as an example.

## Inspect

Show the control plane daemonsets and deployments which will need to be updated.

```sh
$ kubectl get daemonsets -n=kube-system
NAME             DESIRED   CURRENT   NODE-SELECTOR   AGE
kube-apiserver   1         1         master=true     5m
kube-proxy       3         3         <none>          5m
kubelet          3         3         <none>          5m

$ kubectl get deployments -n=kube-system
NAME                      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
kube-controller-manager   1         1         1            1           5m
kube-dns-v20              1         1         1            1           5m
kube-scheduler            1         1         1            1           5m
```

Check the current Kubernetes version.

```sh
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"4", GitVersion:"v1.4.0", GitCommit:"a16c0a7f71a6f93c7e0f222d961f4675cd97a46b", GitTreeState:"clean", BuildDate:"2016-09-26T18:16:57Z", GoVersion:"go1.6.3", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"4", GitVersion:"v1.4.1+coreos.0", GitCommit:"b7a02f46b972c5211e5c04fdb1d5b86ac16c00eb", GitTreeState:"clean", BuildDate:"2016-10-11T20:13:55Z", GoVersion:"go1.6.3", Compiler:"gc", Platform:"linux/amd64"}
```

In this case, Kubernetes is `v1.4.1+coreos.0` and our goal is to upgrade to `v1.4.3+coreos.0`. First, update the control plane pods. Then the kubelets and proxies on all nodes.

**Tip**: Follow along with a QEMU/KVM self-hosted Kubernetes cluster the first time, before upgrading your production bare-metal clusters ([tutorial](bootkube.md)).

## Control Plane

### kube-apiserver

Edit the kube-apiserver daemonset. Change the container image name to `quay.io/coreos/hyperkube:v1.4.3_coreos.0`.

```sh
$ kubectl edit daemonset kube-apiserver -n=kube-system
```

Since daemonsets don't yet support rolling, manually delete each apiserver one by one and wait for each to be re-scheduled.

```sh
$ kubectl get pods -n=kube-system
# WARNING: Self-hosted Kubernetes is still new and this may fail
$ kubectl delete pod kube-apiserver-s62kb -n=kube-system
```

If you only have one, your cluster will be temporarily unavailable. Remember the Hyperkube image is quite large and this can take a minute.

```sh
$ kubectl get pods -n=kube-system
NAME                                       READY     STATUS    RESTARTS   AGE
kube-api-checkpoint-node1.example.com      1/1       Running   0          12m
kube-apiserver-vyg3t                       2/2       Running   0          2m
kube-controller-manager-1510822774-qebia   1/1       Running   2          12m
kube-dns-v20-3531996453-0tlv9              3/3       Running   0          12m
kube-proxy-8jthl                           1/1       Running   0          12m
kube-proxy-bnvgy                           1/1       Running   0          12m
kube-proxy-gkyx8                           1/1       Running   0          12m
kube-scheduler-2099299605-67ezp            1/1       Running   2          12m
kubelet-exe5k                              1/1       Running   0          12m
kubelet-p3g98                              1/1       Running   0          12m
kubelet-quhhg                              1/1       Running   0          12m
```

### kube-scheduler

Edit the scheduler deployment to rolling update the scheduler. Change the container image name for the hyperkube.

```sh
$ kubectl edit deployments kube-scheduler -n=kube-system
```

Wait for the schduler to be deployed.

### kube-controller-manager

Edit the controller-manager deployment to rolling update the controller manager. Change the container image name for the hyperkube.

```sh
$ kubectl edit deployments kube-controller-manager -n=kube-system
```

Wait for the controller manager to be deployed.

```sh
$ kubectl get pods -n=kube-system
NAME                                       READY     STATUS    RESTARTS   AGE
kube-api-checkpoint-node1.example.com      1/1       Running   0          28m
kube-apiserver-vyg3t                       2/2       Running   0          18m
kube-controller-manager-1709527928-zj8c4   1/1       Running   0          4m
kube-dns-v20-3531996453-0tlv9              3/3       Running   0          28m
kube-proxy-8jthl                           1/1       Running   0          28m
kube-proxy-bnvgy                           1/1       Running   0          28m
kube-proxy-gkyx8                           1/1       Running   0          28m
kube-scheduler-2255275287-hti6w            1/1       Running   0          6m
kubelet-exe5k                              1/1       Running   0          28m
kubelet-p3g98                              1/1       Running   0          28m
kubelet-quhhg                              1/1       Running   0          28m
```

### Verify

At this point, the control plane components have been upgraded to v1.4.3.

```sh
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"4", GitVersion:"v1.4.0", GitCommit:"a16c0a7f71a6f93c7e0f222d961f4675cd97a46b", GitTreeState:"clean", BuildDate:"2016-09-26T18:16:57Z", GoVersion:"go1.6.3", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"4", GitVersion:"v1.4.3+coreos.0", GitCommit:"7819c84f25e8c661321ee80d6b9fa5f4ff06676f", GitTreeState:"clean", BuildDate:"2016-10-17T21:19:17Z", GoVersion:"go1.6.3", Compiler:"gc", Platform:"linux/amd64"}
```

Finally, upgrade the kubelets and kube-proxies.

## kubelet and kube-proxy

Show the current kubelet and kube-proxy version on each node.

```sh
$ kubectl get nodes -o yaml | grep 'kubeletVersion\|kubeProxyVersion'
    kubeProxyVersion: v1.4.1+coreos.0
    kubeletVersion: v1.4.1+coreos.0
    kubeProxyVersion: v1.4.1+coreos.0
    kubeletVersion: v1.4.1+coreos.0
    kubeProxyVersion: v1.4.1+coreos.0
    kubeletVersion: v1.4.1+coreos.0
```

Edit the kubelet and kube-proxy daemonsets. Change the container image name for the hyperkube.

```sh
$ kubectl edit daemonset kubelet -n=kube-system
$ kubectl edit daemonset kube-proxy -n=kube-system
```

Since daemonsets don't yet support rolling, manually delete each kubelet and each kube-proxy. The daemonset controller will create new (upgraded) replics.

```sh
$ kubectl get pods -n=kube-system
$ kubectl delete pod kubelet-quhhg
...repeat
$ kubectl delete pod kube-proxy-8jthl -n=kube-system
...repeat

$ kubectl get pods -n=kube-system
NAME                                       READY     STATUS    RESTARTS   AGE
kube-api-checkpoint-node1.example.com      1/1       Running   0          1h
kube-apiserver-vyg3t                       2/2       Running   0          1h
kube-controller-manager-1709527928-zj8c4   1/1       Running   0          47m
kube-dns-v20-3531996453-0tlv9              3/3       Running   0          1h
kube-proxy-6dbne                           1/1       Running   0          1s
kube-proxy-sm4jv                           1/1       Running   0          8s
kube-proxy-xmuao                           1/1       Running   0          14s
kube-scheduler-2255275287-hti6w            1/1       Running   0          49m
kubelet-hfdwr                              1/1       Running   0          38s
kubelet-oia47                              1/1       Running   0          52s
kubelet-s6dab                              1/1       Running   0          59s
```

## Verify

Verify that the kubelet and kube-proxy on each node have been upgraded.

```sh
$ kubectl get nodes -o yaml | grep 'kubeletVersion\|kubeProxyVersion'
      kubeProxyVersion: v1.4.3+coreos.0
      kubeletVersion: v1.4.3+coreos.0
      kubeProxyVersion: v1.4.3+coreos.0
      kubeletVersion: v1.4.3+coreos.0
      kubeProxyVersion: v1.4.3+coreos.0
      kubeletVersion: v1.4.3+coreos.0
```

Now, Kubernetes components have been upgraded to a new version of Kubernetes!

## Going further

Bare-metal or virtualized self-hosted Kubernetes clusters can be upgraded in place in 5-10 minutes. Here is a bare-metal example:

```sh
$ kubectl -n=kube-system get pods
NAME                                       READY     STATUS    RESTARTS   AGE
kube-api-checkpoint-ibm0.lab.dghubble.io   1/1       Running   0          2d
kube-apiserver-j6atn                       2/2       Running   0          5m
kube-controller-manager-1709527928-y05n5   1/1       Running   0          1m
kube-dns-v20-3531996453-zwbl8              3/3       Running   0          2d
kube-proxy-e49p5                           1/1       Running   0          14s
kube-proxy-eu5dc                           1/1       Running   0          8s
kube-proxy-gjrzq                           1/1       Running   0          3s
kube-scheduler-2255275287-96n56            1/1       Running   0          2m
kubelet-9ob0e                              1/1       Running   0          19s
kubelet-bvwp0                              1/1       Running   0          14s
kubelet-xlrql                              1/1       Running   0          24s
```

Check upstream for updates to addons like `kube-dns` or `kube-dashboard` and update them like any other applications. Some kube-system components use version labels and you may wish to clean those up as well.
