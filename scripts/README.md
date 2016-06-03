
# Scripts

## get-coreos

Run the `get-coreos` script to download CoreOS images, verify them, and move them into `examples/assets`.

    ./scripts/get-coreos
    ./scripts/get-coreos channel version

This will create:

    examples/assets/
    └── coreos
        └── 1053.2.0
            ├── CoreOS_Image_Signing_Key.asc
            ├── coreos_production_image.bin.bz2
            ├── coreos_production_image.bin.bz2.sig
            ├── coreos_production_pxe_image.cpio.gz
            ├── coreos_production_pxe_image.cpio.gz.sig
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe.vmlinuz.sig

## libvirt

Create libvirt VM nodes which are configured to boot from the network. The `scripts/libvirt` script will create virtual machines on the `metal0` or `docker0` bridge with known hardware attributes (e.g. UUID, MAC address).

    $ sudo ./scripts/libvirt
    USAGE: libvirt <command>
    Commands:
        create-docker   create libvirt nodes on the docker0 bridge
        create-rkt  create libvirt nodes on a rkt CNI metal0 bridge
        create-uefi create UEFI libvirt nodes on the rkt CNI metal0 bridge
        start       start the libvirt nodes
        reboot      reboot the libvirt nodes
        shutdown    shutdown the libvirt nodes
        poweroff    poweroff the libvirt nodes
        destroy     destroy the libvirt nodes

## k8s-certgen

Generate TLS certificates needed for a multi-node Kubernetes cluster. See the [examples](../examples/README.md#assets).

    $ ./scripts/tls/k8s-certgen -h
    ./scripts/tls/k8s-certgen -h
    Usage: k8s-certgen
    Options:
      -d DEST     Destination for generated files (default: .examples/assets/tls)
      -s SERVER   Reachable Server IP for kubeconfig (e.g. 172.15.0.21)
      -m MASTERS  Master Node Names/Addresses in SAN format (e.g. IP.1=10.3.0.1,IP.2=172.15.0.21).
      -w WORKERS  Worker Node Names/Addresses in SAN format (e.g. IP.1=172.15.0.22,IP.2=172.15.0.23)
      -h          Show help.
