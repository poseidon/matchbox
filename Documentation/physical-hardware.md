
# CoreOS on Baremetal/Physical Hardware

Physical or baremetal hardware can be booted and configured with CoreOS by PXE, iPXE, and Pixiecore network environments. This guide will show how to setup and test network boot environments with physical clients.

DHCP, TFTP, and HTTP services can run on separate network hosts or on the same host in PXE network environments. We'll focus on configuring a single *provisioner* Linux host to run these services as containers.

## Requirements

Client hardware boot firmware must support PXE or iPXE and the client must have at least one PXE-capable network interface. Most boot firmware and network cards support PXE so you're probably ok.

## Inspect

Identify whether the network runs a DHCP service which can be configured or whether you'll need to run a proxyDHCP service. Check a host's routing table to find the gateway where DHCP is likely to be running.

    route -n        # e.g. Gateway 192.168.1.1

## Config Service

Setup `coreos/bootcfg` according to the [docs](bootcfg.md). Pull the `coreos/bootcfg` image, prepare a data volume with `Machine` definitions, `Spec` definitions and cloud configs. Optionally, include a volume of downloaded image assets.

Run the `bootcfg` container to serve configs for any of the network environments we'll discuss next.

    docker run -p 8080:8080 --net=host --name=bootcfg --rm -v $PWD/data:/data:Z -v $PWD/images:/images:Z coreos/bootcfg:latest -address=0.0.0.0:8080 [-log-level=debug]

Note, the `Spec` examples in [data](../data) point `cloud-config-url` kernel options to 172.17.0.2 (the libvirt case). Your kernel cmdline option urls should point to where `bootcfg` runs via IP or DNS name.

## Network Setups

Your network may already have a configurable PXE or iPXE server, configurable DHCP, a DHCP server you cannot modify, or no DHCP server at all. We'll show how to setup each network environment to talk to `bootcfg`, depending on your circumstances.

Several setups use the [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) program and included [coreos/dnsmasq](../dockerfiles/dnsmasq) docker image to run a PXE-enabled DHCP server, proxy DHCP server, or TFTP server as needed. Build the Docker image if needed.

### Configurable iPXE

If your network environment already supports iPXE, edit your iPXE boot script to chainload from `bootcfg`

    # boot.ipxe
    chain http://192.168.1.100:8080/boot.ipxe

Substitute the name or IP and port where `bootcfg` runs.

### Configurable DHCP

If the DHCP server on your network is PXE-enabled and configurable, send the `bootcfg` iPXE endpoint as the boot filename option (e.g. `http://192.168.1.100:8080/boot.ipxe`). Substitute the name or IP and port where `bootcfg` runs.

Optionally, respond to older PXE client firmware with the location of the `undionly.kpxe` boot program on your TFTP server.

With `dnsmasq`, here is an example `dnsmask.conf`

    # dnsmasq.conf
    dhcp-range=192.168.1.1,192.168.1.254,30m
    enable-tftp
    tftp-root=/var/lib/tftpboot
    # set tag "ipxe" if request comes from iPXE ("iPXE" user class)
    dhcp-userclass=set:ipxe,iPXE
    # if PXE request came from regular firmware, serve iPXE firmware (via TFTP)
    dhcp-boot=tag:!ipxe,undionly.kpxe
    # if PXE request came from iPXE, serve an iPXE boot script (via HTTP)
    dhcp-boot=tag:ipxe,http://192.168.1.100:8080/boot.ipxe

### proxyDHCP

If the network already runs a DHCP service, setup a PXE/iPXE network environment alongside it with proxyDHCP and TFTP.

Run DHCP in proxy mode to respond to DHCP requests on the subnet. Optionally, serve the `undionly.pxe` boot file to older, non-iPXE clients (the '#' means not). Detect iPXE clients by the user class sent in their DHCPDISCOVER (or by Option 175) and point them to the `bootcfg` iPXE boot script.

```
sudo docker run --net=host --rm --cap-add=NET_ADMIN coreos/dnsmasq -d -q -i enp0s25 --dhcp-range=192.168.1.1,proxy,255.255.255.0 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --pxe-service=tag:#ipxe,x86PC,"PXE chainload to iPXE",undionly.kpxe --pxe-service=tag:ipxe,x86PC,"iPXE",http://192.168.1.100:8080/boot.ipxe
```

Change the `dhcp-range`, `-i interface`, and boot.ipxe endpoint to match your environment.

In this example, an existing router at 192.168.1.1 runs DHCP to allocate IP addresses between 192.168.1.2 and 192.168.1.254. The proxyDHCP is configured to respond to disocver requests on the 192.168.1.0/24 subnet. `bootcfg` runs on host (192.168.1.100) to serve iPXE boot scripts.

### DHCP

If the network does not already run a DHCP service, you can run one yourself and provide PXE options to baremetal clients. This is the case if your baremetal machines are connected to an isolated switch.

Identify a host machine which should run the DHCP service. If this machine has two NICs, it can serve as a router, using one for the uplink connection and the other to connect to the subnet with baremetal clients.

Run DHCP to allocate IP address leases and TFTP to serve the `undionly.pxe` boot file to older, non-iPXE clients (the '#' means not). Point iPXE clients to the `bootcfg` iPXE boot script.

```
sudo docker run --net=host --rm --cap-add=NET_ADMIN coreos/dnsmasq -d -q -i enp0s20u1 --dhcp-range=192.168.1.101,192.168.1.150 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://192.168.1.100:8080/boot.ipxe
```

Change the `dhcp-range`, `-i interface`, and boot.ipxe endpoint to match your environment.

In this example, a DHCP server is configured to allocate IP addresses between 192.168.1.101 and 192.168.1.150. `bootcfg` runs on host (192.168.1.100) to serve iPXE boot scripts.

You may have to explicitly assign the interface (-i) a network address and mark the interface as up.

    ip addr add 192.168.1.100/24 dev enp0s20u1
    ip link set dev enp0s20u1 up 

## Troubleshooting

**Firewall**: Running DHCP or proxyDHCP with `coreos/dnsmasq` on a host requires that the Firewall allow DHCP and TFTP (for chainloading) services to run.

**Port Collision** Running DHCP or proxyDHCP can cause port already in use collisions depending on what's running. Fedora runs bootp listening on udp/67 for example. Find the service using the port.

    sudo lsof -i :67

Evaluate whether you can configure the existing service or whether you'd like to stop it and test with `coreos/dnsmasq`.

### Alternatives

If you prefer, [Debian](http://www.debian-administration.org/article/478/Setting_up_a_server_for_PXE_network_booting), [Fedora](https://docs.fedoraproject.org/en-US/Fedora/7/html/Installation_Guide/ap-pxe-server.html), and [Ubuntu](https://help.ubuntu.com/community/DisklessUbuntuHowto) provide guides on PXE server setups. This project also includes [Vagrantfiles](vagrant) to quickly setup example Fedora PXE, iPXE, and Pixiecore servers on libvirt.
