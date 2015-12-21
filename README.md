
# CoreOS on Baremetal

CoreOS on Baremetal contains guides for booting and configuring CoreOS clusters on virtual or physical hardware. It includes Dockerfiles and Vagrantfiles for setting up a network boot environment and the `bootcfg` HTTP service for providing configs to machines based on their attributes.

## Guides

[Getting Started](docs/getting-started.md)
[Boot Config Service](docs/bootcfg.md)
[Libvirt Guide](docs/virtual-hardware.md)
[Baremetal Guide](docs/physical-hardware.md)
[bootcfg Config](docs/config.md)
[bootcfg API](docs/api.md)

## Networking

Use PXE, iPXE, or Pixiecore with `bootcfg` to define kernel boot images, options, and cloud configs for machines within your network.

To get started in a libvirt development environment (under Linux), you can run the `bootcfg` container on the docker0 virtual bridge, alongside either the included `ipxe` conatainer or `danderson/pixiecore` and `dhcp` containers. Then you'll be able to boot libvirt VMs on the same bridge or baremetal PXE clients that are attatched to your development machine via a network adapter and added to the bridge. `docker0` defaults to subnet 172.17.0.0/16. See [clients](#clients).

List all your bridges with

    brctl show

and on docker 1.9+, you may inspect the bridge and its interfaces with

    docker network inspect bridge   # docker calls the bridge "bridge"


To get started on a baremetal, run the `bootcfg` container with the `--net=host` argument and configure PXE, iPXE, or Pixiecore. If you do not yet have any of these services running on your network, the [dnsmasq image](dockerfiles/dnsmasq) can be helpful.

### iPXE

Configure your iPXE server to chainload the `bootcfg` iPXE boot script hosted at `$BOOTCFG_BASE_URL/boot.ipxe`. The base URL should be of the for `protocol://hostname:port` where where hostname is the IP address of the config server or a DNS name for it.

    # dnsmasq.conf
    # if request from iPXE, serve bootcfg's iPXE boot script (via HTTP)
    dhcp-boot=tag:ipxe,http://172.17.0.2:8080/boot.ipxe

#### docker0

Try running a PXE/iPXE server alongside `bootcfg` on the docker0 bridge. Use the included `ipxe` container to run an example PXE/iPXE server on that bridge.

The `ipxe` Docker image uses dnsmasq DHCP to point PXE/iPXE clients to the boot config service (hardcoded to http://172.17.0.2:8080). It also runs a TFTP server to serve the iPXE firmware to older PXE clients via chainloading.

    cd dockerfiles/ipxe
    ./docker-build
    ./docker-run

Now create local PXE boot [clients](#clients) as libvirt VMs or by attaching bare metal machines to the docker0 bridge.

### PXE

To use `bootcfg` with PXE, you must [chainload iPXE](http://ipxe.org/howto/chainloading). This can be done by configuring your PXE/DHCP/TFTP server to send `undionly.kpxe` over TFTP to PXE clients.

    # dnsmasq.conf
    enable-tftp
    tftp-root=/var/lib/tftpboot
    # if regular PXE firmware, serve iPXE firmware (via TFTP)
    dhcp-boot=tag:!ipxe,undionly.kpxe

`bootcfg` does not respond to DHCP requests or serve files over TFTP.

### Pixecore

Pixiecore is a ProxyDHCP, TFTP, and HTTP server and calls through to the `bootcfg` API to get a boot config for `pxelinux` to boot. No modification of your existing DHCP server is required in production.

#### docker0

Try running a DHCP server, Pixiecore, and `bootcfg` on the docker0 bridge. Use the included `dhcp` container to run an example DHCP server and the official Pixiecore container image.

    # DHCP
    cd dockerfiles/dhcp
    ./docker-build
    ./docker-run

Start Pixiecore using the script which attempts to detect the IP and port of `bootcfg` on the Docker host or do it manually.

    # Pixiecore
    ./scripts/pixiecore
    # manual
    docker run -v $PWD/images:/images:Z danderson/pixiecore -api http://$BOOTCFG_HOST:$BOOTCFG_PORT/pixiecore

Now create local PXE boot [clients](#clients) as libvirt VMs or by attaching bare metal machines to the docker0 bridge.

## Clients

Once boot services are running, create a PXE boot client VM or attach a bare metal machine to your host.

### VM

Create a VM using the virt-manager UI, select Network Boot with PXE, and for the network selection, choose "Specify Shared Device" with bridge name `docker0`. The VM should PXE boot using the boot configuration determined by its MAC address, which can be tweaked in virt-manager.

### Bare Metal

Link a bare metal machine, which has boot firmware (BIOS) support for PXE, to your host with a network adapter. Get the link and attach it to the bridge.

    ip link show                      # find new link e.g. enp0s20u2
    brctl addif docker0 enp0s20u2

Configure the boot firmware to prefer PXE booting or network booting and restart the machine. It should PXE boot using the boot configuration determined by its MAC address.
