
# GRUB2 Netboot

GRUB netboot support is experimental.

## Requirements

For local development, install the dependencies for libvirt with UEFI.

* [UEFI with QEMU](https://fedoraproject.org/wiki/Using_UEFI_with_QEMU)

## Application Container

Run the `bootcfg` ACI with rkt according to the [development docs](develop.md). Examples contains a [grub.yaml](../../examples/grub.yaml) config with a default machine group for GRUB net booting.

## Client VM

Create a VM with an e1000 or virtio network device.

    sudo virt-install --name uefi-test --pxe --boot=uefi,network --disk pool=default,size=4 --network=bridge=metal0,model=e1000 --memory=1024 --vcpus=1 --os-type=linux --noautoconsole

## Network

On Fedora, add the `metal0` interface to the trusted zone in your firewall configuration.

    sudo firewall-cmd --add-interface=metal0 --zone=trusted

Add a `grub.efi` to `tftpboot` before building the dnsmasq ACI.

    cd contrib/dnsmasq
    ./get-tftp-files
    sudo ./build-aci

Build dnsmasq ACI with `acbuild` and run with rkt.

    sudo rkt --insecure-options=image run dnsmasq.aci --net=metal0:IP=172.15.0.3 -- -d -q --dhcp-range=172.15.0.50,172.15.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-match=set:efi-bc,option:client-arch,7 --dhcp-boot=tag:efi-bc,grub.efi --dhcp-userclass=set:grub,GRUB2 --dhcp-boot=tag:grub,"(http;bootcfg.foo:8080)/grub","172.15.0.2","" --log-queries --log-dhcp --address=/bootcfg.foo/172.15.0.2 --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:pxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe

