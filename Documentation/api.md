
# API

## iPXE Script

Serves a static iPXE boot script which gathers client machine attributes and chain loads to the iPXE endpoint. Configure your DHCP server or iPXE server to boot from this script (e.g. set `dhcp-boot:dhcp-boot=tag:ipxe,http://bootcfg.domain.com/ipxe/boot.ipxe` if using `dnsmasq`).

    GET http://bootcfg.example.com/boot.ipxe
    GET http://bootcfg.example.com/boot.ipxe.0   // for dnsmasq

**Response**

    #!ipxe
    chain config?uuid=${uuid}&mac=${net0/mac:hexhyp}

## iPXE

Finds the spec matching the hardware attribute query parameters and renders the boot config as an iPXE script. Attributes are matched in priority order (UUID, MAC).

    GET http://bootcfg.example.com/ipxe

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    #!ipxe
    kernel /images/coreos/835.9.0/coreos_production_pxe.vmlinuz cloud-config-url=http://172.17.0.2:8080/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp} coreos.autologin
    initrd  /images/coreos/835.9.0/coreos_production_pxe_image.cpio.gz
    boot

The kernel, cmdline kernel options, and initrd are populated from a `Spec`.

## Pixiecore

Finds the spec matching the hardware attribute query parameters and renders the boot config as JSON to implement the Pixiecore API [spec](https://github.com/danderson/pixiecore/blob/master/README.api.md). Currently, Pixiecore only provides the machine's MAC address for matching specs.

    GET http://bootcfg.example.com/pixiecore/v1/boot/:MAC

**URL Parameters**

| Name | Type   | Description |
|------|--------|-------------|
| mac  | string | MAC address |

**Response**

    {
      "kernel":"/images/coreos/877.1.0/coreos_production_pxe.vmlinuz",
      "initrd":["/images/coreos/877.1.0/coreos_production_pxe_image.cpio.gz"],
      "cmdline":{
        "cloud-config-url":"http://bootcfg.example.com/cloud",
        "coreos.autologin":""
      }
    }

## Cloud Config

Finds the spec matching the hardware attribute query parameters and returns the specified cloud config file. Attributes are matched in priority order (UUID, MAC).

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

## API Resources

### Machines

Get a `Machine` definition by id (UUID, MAC).

    http://bootcfg.domain.com/machine/:id

**URL Parameters**

| Name | Type   | Description |
|------|--------|-------------|
| id   | string | machine identifier |

**Response**

```json
{
  "id": "2d9354a2-e8db-4021-bff5-20ffdf443d6f",
  "spec": {
    "id": "",
    "boot": {
      "kernel": "\/images\/coreos\/877.1.0\/coreos_production_pxe.vmlinuz",
      "initrd": [
        "\/images\/coreos\/877.1.0\/coreos_production_pxe_image.cpio.gz"
      ],
      "cmdline": {
        "cloud-config-url": "http:\/\/172.17.0.2:8080\/cloud?uuid=${uuid}&mac=${net0\/mac:hexhyp}",
        "coreos.autologin": ""
      }
    },
    "cloud_id": "node1-cloud.yml"
  },
  "spec_id": ""
}
```

### Specs

Get a `Spec` definition by id (UUID, MAC).

    http://bootcfg.domain.com/spec/:id

**URL Parameters**

| Name | Type   | Description |
|------|--------|-------------|
| id   | string | spec identifier |

**Response**

```json
{
  "id": "orion",
  "boot": {
    "kernel": "\/images\/coreos\/835.9.0\/coreos_production_pxe.vmlinuz",
    "initrd": [
      "\/images\/coreos\/835.9.0\/coreos_production_pxe_image.cpio.gz"
    ],
    "cmdline": {
      "cloud-config-url": "http:\/\/172.17.0.2:8080\/cloud?uuid=${uuid}&mac=${net0\/mac:hexhyp}",
      "coreos.autologin": ""
    }
  },
  "cloud_id": "orion-cloud-config.yml"
}
```

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

