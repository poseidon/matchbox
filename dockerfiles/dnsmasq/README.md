
# coreos/dnsmasq

[coreos/dnsmasq](https://quay.io/repository/coreos/dnsmasq) is a convenience entrypoint to [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) for running DHCP, proxy DHCP, and TFTP without making changes to the host `/etc/dnsmasq.conf`.

The image bundles `undionly.kpxe` which chainloads PXE clients to iPXE.

## Usage

Build the image

    cd dockerfiles/dnsmasq
    ./docker-build

Run `dnsmasq` on a host in proxyDHCP mode to chainload iPXE.

    docker run --net=host --rm --cap-add=NET_ADMIN coreos/dnsmasq -d -i enp0s25 --dhcp-range=192.168.86.0,proxy,255.255.255.0 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --pxe-service=tag:#ipxe,x86PC,"PXE chainload to iPXE",undionly.kpxe




