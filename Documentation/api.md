
# API

## iPXE Script

Serves a static iPXE boot script which gathers client machine attributes and chainloads to the iPXE endpoint. Use DHCP/TFTP to point iPXE clients to this endpoint as the next-server.

    GET http://bootcfg.foo/boot.ipxe
    GET http://bootcfg.foo/boot.ipxe.0   // for dnsmasq

**Response**

    #!ipxe
    chain ipxe?uuid=${uuid}&mac=${net0/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}

## iPXE

Finds the profile for the machine and renders the network boot config (kernel, options, initrd) as an iPXE script.

    GET http://bootcfg.foo/ipxe

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    #!ipxe
    kernel /assets/coreos/1032.0.0/coreos_production_pxe.vmlinuz coreos.config.url=http://bootcfg.foo:8080/ignition?uuid=${uuid}&mac=${net0/mac:hexhyp} coreos.first_boot=1 coreos.autologin
    initrd  /assets/coreos/1032.0.0/coreos_production_pxe_image.cpio.gz
    boot

## GRUB2

Finds the profile for the machine and renders the network boot config as a GRUB config. Use DHCP/TFTP to point GRUB clients to this endpoint as the next-server.

    GET http://bootcfg.foo/grub

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    default=0
    timeout=1
    menuentry "CoreOS" {
    echo "Loading kernel"
    linuxefi "(http;bootcfg.foo:8080)/assets/coreos/1032.0.0/coreos_production_pxe.vmlinuz" "coreos.autologin" "coreos.config.url=http://bootcfg.foo:8080/ignition" "coreos.first_boot"
    echo "Loading initrd"
    initrdefi "(http;bootcfg.foo:8080)/assets/coreos/1032.0.0/coreos_production_pxe_image.cpio.gz"
    }

## Pixiecore

Finds the profile matching the machine and renders the network boot config as JSON to implement the [Pixiecore API](https://github.com/danderson/pixiecore/blob/master/README.api.md). Currently, Pixiecore only provides the machine's MAC address for matching.

    GET http://bootcfg.foo/pixiecore/v1/boot/:MAC

**URL Parameters**

| Name | Type   | Description |
|------|--------|-------------|
| mac  | string | MAC address |

**Response**

    {
      "kernel":"/assets/coreos/1032.0.0/coreos_production_pxe.vmlinuz",
      "initrd":["/assets/coreos/1032.0.0/coreos_production_pxe_image.cpio.gz"],
      "cmdline":{
        "cloud-config-url":"http://bootcfg.foo/cloud?mac=ADDRESS",
        "coreos.autologin":""
      }
    }

## Cloud Config

Finds the profile matching the machine and renders the corresponding Cloud-Config with metadata.

    GET http://bootcfg.foo/cloud

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

Finds the profile matching the machine and renders the corresponding Ignition Config with metadata.

    GET http://bootcfg.foo/ignition

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

## Metadata

Finds the matching machine group and renders the selectors and metadata as a `plain/text` file.

    GET http://bootcfg.foo/metadata

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    IPV4_ADDRESS=172.15.0.21
    NETWORKD_ADDRESS=172.15.0.21/16
    NETWORKD_GATEWAY=172.15.0.1
    NETWORKD_NAME=ens3
    ETCD_NAME=node1
    FLEET_METADATA=role=etcd,name=node1
    UUID=16e7d8a7-bfa9-428b-9117-363341bb330b
    ETCD_INITIAL_CLUSTER=node1=http://172.15.0.21:2380,node2=http://172.15.0.22:2380,node3=http://172.15.0.23:2380
    NETWORKD_DNS=172.15.0.3

## OpenPGP Signatures

OpenPGPG signature endpoints serve detached binary and ASCII armored signatures of rendered configs, if enabled. See [OpenPGP Signing](openpgp.md).

| Endpoint   | Signature Endpoint | ASCII Signature Endpoint |
|------------|--------------------|-------------------------|
| iPXE       | `http://bootcfg.foo/ipxe.sig` | `http://bootcfg.foo/ipxe.asc` |
| Pixiecore  | `http://bootcfg/pixiecore/v1/boot.sig/:MAC` | `http://bootcfg/pixiecore/v1/boot.asc/:MAC` |
| GRUB2      | `http://bootcf.foo/grub.sig` | `http://bootcfg.foo/grub.asc` |
| Ignition   | `http://bootcfg.foo/ignition.sig` | `http://bootcfg.foo/ignition.asc` |
| Cloud-Config | `http://bootcfg.foo/cloud.sig` | `http://bootcfg.foo/cloud.asc` |
| Metadata   | `http://bootcfg.foo/metadata.sig` | `http://bootcfg.foo/metadata.asc` |

Get a config and its detached ASCII armored signature.

    GET http://bootcfg.foo/ipxe?label=value
    GET http://bootcfg.foo/ipxe.asc?label=value

**Response**

```
-----BEGIN PGP SIGNATURE-----

wsBcBAEBCAAQBQJWoDHyCRCzUpbPLRRcKAAAqQ8IAGD+eC9kzc/U7h9tgwvvWwm9
suTmVSGlzC5RwTRXg6CKuW31m3WAin2b5zWRPa7MxxanYMhhBbOfrqg/4xi1tfdE
w7ipmmgftl3re0np75Jt9K1rwGXUHTCs3yooz/zvqSvNSobG13FL5tp+Jl7a22wE
+W7x9BukTytVgNLt3IDIxsJ/rAEYUm4zySftooDbFVKj/SK5w8xg4zLmE6Jxz6wp
eaMlL1TEXy3NaFR0+hgbqM/tgeV2j6pmho8yaPF63iPnksH+gdmPiwasCfpSaJyr
NO+p24BL3PHZyKw0nsrm275C913OxEVgnNZX7TQltaweW23Cd1YBNjcfb3zv+Zo=
=mqZK
-----END PGP SIGNATURE-----
```

## Assets

If you need to serve static assets (e.g. kernel, initrd), `bootcfg` can serve arbitrary assets from the `-assets-path`.

    bootcfg.foo/assets/
    └── coreos
        └── 1032.0.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz
        └── 983.0.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

