
# Scripts

## get-coreos

Run the `get-coreos` script to quickly download CoreOS kernel and initrd images, verify them, and move them into `assets`.

    ./scripts/get-coreos                 # beta, 899.6.0
    ./scripts/get-coreos alpha 942.0.0

This will create:

    assets/
    └── coreos
        └── 899.6.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz
        └── 942.0.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

## libvirt

Create libvirt VM nodes which are configured to boot from the network or from disk (empty). The `scripts/libvirt` script will create virtual machines on the `metal0` or `docker0` bridge with known hardware attributes (e.g. UUID, MAC address).

    $ sudo ./scripts/libvirt
    USAGE: libvirt <command>
    Commands:
        create-docker  create 4 libvirt nodes on the docker0 bridge
        create-rkt     create 4 libvirt nodes on a rkt CNI metal0 bridge
        start          start the 4 libvirt nodes
        reboot         reboot the 4 libvirt nodes
        shutdown       shutdown the 4 libvirt nodes
        poweroff       poweroff the 4 libvirt nodes
        destroy        destroy the 4 libvirt nodes
        delete-disks   delete the allocated disks

