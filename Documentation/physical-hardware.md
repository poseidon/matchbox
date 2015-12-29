
# CoreOS on Baremetal/Physical Hardware

Physical or baremetal hardware can be booted and configured with CoreOS by PXE, iPXE, and Pixiecore network environments. This guide will show how to setup and test network boot environments with physical clients.

DHCP, TFTP, and HTTP services can run on separate network hosts or on the same host in PXE network environments. We'll focus on configuring a single *provisioner* Linux host to run these services as containers.

A [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) Docker image [is included](../dockerfiles/dnsmasq) for running dnsmasq DHCP, proxyDHCP, and TFTP for trying different network setups without requiring changes to `dnsmasq.conf` on your host. It also bundles `undionly.kpxe`  which is used for chainloading PXE clients to iPXE.

## Requirements

First, be sure that your client hardware has at least one PXE-capable network interface and that its boot firmware supports PXE. Most network cards and boot firmware today support PXE so you're probably ok.

Next, identify whether the network runs a DHCP service which can be configured or whether you'll need to run a proxyDHCP service. Check a host's routing table to find the gateway where DHCP is likely to be running.

    route -n        # e.g. Gateway 192.168.1.1

## Boot Config Service

Pull the `coreos/bootcfg` image, prepare machine configs, and download image assets. Check [boot config service](bootcfg.md) usage for details.

Run the `bootcfg` container to serve configs for any of the network environments we'll discuss next.

    docker run -p 8080:8080 --net=host --name=bootcfg --rm -v $PWD/data:/data:Z -v $PWD/images:/images:Z coreos/bootcfg:latest -address=0.0.0.0:8080

Check that `/boot.ipxe`, `/ipxe`, `/cloud`, and `/images` respond to requests.

Note, that the example `cloud-config-url`s in the [data](../data) directory assume the provisioner runs on 172.17.0.2 (the libvirt case). Update the address to point to the provisioner where `bootcfg` runs or assign a DNS name.

## PXE/iPXE

### Configurable DHCP

If the DHCP server on your network is PXE-enabled and configurable, return the boot filename http://$BOOTCFG_HOST:$BOOTCFG_PORT/boot.ipxe to iPXE client firmware. Substitute the hostname and port at which you run the `bootcfg` service. Respond to older PXE client firmware with the location of the `undionly.kpxe` boot program on your TFTP server.

These steps depend on the DHCP and TFTP servers you use.

### proxyDHCP

If the network already runs a DHCP service, you can setup a PXE/iPXE network environment alongside it with proxyDHCP and TFTP. Run DHCP in proxy mode to respond to DHCP requests on the subnet and serve the `undionly.pxe` boot file to older, non-iPXE clients (the '#' means not). Detect iPXE clients by the user class sent in their DHCPDISCOVER (or by Option 175) and point them to the `bootcfg` iPXE boot script.

```
sudo docker run --net=host --rm --cap-add=NET_ADMIN coreos/dnsmasq -d -q -i enp0s25 --dhcp-range=192.168.1.1,proxy,255.255.255.0 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --pxe-service=tag:#ipxe,x86PC,"PXE chainload to iPXE",undionly.kpxe --pxe-service=tag:ipxe,x86PC,"iPXE",http://192.168.1.100:8080/boot.ipxe
```

In this example, a router at 192.168.1.1 ran DHCP to allocate IP addresses between 192.168.1.2 and 192.168.1.254 so proxyDHCP was configured to respond to disocver requests on the 192.168.1.0/24 subnet. The provisioner (192.168.1.100) runs `bootcfg` to provide iPXE boot scripts. If you wish, assign the provisioner a DNS name.

### DHCP

If the network does not already run a DHCP service, you can run one yourself and provide PXE options to baremetal clients. This is the case if your baremetal machines are connected to an isolated switch.

Identify a host machine which should run the DHCP service. If this machine has two NICs, it can serve as a router, using one for the uplink connection and the other to connect to the subnet with baremetal clients.

Run DHCP to allocate IP address leases and TFTP to serve the `undionly.pxe` boot file to older, non-iPXE clients (the '#' means not). Point iPXE clients to the `bootcfg` iPXE boot script.

```
sudo docker run --net=host --rm --cap-add=NET_ADMIN coreos/dnsmasq -d -q -i enp0s20u1 --dhcp-range=192.168.1.101,192.168.1.150 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://192.168.1.100:8080/boot.ipxe
```

Like in the previous example, the boot file `http://192.168.1.100:8080/boot.ipxe` points to where `bootcfg` runs and you may choose to assign a DNS name for this endpoint.

Note that the `-i enp0s20u1` flag specifies the interface on which dnsmasq should listen. This should be the interface on the subnet with the baremetal clients you wish to boot.

You may have to explicitly assign the interface a network address and mark the interface as up.

    ip addr add 192.168.1.100/24 dev enp0s20u1
    ip link set dev enp0s20u1 up 

### Alternatives

If you prefer, [Debian](http://www.debian-administration.org/article/478/Setting_up_a_server_for_PXE_network_booting), [Fedora](https://docs.fedoraproject.org/en-US/Fedora/7/html/Installation_Guide/ap-pxe-server.html), and [Ubuntu](https://help.ubuntu.com/community/DisklessUbuntuHowto) provide guides on PXE server setups. This project also includes [Vagrantfiles](vagrant) to quickly setup example Fedora PXE, iPXE, and Pixiecore servers on libvirt.
