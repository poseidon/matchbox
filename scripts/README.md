# Scripts

## get-fedora-coreos

Run the `get-fedora-coreos` script to download Fedora CoreOS images, verify them, and move them into `examples/assets`.

```
./scripts/get-fedora-coreos
./scripts/get-fedora-coreos stream version dest
```

This will create:

```
examples/assets/fedora-coreos/
├── fedora-coreos-36.20220906.3.2-live-initramfs.x86_64.img
├── fedora-coreos-36.20220906.3.2-live-kernel-x86_64
├── fedora-coreos-36.20220906.3.2-live-rootfs.x86_64.img
```

## get-flatcar

Run the `get-flatcar` script to download Flatcar Linux images, verify them, and move them into `examples/assets`.

```
./scripts/get-flatcar
./scripts/get-flatcar channel version dest
```

This will create:

```
examples/assets/flatcar/
└── 3227.2.0
    ├── Flatcar_Image_Signing_Key.asc
    ├── flatcar_production_image.bin.bz2
    ├── flatcar_production_image.bin.bz2.sig
    ├── flatcar_production_pxe_image.cpio.gz
    ├── flatcar_production_pxe_image.cpio.gz.sig
    ├── flatcar_production_pxe.vmlinuz
    ├── flatcar_production_pxe.vmlinuz.sig
    └── version.txt
```

## libvirt

Create QEMU/KVM VMs which are configured to boot from the network. The `scripts/libvirt` script will create virtual machines on the `metal0` or `docker0` bridge with known hardware attributes (e.g. UUID, MAC address).

    $ sudo ./scripts/libvirt
    USAGE: libvirt <command>
    Commands:
        create      create QEMU/KVM nodes on the docker0 bridge
        start       start the QEMU/KVM nodes
        reboot      reboot the QEMU/KVM nodes
        shutdown    shutdown the QEMU/KVM nodes
        poweroff    poweroff the QEMU/KVM nodes
        destroy     destroy the QEMU/KVM nodes
