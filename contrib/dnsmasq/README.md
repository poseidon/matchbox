# dnsmasq [![Docker Repository on Quay](https://quay.io/repository/coreos/dnsmasq/status "Docker Repository on Quay")](https://quay.io/repository/coreos/dnsmasq)

`dnsmasq` provides a container image for running DHCP, proxy DHCP, DNS, and/or TFTP with [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html). Use it to test different network setups with clusters of network bootable machines.

The image bundles `undionly.kpxe`, `ipxe.efi`, and `grub.efi` (experimental) for chainloading BIOS and UEFI clients to iPXE.

## Usage

Run the container image as a DHCP, DNS, and TFTP service.

```sh
sudo rkt run --net=host quay.io/coreos/dnsmasq \
  --caps-retain=CAP_NET_ADMIN,CAP_NET_BIND_SERVICE,CAP_SETGID,CAP_SETUID,CAP_NET_RAW \
  -- -d -q \
  --dhcp-range=192.168.1.3,192.168.1.254 \
  --enable-tftp \
  --tftp-root=/var/lib/tftpboot \
  --dhcp-match=set:bios,option:client-arch,0 \
  --dhcp-boot=tag:bios,undionly.kpxe \
  --dhcp-match=set:efi32,option:client-arch,6 \
  --dhcp-boot=tag:efi32,ipxe.efi \
  --dhcp-match=set:efibc,option:client-arch,7 \
  --dhcp-boot=tag:efibc,ipxe.efi \
  --dhcp-match=set:efi64,option:client-arch,9 \
  --dhcp-boot=tag:efi64,ipxe.efi \
  --dhcp-userclass=set:ipxe,iPXE \
  --dhcp-boot=tag:ipxe,http://matchbox.example.com:8080/boot.ipxe \
  --address=/matchbox.example.com/192.168.1.2 \
  --log-queries \
  --log-dhcp
```

```sh
sudo docker run --rm --cap-add=NET_ADMIN --net=host quay.io/coreos/dnsmasq \
  -d -q \
  --dhcp-range=192.168.1.3,192.168.1.254 \
  --enable-tftp --tftp-root=/var/lib/tftpboot \
  --dhcp-match=set:bios,option:client-arch,0 \
  --dhcp-boot=tag:bios,undionly.kpxe \
  --dhcp-match=set:efi32,option:client-arch,6 \
  --dhcp-boot=tag:efi32,ipxe.efi \
  --dhcp-match=set:efibc,option:client-arch,7 \
  --dhcp-boot=tag:efibc,ipxe.efi \
  --dhcp-match=set:efi64,option:client-arch,9 \
  --dhcp-boot=tag:efi64,ipxe.efi \
  --dhcp-userclass=set:ipxe,iPXE \
  --dhcp-boot=tag:ipxe,http://matchbox.example.com:8080/boot.ipxe \
  --address=/matchbox.example/192.168.1.2 \
  --log-queries \
  --log-dhcp
```

Press ^] three times to stop the rkt pod. Press ctrl-C to stop the Docker container.

## Configuration Flags

Configuration arguments can be provided as flags. Check the dnsmasq [man pages](http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html) for a complete list.

| flag     | description | example |
|----------|-------------|---------|
| --dhcp-range | Enable DHCP, lease given range | `172.18.0.50,172.18.0.99`, `192.168.1.1,proxy,255.255.255.0` |
| --dhcp-boot | DHCP next server option | `http://matchbox.foo:8080/boot.ipxe` |
| --enable-tftp | Enable serving from tftp-root over TFTP | NA |
| --address | IP address for a domain name | /matchbox.foo/172.18.0.2 |

## Development

Build a container image locally.

```
make docker-image
```

Run the image with Docker on the `docker0` bridge (default).

```
sudo docker run --rm --cap-add=NET_ADMIN coreos/dnsmasq -d -q
```

