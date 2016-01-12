
# API

## iPXE Script

Serves a static iPXE boot script which gathers client machine attributes and chainloads to the iPXE endpoint. Configure your DHCP server or iPXE server to boot from this script (e.g. set `dhcp-boot:dhcp-boot=tag:ipxe,http://bootcfg.domain.com/ipxe/boot.ipxe` if using `dnsmasq`).

    GET http://bootcfg.example.com/boot.ipxe
    GET http://bootcfg.example.com/boot.ipxe.0   // for dnsmasq

**Response**

    #!ipxe
    chain ipxe?uuid=${uuid}&mac=${net0/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}

## iPXE

Finds the spec matching the attribute query parameters and renders the boot config as an iPXE script.

    GET http://bootcfg.example.com/ipxe

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    #!ipxe
    kernel /assets/coreos/835.9.0/coreos_production_pxe.vmlinuz cloud-config-url=http://172.17.0.2:8080/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp} coreos.autologin
    initrd  /assets/coreos/835.9.0/coreos_production_pxe_image.cpio.gz
    boot

The kernel, cmdline kernel options, and initrd are populated from a `Spec`.

## Pixiecore

Finds the spec matching the attribute query parameters and renders the boot config as JSON to implement the Pixiecore API [spec](https://github.com/danderson/pixiecore/blob/master/README.api.md). Currently, Pixiecore only provides the machine's MAC address for matching specs.

    GET http://bootcfg.example.com/pixiecore/v1/boot/:MAC

**URL Parameters**

| Name | Type   | Description |
|------|--------|-------------|
| mac  | string | MAC address |

**Response**

    {
      "kernel":"/assets/coreos/877.1.0/coreos_production_pxe.vmlinuz",
      "initrd":["/assets/coreos/877.1.0/coreos_production_pxe_image.cpio.gz"],
      "cmdline":{
        "cloud-config-url":"http://bootcfg.example.com/cloud",
        "coreos.autologin":""
      }
    }

## Cloud Config

Finds the spec matching the attribute query parameters and returns the corresponding cloud config file.

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

## Ignition Config

Finds the spec matching the attribute query parameters and returns the corresponding ignition config JSON.

    GET http://bootcfg.example.com/ignition

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    {
      "ignitionVersion": 1,
      "storage": {},
      "systemd": {
        "units": [
          {
            "name": "hello.service",
            "enable": true,
            "contents": "[Service]\nType=oneshot\nExecStart=\/usr\/bin\/echo Hello World\n\n[Install]\nWantedBy=multi-user.target"
          }
        ]
      },
      "networkd": {},
      "passwd": {}
    }


## API Resources

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
    "kernel": "\/assets\/coreos\/835.9.0\/coreos_production_pxe.vmlinuz",
    "initrd": [
      "\/assets\/coreos\/835.9.0\/coreos_production_pxe_image.cpio.gz"
    ],
    "cmdline": {
      "cloud-config-url": "http:\/\/172.17.0.2:8080\/cloud?uuid=${uuid}&mac=${net0\/mac:hexhyp}",
      "coreos.autologin": ""
    }
  },
  "cloud_id": "orion-cloud-config.yml"
}
```

## Assets

If you need to host static assets (e.g. kernel, initrd) within your network, bootcfg server's `/assets/` route serves free-form static assets. Set the `-assets-path` when starting the bootcfg server. Here is an example:

    assets/
    └── coreos
        └── 835.9.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz
        └── 877.1.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

