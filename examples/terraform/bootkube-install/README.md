# Kubernetes

The Kubernetes example shows how to use Matchbox to network boot and provision a 3 node Kubernetes v1.10.3 cluster. This example uses [Terraform](https://www.terraform.io/intro/index.html) and a module provided by [Typhoon](https://github.com/poseidon/typhoon) to describe cluster resources. [kubernetes-incubator/bootkube](https://github.com/kubernetes-incubator/bootkube) is run once to bootstrap the Kubernetes control plane.

## Requirements

Follow the getting started [tutorial](../../../Documentation/getting-started.md) to learn about matchbox and set up an environment that meets the requirements:

* Matchbox v0.6+ [installation](../../../Documentation/deployment.md) with gRPC API enabled
* Matchbox provider credentials `client.crt`, `client.key`, and `ca.crt`
* PXE [network boot](../../../Documentation/network-setup.md) environment
* Terraform v0.10.x or v0.11.x and [terraform-provider-matchbox](https://github.com/coreos/terraform-provider-matchbox) installed locally on your system
* Machines with known DNS names and MAC addresses

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

cluster_name = "demo"
os_channel   = "coreos-stable"
os_version   = "1576.5.0"
```

Provide an ordered list of controller names, MAC addresses, and domain names. Provide an ordered list of worker names, MAC addresses, and domain names.

```hcl
controller_names = ["node1"]
controller_macs = ["52:54:00:a1:9c:ae"]
controller_domains = ["node1.example.com"]
worker_names = ["node2", "node3"]
worker_macs = ["52:54:00:b2:2f:86", "52:54:00:c3:61:77"]
worker_domains = ["node2.example.com", "node3.example.com"]
```

Provide an `assets_dir` for generated manifests and a DNS name which you've setup to resolves to controller(s) (e.g. round-robin). Worker nodes and your kubeconfig will communicate via this endpoint.

```hcl
k8s_domain_name = "cluster.example.com"
asset_dir = "assets"
```

Note: The `cached-container-linux-install` profile will PXE boot and install Container Linux from matchbox [assets](https://github.com/coreos/matchbox/blob/master/Documentation/api.md#assets). If you have not populated the assets cache, use the `container-linux-install` profile to use public images (slower).

### Optional

You may set certain optional variables to override defaults. Set `networking` to either "flannel" or "calico" to set the networking provider. [Check upstream](https://typhoon.psdn.io/bare-metal/) for the full list of options.

```hcl
# Optional (defaults)
# cached_install = "false"
# install_disk = "/dev/sda"
# container_linux_oem = ""
# networking = "flannel"
```

The default is to create a Kubernetes cluster with 1 controller and 2 workers as an example, but check `multi-controller.tfvars.example` for an example which defines 3 controllers and 1 worker.

## Terraform

Initialize Terraform from the `bootkube-install` directory.

```sh
terraform init
```

Get or update Terraform modules.

```sh
$ terraform get            # downloads missing modules
$ terraform get --update   # updates all modules
Get: git::https://github.com/poseidon/typhoon (update)
Get: git::https://github.com/poseidon/bootkube-terraform.git?ref=v0.11.0 (update)
```

Plan the resources to be created.

```sh
$ terraform plan
Plan: 55 to add, 0 to change, 0 to destroy.
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
node1.example.com   Ready     11m       v1.10.0
node2.example.com   Ready     11m       v1.10.0
node3.example.com   Ready     11m       v1.10.0

$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                       READY     STATUS    RESTARTS   AGE
kube-system   kube-flannel-fqp7f                         2/2       Running   1          11m
kube-system   kube-flannel-gnjrm                         2/2       Running   0          11m
kube-system   kube-flannel-llbgt                         2/2       Running   0          11m
kube-system   kube-apiserver-7336w                       1/1       Running   0          11m
kube-system   kube-controller-manager-3271970485-b9chx   1/1       Running   0          11m
kube-system   kube-controller-manager-3271970485-v30js   1/1       Running   1          11m
kube-system   kube-dns-1187388186-mx9rt                  3/3       Running   0          11m
kube-system   kube-proxy-50sd4                           1/1       Running   0          11m
kube-system   kube-proxy-bczhp                           1/1       Running   0          11m
kube-system   kube-proxy-mp2fw                           1/1       Running   0          11m
kube-system   kube-scheduler-3895335239-fd3l7            1/1       Running   1          11m
kube-system   kube-scheduler-3895335239-hfjv0            1/1       Running   0          11m
kube-system   pod-checkpointer-wf65d                     1/1       Running   0          11m
kube-system   pod-checkpointer-wf65d-node1.example.com   1/1       Running   0          11m
```

## Addons

Install **important** cluster [addons](../../../Documentation/cluster-addons.md).

## Going Further

Learn more about [matchbox](../../../Documentation/matchbox.md) or explore the other [example](../) clusters.
