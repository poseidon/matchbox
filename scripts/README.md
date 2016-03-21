
# Scripts

## get-coreos

Run the `get-coreos` script to download CoreOS kernel and initrd images, verify them, and move them into `assets`.

    ./scripts/get-coreos
    ./scripts/get-coreos channel version

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
        create-docker   create libvirt nodes on the docker0 bridge
        create-rkt      create libvirt nodes on a rkt CNI metal0 bridge
        start           start the libvirt nodes
        reboot          reboot the libvirt nodes
        shutdown        shutdown the libvirt nodes
        poweroff        poweroff the libvirt nodes
        destroy         destroy the libvirt nodes
        delete-disks    delete the allocated disks

