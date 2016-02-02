
# rkt Tutorial

**Like the Docker libvirt setup, the rkt libvirt setup is meant for local development and testing**

Get started on your laptop with `rkt`, `libvirt`, and `virt-manager` (tested on Fedora 23).

Install a [rkt](https://github.com/coreos/rkt/releases) and [acbuild](https://github.com/appc/acbuild/releases) release and get the `libvirt` and `virt-manager` packages.

    sudo dnf install virt-manager virt-install

Clone the source.

    git clone https://github.com/coreos/coreos-baremetal.git
    cd coreos-baremetal

Currently, rkt and the Fedora/RHEL/CentOS SELinux policies aren't supported. See the [issue](https://github.com/coreos/rkt/issues/1727) tracking the work and policy changes. To test these examples on your laptop, set SELinux enforcement to permissive if you are comfortable (`sudo setenforce 0`). Enable it again when you are finished.

Download the CoreOS network boot images to assets/coreos:

    cd coreos-baremetal
    ./scripts/get-coreos

Define the `metal0` virtual bridge with [CNI](https://github.com/appc/cni).

    cat > /etc/rkt/net.d/20-metal.conf << EOF
    {
      "name": "metal0",
      "type": "bridge",
      "bridge": "metal0",
      "isGateway": true,
      "ipam": {
        "type": "host-local",
        "subnet": "172.15.0.0/16",
        "routes" : [ { "dst" : "172.15.0.0/16" } ]
       }
    }
    EOF

Run the config server on `metal0` with the IP address corresponding to the examples (or add DNS).

    sudo rkt --insecure-options=image fetch docker://quay.io/coreos/bootcfg

The insecure flag is needed because Docker images do not support signature verification.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/assets --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/data --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg -- -address=0.0.0.0:8080 -log-level=debug -config /data/etcd-rkt.yaml

If you get an error about the IP being assigned already.

    sudo rkt gc --grace-period=0
    sudo rkt list            # should be empty

Create 5 VM nodes on the `metal0` bridge, which have known "hardware" attributes that match the examples.

    sudo ./scripts/libvirt create-rkt
    # if you previously tried the docker examples, cleanup first
    sudo ./scripts/libvirt shutdown
    sudo ./scripts/libvirt destroy

In your firewall settings, configure the `metal0` interface as trusted.

Build an dnsmasq ACI and run it to create a DNS server, TFTP server, and DHCP server which points network boot clients to the config server started above.

    cd contrib/dnsmasq
    sudo ./build-aci

Run `dnsmasq.aci` to create a DHCP and TFTP server pointing to config server.

    sudo rkt --insecure-options=image run dnsmasq.aci --net=metal0 -- -d -q --dhcp-range=172.15.0.50,172.15.0.99 --enable-tftp --tftp-root=/var/lib/tftpboot --dhcp-userclass=set:ipxe,iPXE --dhcp-boot=tag:#ipxe,undionly.kpxe --dhcp-boot=tag:ipxe,http://bootcfg.foo:8080/boot.ipxe --log-queries --log-dhcp --dhcp-option=3,172.15.0.1 --address=/bootcfg.foo/172.15.0.2

Reboot Nodes

    sudo ./script/libvirt poweroff
    sudo ./script/libvirt start