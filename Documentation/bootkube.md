
# Self-Hosted Kubernetes

The self-hosted Kubernetes examples provision a 3 node cluster with an etcd cluster, flannel, and a special "runonce" host Kublet. The CoreOS [bootkube](https://github.com/coreos/bootkube) reverse-tunnel tool is used to start the kubelet, api server, scheduler, and controller manager as pods which can be managed via kubectl (self-hosted). The `bootkube` example PXE boots and provisions CoreOS nodes for this purpose, while `bootkube-install` does the CoreOS install to disk as well.

## Requirements

Ensure that you've gone through the `bootcfg` [Getting Started Guide](getting-started-rkt.md) guide and understand the basics. In particular, you should be able to:

* Use rkt to start `bootcfg`
* Create a network boot environment with `coreos/dnsmasq`
* Create the example libvirt client VMs

Build and install [bootkube](https://github.com/coreos/bootkube).

## Examples

The examples statically assign IP addresses for client VMs on the `metal0` CNI bridge used by rkt. You can use the same examples for real hardware, but you'll need to update addresses and the NIC name.

* [bootkube](../examples/groups/bootkube) - iPXE boot a bootkube-ready cluster (use rkt)
* [bootkube-install](../examples/groups/bootkube-install) - Install a bootkube-ready cluster (use rkt)

### Assets

Download the CoreOS PXE image referenced in the target [profile](../examples/profiles).

    ./scripts/get-coreos alpha 983.0.0

Use the `bootkube` tool to render Kubernetes manifests and credentials into an output directory.

    bootkube render --outdir assets --apiserver-cert-ip-addrs=172.15.0.21,10.3.0.1 --api-servers=https://172.15.0.21:6443 --etcd-servers=http://172.15.0.21:2379

Manually copy the `certificate-authority-data` from the generated `assets/auth/kubeconfig.yaml` file and paste it in each machine group definition in `examples/groups/bootkube` or `examples/groups/bootkube-install` (i.e. `node1json`, `node2.json`, `node3.json`).

Manually add your SSH public key to each machine group definition [as shown](../examples/README.md#ssh-keys).

## Containers

Run the latest `bootcfg` ACI with rkt.

    sudo rkt --insecure-options=image run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug

Create the machine groups at runtime (or you could mount them).

    bootcmd --endpoints 172.15.0.2:8081 group create -f examples/groups/bootkube-install/install.json
    bootcmd --endpoints 172.15.0.2:8081 group create -f examples/groups/bootkube-install/node1.json
    bootcmd --endpoints 172.15.0.2:8081 group create -f examples/groups/bootkube-install/node2.json
    bootcmd --endpoints 172.15.0.2:8081 group create -f examples/groups/bootkube-install/node3.json

Check the loaded profiles and groups.

    bootcmd --endpoints 172.15.0.2:8081 group list
    bootcmd --endpoints 172.15.0.2:8081 profile list

Create a network boot environment with `coreos/dnsmasq`. Finally, create client VMs with `scripts/libvirt` or power-on your physical machines. Client machines should boot and provision themselves.

Revisit [bootcfg with rkt](getting-started-rkt.md) for help.

## bootkube

You're ready to use the [bootkube](https://github.com/coreos/bootkube) control plane tool from your laptop. This will run a temporary control plane on your machine and make it available to your cluster via an SSH reverse-tunnel to bootstrap the control plane.

    bootkube start --ssh-user=core --ssh-keyfile=${HOME}/.ssh/id_rsa --remote-address=172.15.0.21:22 --remote-etcd-address=172.15.0.21:2379 --manifest-dir=assets/manifests --apiserver-key=assets/tls/apiserver.key --apiserver-cert=assets/tls/apiserver.crt --ca-cert=assets/tls/ca.crt --service-account-key=assets/tls/service-account.key --token-auth-file=assets/auth/token-auth.csv

After a few minutes, the kubelet pod started by `bootkube` will take over in place of the host kubelet.

    I0425 12:38:23.746330   29538 status.go:87] Pod status kubelet: Running
    I0425 12:38:23.746361   29538 status.go:87] Pod status kube-apiserver: Running
    I0425 12:38:23.746370   29538 status.go:87] Pod status kube-scheduler: Running
    I0425 12:38:23.746378   29538 status.go:87] Pod status kube-controller-manager: Running

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster. Verify that the cluster is accessible and that the kubelet, apiserver, scheduler, and controller-manager are running as pods.

    $ kubectl --kubeconfig=assets/auth/kubeconfig.yaml get nodes
    NAME          STATUS    AGE
    172.15.0.21   Ready     3m
    172.15.0.22   Ready     3m
    172.15.0.23   Ready     3m

    $ kubectl --kubeconfig=assets/auth/kubeconfig.yaml get pods --all-namespaces
    NAMESPACE     NAME                            READY     STATUS    RESTARTS   AGE
    kube-system   kube-apiserver-z6f9e            1/1       Running   0          3m
    kube-system   kube-controller-manager-slia7   1/1       Running   0          3m
    kube-system   kube-dns-v11-1943182923-7embo   4/4       Running   0          3m
    kube-system   kube-proxy-ed6af                1/1       Running   0          3m
    kube-system   kube-proxy-fbhfq                1/1       Running   0          3m
    kube-system   kube-proxy-fzy8r                1/1       Running   0          3m
    kube-system   kube-scheduler-djgdh            1/1       Running   0          3m
    kube-system   kubelet-crywn                   1/1       Running   0          3m
    kube-system   kubelet-rmdjq                   1/1       Running   0          3m
    kube-system   kubelet-wjj0g                   1/1       Running   0          3m

