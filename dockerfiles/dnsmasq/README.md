

Use this container image to run the `dnsmasq` command on a host, in an isolated container which includes iPXE's undionly.kpxe boot file.

## Usage

Build the image

    ./docker-build

Run `dnsmasq` on a host in proxyDHCP mode to chainload iPXE.

    docker run --net=host --rm --cap-add=NET_ADMIN coreos/dnsmasq -d -i enp0s25 --dhcp-range=192.168.86.0,proxy,255.255.255.0 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --pxe-service=tag:#ipxe,x86PC,"PXE chainload to iPXE",undionly.kpxe




