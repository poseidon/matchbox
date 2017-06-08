
# Scripts

## get-coreos

Run the `get-coreos` script to download CoreOS Container Linux images, verify them, and move them into `examples/assets`.

    ./scripts/get-coreos
    ./scripts/get-coreos channel version

This will create:

    examples/assets/
    └── coreos
        └── 1153.0.0
            ├── CoreOS_Image_Signing_Key.asc
            ├── coreos_production_image.bin.bz2
            ├── coreos_production_image.bin.bz2.sig
            ├── coreos_production_pxe_image.cpio.gz
            ├── coreos_production_pxe_image.cpio.gz.sig
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe.vmlinuz.sig

## libvirt

Create QEMU/KVM VMs which are configured to boot from the network. The `scripts/libvirt` script will create virtual machines on the `metal0` or `docker0` bridge with known hardware attributes (e.g. UUID, MAC address).

    $ sudo ./scripts/libvirt
    USAGE: libvirt <command>
    Commands:
        create      create QEMU/KVM nodes on a rkt CNI metal0 bridge
        create-rkt  create QEMU/KVM nodes on a rkt CNI metal0 bridge
        create-docker   create QEMU/KVM nodes on the docker0 bridge
        create-uefi create UEFI QEMU/KVM nodes on the rkt CNI metal0 bridge
        start       start the QEMU/KVM nodes
        reboot      reboot the QEMU/KVM nodes
        shutdown    shutdown the QEMU/KVM nodes
        poweroff    poweroff the QEMU/KVM nodes
        destroy     destroy the QEMU/KVM nodes

## k8s-certgen

Generate TLS certificates needed for a multi-node Kubernetes cluster. See the [examples](../examples/README.md#assets).

    $ ./scripts/tls/k8s-certgen -h
    Usage: k8s-certgen
    Options:
      -d DEST     Destination for generated files (default: .examples/assets/tls)
      -s SERVER   Reachable Server IP for kubeconfig (e.g. node1.example.com)
      -m MASTERS  Controller Node Names/Addresses in SAN format (e.g. IP.1=10.3.0.1,DNS.1=node1.example.com)
      -w WORKERS  Worker Node Names/Addresses in SAN format (e.g. DNS.1=node2.example.com,DNS.2=node3.example.com)
      -h          Show help

