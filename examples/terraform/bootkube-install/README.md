# Kubernetes

The Kubernetes example shows how to use Matchbox to network boot and provision a 3 node Kubernetes v1.13.2 cluster. This example uses [Terraform](https://www.terraform.io/intro/index.html) and a module provided by [Typhoon](https://github.com/poseidon/typhoon) to describe cluster resources. [kubernetes-incubator/bootkube](https://github.com/kubernetes-incubator/bootkube) is run once to bootstrap the Kubernetes control plane.

## Requirements

Follow the getting started [tutorial](../../../Documentation/getting-started.md) to learn about matchbox and set up an environment that meets the requirements:

* Matchbox v0.6+ [installation](../../../Documentation/deployment.md) with gRPC API enabled
* Matchbox provider credentials `client.crt`, `client.key`, and `ca.crt`
* PXE [network boot](../../../Documentation/network-setup.md) environment
* Terraform v0.11.x, [terraform-provider-matchbox](https://github.com/coreos/terraform-provider-matchbox), and [terraform-provider-ct](https://github.com/coreos/terraform-provider-ct) installed locally
* Machines with known DNS names and MAC addresses

If you prefer to provision QEMU/KVM VMs on your local Linux machine, set up the matchbox [development environment](../../../Documentation/getting-started-docker.md).

```sh
sudo ./scripts/devnet create
```

## Terraform Setup

Install [Terraform](https://www.terraform.io/downloads.html) v0.11.x on your system.

```sh
$ terraform version
Terraform v0.11.7
```

Add the [terraform-provider-matchbox](https://github.com/coreos/terraform-provider-matchbox) plugin binary for your system to `~/.terraform.d/plugins/`, noting the final name.

```sh
wget https://github.com/coreos/terraform-provider-matchbox/releases/download/v0.2.2/terraform-provider-matchbox-v0.2.2-linux-amd64.tar.gz
tar xzf terraform-provider-matchbox-v0.2.2-linux-amd64.tar.gz
mv terraform-provider-matchbox-v0.2.2-linux-amd64/terraform-provider-matchbox ~/.terraform.d/plugins/terraform-provider-matchbox_v0.2.2
```

Add the [terraform-provider-ct](https://github.com/coreos/terraform-provider-ct) plugin binary for your system to `~/.terraform.d/plugins/`, noting the final name.

```sh
wget https://github.com/coreos/terraform-provider-ct/releases/download/v0.3.0/terraform-provider-ct-v0.3.0-linux-amd64.tar.gz
tar xzf terraform-provider-ct-v0.3.0-linux-amd64.tar.gz
mv terraform-provider-ct-v0.3.0-linux-amd64/terraform-provider-ct ~/.terraform.d/plugins/terraform-provider-ct_v0.3.0
```

## Usage

Clone the [matchbox](https://github.com/coreos/matchbox) project and take a look at the cluster examples.

```sh
$ git clone https://github.com/coreos/matchbox.git
$ cd matchbox/examples/terraform/bootkube-install
```

Configure the Matchbox provider to use your Matchbox API endpoint and client certificate in a `providers.tf` file.

```
provider "matchbox" {
  version = "0.2.2"
  endpoint    = "matchbox.example.com:8081"
  client_cert = "${file("~/.matchbox/client.crt")}"
  client_key  = "${file("~/.matchbox/client.key")}"
  ca          = "${file("~/.matchbox/ca.crt")}"
}

provider "ct" {
  version = "0.3.0"
}
...
```

Copy the `terraform.tfvars.example` file to `terraform.tfvars`. It defines a few variables needed for examples. Set your `ssh_authorized_key` to use in the cluster definition.

Note: With `cached_install="true"`, machines will PXE boot and install Container Linux from matchbox [assets](https://github.com/coreos/matchbox/blob/master/Documentation/api.md#assets). For convenience, `scripts/get-coreos` can download needed images.

## Terraform

Initialize Terraform from the `bootkube-install` directory.

```sh
terraform init
```

Plan the resources to be created.

```sh
$ terraform plan
Plan: 75 to add, 0 to change, 0 to destroy.
```

Terraform will configure matchbox with profiles (e.g. `cached-container-linux-install`, `bootkube-controller`, `bootkube-worker`) and add groups to match machines by MAC address to a profile. These resources declare that each machine should PXE boot and install Container Linux to disk. `node1` will provision itself as a controller, while `node2` and `node3` provision themselves as workers.

The module referenced in `cluster.tf` will also generate bootkube assets to `assets_dir` (exactly like the [bootkube](https://github.com/kubernetes-incubator/bootkube) binary would). These assets include Kubernetes bootstrapping and control plane manifests as well as a kubeconfig you can use to access the cluster. 

### ssh-agent

Initial bootstrapping requires `bootkube.service` be started on one controller node. Terraform uses `ssh-agent` to automate this step. Add your SSH private key to `ssh-agent`, otherwise `terraform apply` will hang.

```sh
ssh-add ~/.ssh/id_rsa
ssh-add -L
```

### Apply

Apply the changes.

```sh
$ terraform apply
module.cluster.null_resource.copy-secrets.0: Still creating... (5m0s elapsed)
module.cluster.null_resource.copy-secrets.1: Still creating... (5m0s elapsed)
module.cluster.null_resource.copy-secrets.2: Still creating... (5m0s elapsed)
...
module.cluster.null_resource.bootkube-start: Still creating... (8m40s elapsed)
...
```

Apply will then loop until it can successfully copy credentials to each machine and start the one-time Kubernetes bootstrap service. Proceed to the next step while this loops.

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

## Verify

[Install kubectl](https://coreos.com/kubernetes/docs/latest/configure-kubectl.html) on your laptop. Use the generated kubeconfig to access the Kubernetes cluster. Verify that the cluster is accessible and that the apiserver, scheduler, and controller-manager are running as pods.

```sh
$ export KUBECONFIG=assets/auth/kubeconfig
$ kubectl get nodes
NAME                STATUS    AGE       VERSION
node1.example.com   Ready     11m       v1.13.2
node2.example.com   Ready     11m       v1.13.2
node3.example.com   Ready     11m       v1.13.2

$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                       READY     STATUS    RESTARTS   AGE
kube-system   coredns-1187388186-mx9rt                   3/3       Running   0          11m
kube-system   coredns-1187388186-dsfk3                   3/3       Running   0          11m
kube-system   flannel-fqp7f                              2/2       Running   1          11m
kube-system   flannel-gnjrm                              2/2       Running   0          11m
kube-system   flannel-llbgt                              2/2       Running   0          11m
kube-system   kube-apiserver-7336w                       1/1       Running   0          11m
kube-system   kube-controller-manager-3271970485-b9chx   1/1       Running   0          11m
kube-system   kube-controller-manager-3271970485-v30js   1/1       Running   1          11m
kube-system   kube-proxy-50sd4                           1/1       Running   0          11m
kube-system   kube-proxy-bczhp                           1/1       Running   0          11m
kube-system   kube-proxy-mp2fw                           1/1       Running   0          11m
kube-system   kube-scheduler-3895335239-fd3l7            1/1       Running   1          11m
kube-system   kube-scheduler-3895335239-hfjv0            1/1       Running   0          11m
kube-system   pod-checkpointer-wf65d                     1/1       Running   0          11m
kube-system   pod-checkpointer-wf65d-node1.example.com   1/1       Running   0          11m
```

## Optional

Several Terraform module variables can override cluster defaults. [Check upstream](https://typhoon.psdn.io/bare-metal/) for the full list of options.

```hcl
...
cached_install = "false"
install_disk = "/dev/sda"
networking = "calico"
```

## Addons

Install **important** cluster [addons](../../../Documentation/cluster-addons.md).

## Going Further

Learn more about [matchbox](../../../Documentation/matchbox.md) or explore the other [example](../) clusters.
