# Getting started

In this tutorial, we'll use `matchbox` with Terraform to provision Fedora CoreOS or Flatcar Linux machines.

We'll install the `matchbox` service, setup a PXE network boot environment, and use Terraform configs to declare infrastructure and apply resources on `matchbox`.

## matchbox

Install `matchbox` on a host server or Kubernetes cluster. Generate TLS credentials and enable the gRPC API as directed. Save the `ca.crt`, `client.crt`, and `client.key` on your local machine (e.g. `~/.matchbox`).

* Installing on a [Linux distro](deployment.md)
* Installing on [Kubernetes](deployment.md#kubernetes)
* Running with [docker](deployment.md#docker)

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

Install [Terraform](https://www.terraform.io/downloads.html) v0.13+ on your system.

```sh
$ terraform version
Terraform v0.13.3
```

### Examples

Clone the matchbox source.

```sh
$ git clone https://github.com/poseidon/matchbox.git
$ cd matchbox/examples/terraform
```

Select from the Terraform [examples](https://github.com/poseidon/matchbox/tree/master/examples/terraform). For example,

* `fedora-coreos-install` - PXE boot, install Fedora CoreOS to disk, reboot, and machines come up with your SSH authorized key set
* `flatcar-install` - PXE boot, install Flatcar Linux to disk, reboot, and machines come up with your SSH authorized key set

These aren't exactly full clusters, but they show declarations and network provisioning.

```sh
$ cd fedora-coreos-install    # or flatcar-install
```

!!! note
    Fedora CoreOS images are only served via HTTPS, so your iPXE firmware must be compiled to support HTTPS downloads.

Let's review the terraform config and learn a bit about Matchbox.

### Provider

Matchbox is configured as a provider platform for bare-metal resources.

```tf
// Configure the matchbox provider
provider "matchbox" {
  endpoint    = var.matchbox_rpc_endpoint
  client_cert = file("~/.matchbox/client.crt")
  client_key  = file("~/.matchbox/client.key")
  ca          = file("~/.matchbox/ca.crt")
}

terraform {
  required_providers {
    ct = {
      source  = "poseidon/ct"
      version = "0.7.1"
    }
    matchbox = {
      source = "poseidon/matchbox"
      version = "0.4.1"
    }
  }
}
```

### Profiles

Machine profiles specify the kernel, initrd, kernel args, Ignition Config, and other configs (e.g. templated Container Linux Config, Cloud-config, generic) used to network boot and provision a bare-metal machine. The profile below would PXE boot machines using a Fedora CoreOS kernel and initrd (see [assets](api-http.md#assets) to learn about caching for speed), perform a disk install, reboot (first boot from disk), and use a [Fedora CoreOS Config](https://github.com/coreos/fcct/blob/master/docs/configuration-v1_1.md) to generate an Ignition config to provision.

```tf
// Fedora CoreOS profile
resource "matchbox_profile" "fedora-coreos-install" {
  name  = "worker"
  kernel = "https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-kernel-x86_64"
  initrd = [
    "https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-initramfs.x86_64.img"
  ]

  args = [
    "coreos.live.rootfs_url=https://builds.coreos.fedoraproject.org/prod/streams/${var.os_stream}/builds/${var.os_version}/x86_64/fedora-coreos-${var.os_version}-live-rootfs.x86_64.img",
    "coreos.inst.install_dev=/dev/sda",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "console=tty0",
    "console=ttyS0",
  ]

  raw_ignition = data.ct_config.worker-ignition.rendered
}

data "ct_config" "worker-ignition" {
  content  = data.template_file.worker-config.rendered
  strict   = true
}

data "template_file" "worker-config" {
  template = file("fcc/fedora-coreos.yaml")
  vars = {
    ssh_authorized_key     = var.ssh_authorized_key
  }
}
```

### Groups

Matcher groups match machines based on labels like MAC, UUID, etc. to different profiles and templates in machine-specific values. The group below does not have a `selector` block, so any machines which network boot from Matchbox will match this group and be provisioned using the `fedora-coreos-install` profile. Machines are matched to the most specific matching group.

```tf
// Default matcher group for machines
resource "matchbox_group" "default" {
  name    = "default"
  profile = matchbox_profile.fedora-coreos-install.name
}
```

### Variables

Some Terraform [variables](https://www.terraform.io/docs/configuration/variables.html) are used in the examples. A quick way to set their value is by creating a `terraform.tfvars` file.

```
cp terraform.tfvars.example terraform.tfvars
```

```tf
matchbox_http_endpoint = "http://matchbox.example.com:8080"
matchbox_rpc_endpoint = "matchbox.example.com:8081"
ssh_authorized_key = "YOUR_SSH_KEY"
```

### Apply

Initialize the Terraform workspace. Then plan and apply the resources.

```
terraform init
```

```
$ terraform apply
Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
```

Matchbox serves configs to machines and respects query parameters, if you're interested:

* iPXE default - [/ipxe](http://matchbox.example.com:8080/ipxe)
* Ignition default - [/ignition](http://matchbox.example.com:8080/ignition)
* Ignition post-install - [/ignition?os=installed](http://matchbox.example.com:8080/ignition?os=installed)

## Network

Matchbox can integrate with many on-premise network setups. It does not seek to be the DHCP server, TFTP server, or DNS server for the network. Instead, matchbox serves iPXE scripts as the entrypoint for provisioning network booted machines. PXE clients are supported by chainloading iPXE firmware.

In the simplest case, an iPXE-enabled network can chain to Matchbox,

```
# /var/www/html/ipxe/default.ipxe
chain http://matchbox.foo:8080/boot.ipxe
```

Read [network-setup.md](network-setup.md) for the complete range of options. Network admins have a great amount of flexibility:

* May keep using existing DHCP, TFTP, and DNS services
* May configure subnets, architectures, or specific machines to delegate to matchbox
* May place matchbox behind a menu entry (timeout and default to matchbox)

If you've never setup a PXE-enabled network before or you're trying to setup a home lab, checkout the [quay.io/poseidon/dnsmasq](https://quay.io/repository/poseidon/dnsmasq) container image [copy-paste examples](https://github.com/poseidon/matchbox/blob/master/docs/network-setup.md#poseidondnsmasq) and see the section about [proxy-DHCP](https://github.com/poseidon/matchbox/blob/master/docs/network-setup.md#proxy-dhcp).

## Boot

Its time to network boot your machines. Use the BMC's remote management capabilities (may be vendor-specific) to set the boot device (on the next boot only) to PXE and power on each machine.

```sh
$ ipmitool -H node1.example.com -U USER -P PASS power off
$ ipmitool -H node1.example.com -U USER -P PASS chassis bootdev pxe
$ ipmitool -H node1.example.com -U USER -P PASS power on
```

Each machine should chainload iPXE, delegate to Matchbox, receive its iPXE config (or other supported configs) and begin the provisioning process. The examples assume machines are configured to boot from disk first and PXE only when requested, but you can write profiles for different cases.

Once the install completes and the machine reboots, you can SSH.

```ssh
$ ssh core@node1.example.com
```

To re-provision the machine for another purpose, run `terraform apply` and PXE boot machines again.

## Going Further

Matchbox can be used to provision multi-node Fedora CoreOS or Flatcar Linux clusters at one or many on-premise sites if deployed in an HA way. Machines can be matched individually by MAC address, UUID, region, or other labels you choose. Installs can be made much faster by caching images in the built-in HTTP [assets](api-http.md#assets) server.

[Ignition](https://github.com/coreos/ignition) can be used to partition disks, create file systems, write systemd units, write networkd configs or regular files, and create users. Nodes can be network provisioned into a complete cluster system that meets your needs. For example, see [Typhoon](https://typhoon.psdn.io/fedora-coreos/bare-metal/).

