# Getting started

In this tutorial, we'll show how to use terraform with `matchbox` to provision Container Linux machines.

You'll install the `matchbox` service, setup a PXE network boot environment, and then use terraform configs to describe your infrastructure and the terraform CLI to create those resources on `matchbox`.

## matchbox

Install `matchbox` on a dedicated server or Kubernetes cluster. Generate TLS credentials and enable the gRPC API as directed. Save the `ca.crt`, `client.crt`, and `client.key` on your local machine (e.g. `~/.matchbox`).

* Installing on [Container Linux / other distros](deployment.md)
* Installing on [Kubernetes](deployment.md#kubernetes)
* Running with [rkt](deployment.md#rkt) / [docker](deployment.md#docker)

Verify the matchbox read-only HTTP endpoints are accessible.

```sh
$ curl http://matchbox.example.com:8080
matchbox
```

Verify your TLS client certificate and key can be used to access the gRPC API.

```sh
$ openssl s_client -connect matchbox.example.com:8081 \
  -CAfile ~/.matchbox/ca.crt \
  -cert ~/.matchbox/client.crt \
  -key ~/.matchbox/client.key
```

## Terraform

Install [Terraform][terraform-dl] v0.9+ on your system.

```sh
$ terraform version
Terraform v0.9.4
```

Add the `terraform-provider-matchbox` plugin binary on your system.

```sh
$ wget https://github.com/coreos/terraform-provider-matchbox/releases/download/v0.1.0/terraform-provider-matchbox-v0.1.0-linux-amd64.tar.gz
$ tar xzf terraform-provider-matchbox-v0.1.0-linux-amd64.tar.gz
```

Add the plugin to your `~/.terraformrc`.

```hcl
providers {
  matchbox = "/path/to/terraform-provider-matchbox"
}
```

## First cluster

Clone the matchbox source and take a look at the Terraform examples.

```sh
$ git clone https://github.com/coreos/matchbox.git
$ cd matchbox/examples/terraform
```

Let's start with the `simple-install` example. With `simple-install`, any machines which PXE boot from matchbox will install Container Linux to `dev/sda`, reboot, and have your SSH key set. Its not much of a cluster, but we'll get to that later.

```sh
$ cd simple-install
```

Configure the variables in `variables.tf` by creating a `terraform.tfvars` file.

```hcl
matchbox_http_endpoint = "http://matchbox.example.com:8080"
matchbox_rpc_endpoint = "matchbox.example.com:8081"
ssh_authorized_key = "YOUR_SSH_KEY"
```

Terraform can now interact with the matchbox service and create resources.

```sh
$ terraform plan
Plan: 4 to add, 0 to change, 0 to destroy.
```

Let's review the terraform config and learn a bit about matchbox.

#### Provider

Matchbox is configured as a provider platform for bare-metal resources.

```hcl
// Configure the matchbox provider
provider "matchbox" {
  endpoint = "${var.matchbox_rpc_endpoint}"
  client_cert = "${file("~/.matchbox/client.crt")}"
  client_key = "${file("~/.matchbox/client.key")}"
  ca         = "${file("~/.matchbox/ca.crt")}"
}
```

#### Profiles

Machine profiles specify the kernel, initrd, kernel args, Container Linux Config, Cloud-config, or other configs used to network boot and provision a bare-metal machine. This profile will PXE boot machines using the current stable Container Linux kernel and initrd (see [assets](api.md#assets) to learn about caching for speed) and supply a Container Linux Config specifying that a disk install and reboot should be performed. Learn more about [Container Linux configs](https://coreos.com/os/docs/latest/configuration.html).

```hcl
// Create a CoreOS-install profile
resource "matchbox_profile" "coreos-install" {
  name = "coreos-install"
  kernel = "https://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz"
  initrd = [
    "https://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe_image.cpio.gz"
  ]
  args = [
    "coreos.config.url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "coreos.first_boot=yes",
    "console=tty0",
    "console=ttyS0",
  ]
  container_linux_config = "${file("./cl/coreos-install.yaml.tmpl")}"
}
```

#### Groups

Matcher groups match machines based on labels like MAC, UUID, etc. to different profiles and templates in machine-specific values. This group does not have a `selector` block, so any machines which network boot from matchbox will match this group and be provisioned using the `coreos-install` profile. Machines are matched to the most specific matching group.

```hcl
resource "matchbox_group" "default" {
  name = "default"
  profile = "${matchbox_profile.coreos-install.name}"
  # no selector means all machines can be matched
  metadata {
    ignition_endpoint = "${var.matchbox_http_endpoint}/ignition"
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
```

### Apply

Apply the terraform configuration.

```sh
$ terraform apply
Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
```

Matchbox serves configs to machines and respects query parameters, if you're interested:

* iPXE default - [/ipxe](http://matchbox.example.com:8080/ipxe)
* Ignition default - [/ignition](http://matchbox.example.com:8080/ignition)
* Ignition post-install - [/ignition?os=installed](http://matchbox.example.com:8080/ignition?os=installed)
* GRUB default - [/grub](http://matchbox.example.com:8080/grub)

## Network

Matchbox can integrate with many on-premise network setups. It does not seek to be the DHCP server, TFTP server, or DNS server for the network. Instead, matchbox serves iPXE scripts and GRUB configs as the entrypoint for provisioning network booted machines. PXE clients are supported by chainloading iPXE firmware.

In the simplest case, an iPXE-enabled network can chain to matchbox,

```
# /var/www/html/ipxe/default.ipxe
chain http://matchbox.foo:8080/boot.ipxe
```

Read [network-setup.md](network-setup.md) for the complete range of options. Network admins have a great amount of flexibility:

* May keep using existing DHCP, TFTP, and DNS services
* May configure subnets, architectures, or specific machines to delegate to matchbox
* May place matchbox behind a menu entry (timeout and default to matchbox)

If you've never setup a PXE-enabled network before or you're trying to setup a home lab, checkout the [quay.io/coreos/dnsmasq](https://quay.io/repository/coreos/dnsmasq) container image [copy-paste examples](https://github.com/coreos/matchbox/blob/master/Documentation/network-setup.md#coreosdnsmasq) and see the section about [proxy-DHCP](https://github.com/coreos/matchbox/blob/master/Documentation/network-setup.md#proxy-dhcp).

## Boot

Its time to network boot your machines. Use the BMC's remote management capablities (may be vendor-specific) to set the boot device (on the next boot only) to PXE and power on each machine.

```sh
$ ipmitool -H node1.example.com -U USER -P PASS power off
$ ipmitool -H node1.example.com -U USER -P PASS chassis bootdev pxe
$ ipmitool -H node1.example.com -U USER -P PASS power on
```

Each machine should chainload iPXE, delegate to `matchbox`, receive its iPXE config (or other supported configs) and begin the provisioning process. The `simple-install` example assumes your machines are configured to boot from disk first and PXE only when requested, but you can write profiles for different cases.

Once the Container Linux install completes and the machine reboots you can SSH,

```ssh
$ ssh core@node1.example.com
```

To re-provision the machine for another purpose, run `terraform apply` and PXE boot it again.

## Going Further

Matchbox can be used to provision multi-node Container Linux clusters at one or many on-premise sites if deployed in an HA way. Machines can be matched individually by MAC address, UUID, region, or other labels you choose. Installs can be made much faster by caching images in the built-in HTTP [assets](api.md#assets) server.

[Container Linux configs](https://coreos.com/os/docs/latest/configuration.html) can be used to partition disks and filesystems, write systemd units, write networkd configs or regular files, and create users. Container Linux nodes can be provisioned into a system that meets your needs. Checkout the examples which create a 3 node [etcd](../examples/terraform/etcd3-install) cluster or a 3 node [Kubernetes](../examples/terraform/bootkube-install) cluster.

[terraform-dl]: https://www.terraform.io/downloads.html
