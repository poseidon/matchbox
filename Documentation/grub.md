
# GRUB2 Netboot

Use GRUB to network boot UEFI hardware.

## Requirements

For local development, install the dependencies for libvirt with UEFI.

* [UEFI with QEMU](https://fedoraproject.org/wiki/Using_UEFI_with_QEMU)

Ensure that you've gone through the [matchbox with rkt](getting-started-rkt.md) and [matchbox](matchbox.md) guides and understand the basics.

## Containers

Run `matchbox` with rkt, but mount the [grub](../examples/groups/grub) group example.

## Network

On Fedora, add the `metal0` interface to the trusted zone in your firewall configuration.

```sh
$ sudo firewall-cmd --add-interface=metal0 --zone=trusted
```

Run the `coreos.com/dnsmasq` ACI with rkt.

```sh
$ sudo rkt run coreos.com/dnsmasq:v0.3.0 --net=metal0:IP=172.18.0.3 -- -d -q --dhcp-range=172.18.0.50,172.18.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-match=set:efi-bc,option:client-arch,7 --dhcp-boot=tag:efi-bc,grub.efi --dhcp-userclass=set:grub,GRUB2 --dhcp-boot=tag:grub,"(http;matchbox.foo:8080)/grub","172.18.0.2" --log-queries --log-dhcp --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:pxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://matchbox.foo:8080/boot.ipxe --address=/matchbox.foo/172.18.0.2
```

## Client VM

Create UEFI VM nodes which have known hardware attributes.

```sh
$ sudo ./scripts/libvirt create-uefi
```

## Docker

If you use Docker, run `matchbox` according to [matchbox with Docker](getting-started-docker.md), but mount the [grub](../examples/groups/grub) group example. Then start the `coreos/dnsmasq` Docker image, which bundles a `grub.efi`.

```sh
$ sudo docker run --rm --cap-add=NET_ADMIN quay.io/coreos/dnsmasq -d -q --dhcp-range=172.17.0.43,172.17.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-match=set:efi-bc,option:client-arch,7 --dhcp-boot=tag:efi-bc,grub.efi --dhcp-userclass=set:grub,GRUB2 --dhcp-boot=tag:grub,"(http;matchbox.foo:8080)/grub","172.17.0.2" --log-queries --log-dhcp --dhcp-option=3,172.17.0.1 --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:pxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://matchbox.foo:8080/boot.ipxe --address=/matchbox.foo/172.17.0.2
```

Create a VM to verify the machine network boots.

```sh
$ sudo virt-install --name uefi-test --pxe --boot=uefi,network --disk pool=default,size=4 --network=bridge=docker0,model=e1000 --memory=1024 --vcpus=1 --os-type=linux --noautoconsole
```
