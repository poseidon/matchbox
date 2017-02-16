
# dnsmasq

[![Docker Repository on Quay](https://quay.io/repository/coreos/dnsmasq/status "Docker Repository on Quay")](https://quay.io/repository/coreos/dnsmasq)

`dnsmasq` provides an App Container Image (ACI) or Docker image for running DHCP, proxy DHCP, DNS, and/or TFTP with [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) in a container/pod. Use it to test different network setups with clusters of network bootable machines.

The image bundles `undionly.kpxe` which chainloads PXE clients to iPXE and `grub.efi` (experimental) which chainloads UEFI architectures to GRUB2.

## Usage

Run the `coreos.com/dnsmasq` ACI with rkt.

    sudo rkt trust --prefix coreos.com/dnsmasq
    # gpg key fingerprint is: 18AD 5014 C99E F7E3 BA5F  6CE9 50BD D3E0 FC8A 365E
    sudo rkt run coreos.com/dnsmasq:v0.3.0

Press ^] three times to kill the container.

Alternately, Docker can be used.

    docker pull quay.io/coreos/dnsmasq
    docker run --cap-add NET_ADMIN quay.io/coreos/dnsmasq

## Configuration Flags

Configuration arguments can be provided at the command line. Check the dnsmasq [man pages](http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html) for a complete list, but here are some important flags.

| flag     | description | example |
|----------|-------------|---------|
| --dhcp-range | Enable DHCP, lease given range | `172.18.0.50,172.18.0.99`, `192.168.1.1,proxy,255.255.255.0` |
| --dhcp-boot | DHCP next server option | `http://matchbox.foo:8080/boot.ipxe` |
| --enable-tftp | Enable serving from tftp-root over TFTP | NA |
| --address | IP address for a domain name | /matchbox.foo/172.18.0.2 |

## ACI

Build a `dnsmasq` ACI with the build script which uses [acbuild](https://github.com/appc/acbuild).

    cd contrib/dnsmasq
    ./get-tftp-files
    sudo ./build-aci

Run `dnsmasq.aci` with rkt to run DHCP/proxyDHCP/TFTP/DNS services.

DHCP+TFTP+DNS on the `metal0` bridge:

    sudo rkt --insecure-options=image run dnsmasq.aci --net=metal0 -- -d -q --dhcp-range=172.18.0.50,172.18.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://matchbox.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.18.0.1 --address=/matchbox.foo/172.18.0.2

## Docker

Build a Docker image locally using the tag `latest`.

    cd contrib/dnsmasq
    ./get-tftp-files
    sudo ./build-docker

Run the Docker image to run DHCP/proxyDHCP/TFTP/DNS services.

DHCP+TFTP+DNS on the `docker0` bridge:

    sudo docker run --rm --cap-add=NET_ADMIN quay.io/coreos/dnsmasq -d -q --dhcp-range=172.17.0.43,172.17.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://matchbox.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.17.0.1 --address=/matchbox.foo/172.17.0.2
