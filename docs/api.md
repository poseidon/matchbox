
# bootcfg API

## iPXE Script

Serves a static iPXE boot script which gathers client machine attributes and chain loads to the iPXE boot config endpoint. Configure your iPXE server or DHCP options to boot from this script (e.g. set `dhcp-boot:dhcp-boot=tag:ipxe,http://provisioner.example.com/ipxe/boot.ipxe` if using `dnsmasq`).

    GET http://bootcfg.example.com/boot.ipxe

**Response**

    #!ipxe
    chain config?uuid=${uuid}

## iXPE Boot Config

Renders an iPXE boot config script based on hardware attribute query parameters and registered boot configs. Hardware attributes are matched in priority order: UUID, MAC address, default (if set).

    GET http://bootcfg.example.com/ipxe

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    #!ipxe
    kernel kernel-for-machine.vmlinuz cmd-line-options
    initrd initrd-for-machine
    boot

`kernel-for-machine`, `cmd-line-options`, and `initrf-for-machine` are populated from a boot config.

## Pixiecore Boot Config

Renders a Pixiecore boot config as JSON to implement the Pixiecore API [spec](https://github.com/danderson/pixiecore/blob/master/README.api.md) for Pixiecore setups. Currently, Pixiecore only provides MAC addresses for mapping to boot configs.

    GET http://bootcfg.example.com/pixiecore/v1/boot/:MAC

**URL Parameters**

| Name | Type   | Description |
|------|--------|-------------|
| mac  | string | MAC address |

**Response**

    {
      "kernel": "\/images\/coreos_production_pxe.vmlinuz",
      "initrd": [
        "\/images\/coreos_production_pxe_image.cpio.gz"
      ],
      "cmdline": {
        "coreos.autologin": ""
      }
    }

## Cloud Configs

Serves cloud configs based on hardware attribute query parameters and registered cloud configs. Hardware attributes are matched in priority order: UUID, MAC address, default (if set).

    GET http://bootcfg.example.com/cloud

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    #cloud-config
    coreos:
      units:
        - name: etcd2.service
          command: start
        - name: fleet.service
          command: start

## Image Assets

If you need to host kernel and initrd images within your network, bootcfg server's `/images/` route serves free-form static assets. Set the `-images-path` when starting the bootcfg server. Here is an example images directory layout:

    images/
    └── coreos
        └── 835.9.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz
        └── 877.1.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

