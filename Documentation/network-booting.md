
# Network boot environments

This guide reviews network boot protocols and the different ways client machines can be PXE booted.

## PXE

The Preboot eXecution Environment (PXE) defines requirements for consistent, hardware-independent network-based machine booting and configuration. Formally, PXE specifies pre-boot protocol services that client NIC firmware must provide (DHCP, TFTP, UDP/IP), specifies boot firmware requirements, and defines a client-server protocol for obtaining a network boot program (NBP) which automates OS installation and configuration.

![PXE protocol](img/pxelinux.png)

At power-on, if a client machine's BIOS or UEFI boot firmware is set to perform network booting, the network interface card's PXE firmware broadcasts a DHCPDISCOVER packet identifying itself as a PXEClient to the network environment.

The network environment can be set up in a number of ways, which we'll discuss. In the simplest, a PXE-enabled DHCP Server responds with a DHCPOFFER with Options, which include a TFTP server IP ("next server") and the name of an NBP ("boot filename") to download (e.g. pxelinux.0). PXE firmware then downloads the NBP over TFTP and starts it. Finally, the NBP loads configs, scripts, and/or images it requires to run an OS.

### Network boot programs

Machines can be booted and configured with CoreOS using several network boot programs and approaches. Let's review them. If you're new to network booting or unsure which to choose, iPXE is a reasonable and flexible choice.

#### PXELINUX

[PXELINUX](http://www.syslinux.org/wiki/index.php/PXELINUX) is a common network boot program which loads a config file from `mybootdir/pxelinux.cfg/`  over TFTP. The file is chosen based on the client's UUID, MAC address, IP address, or a default.

```sh
$ mybootdir/pxelinux.cfg/b8945908-d6a6-41a9-611d-74a6ab80b83d
$ mybootdir/pxelinux.cfg/default
```

Here is an example PXE config file which boots a CoreOS image hosted on the TFTP server.

```
default coreos
prompt 1
timeout 15

display boot.msg

label coreos
  menu default
  kernel coreos_production_pxe.vmlinuz
  append initrd=coreos_production_pxe_image.cpio.gz cloud-config-url=http://example.com/pxe-cloud-config.yml
```

PXELINUX then downloads the specified kernel and init RAM filesystem images with TFTP.

This approach has a number of drawbacks. TFTP can be slow, managing config files can be tedious, and using different ignition or cloud configs on different machines requires separate pxelinux configs. These limitations spurred the development of various enhancements to PXE, discussed next.

#### iPXE

[iPXE](http://ipxe.org/) is an enhanced implementation of the PXE client firmware and a network boot program which uses iPXE scripts rather than config files and can download scripts and images with HTTP.

![iPXE flow](img/ipxe.png)

A DHCPOFFER to iPXE client firmware specifies an HTTP boot script such as `http://matchbox.foo/boot.ipxe`.

Here is an example iPXE script for booting the remote CoreOS stable image.

```
#!ipxe

set base-url http://stable.release.core-os.net/amd64-usr/current
kernel ${base-url}/coreos_production_pxe.vmlinuz cloud-config-url=http://provisioner.example.net/cloud-config.yml
initrd ${base-url}/coreos_production_pxe_image.cpio.gz
boot
```

A TFTP server is used only to provide the `undionly.kpxe` boot program to older PXE firmware in order to bootstrap into iPXE.

CoreOS `matchbox` can render signed iPXE scripts to machines based on their hardware attributes. Setup involves configuring your DHCP server to point iPXE clients to the `matchbox` [iPXE endpoint](api.md#ipxe).

## DHCP

Many networks have DHCP services which are impractical to modify or disable. Company DHCP servers are governed by network admin policies and home/office networks often have routers running a DHCP service which cannot supply PXE options to PXE clients.

To address this, PXE client firmware listens for DHCPOFFERs from a non-PXE DHCP server *and* a PXE-enabled **proxyDHCP server** configured to respond with the next server and boot filename only. Client firmware combines the two responses as if they had come from a single PXE-enabled DHCP server.

![Proxy DHCP flow](img/proxydhcp.png)
