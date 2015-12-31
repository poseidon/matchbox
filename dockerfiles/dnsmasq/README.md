
# coreos/dnsmasq

The coreos/dnsmasq Docker image provides an entrypoint to [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) for running dnsmasq DHCP, proxyDHCP, and TFTP. This is useful for testing different network setups without requiring changes to `dnsmasq.conf` on your host. The image also bundles `undionly.kpxe` which can be used to chainload PXE clients to iPXE.

## Usage

Build the image

    cd dockerfiles/dnsmasq
    ./docker-build

Run `dnsmasq` on a host in proxyDHCP mode to chainload iPXE.

    docker run --net=host --rm --cap-add=NET_ADMIN coreos/dnsmasq -d -i enp0s25 --dhcp-range=192.168.86.0,proxy,255.255.255.0 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --pxe-service=tag:#ipxe,x86PC,"PXE chainload to iPXE",undionly.kpxe




