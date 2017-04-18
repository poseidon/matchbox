# Must be run as root/sudoed

bootconfig: rkt --insecure-options=image run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/bootcfg-data --mount volume=assets,target=/var/lib/bootcfg/assets --volume assets,kind=host,source=$PWD/assets bootcfg.aci -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug

dns: rkt run coreos.com/dnsmasq:v0.2.0 --net=metal0:IP=172.15.0.3 -- -d -q --dhcp-range=172.15.0.50,172.15.0.99  --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.15.0.1 --address=/bootcfg.foo/172.15.0.2
