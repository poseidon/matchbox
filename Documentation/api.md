
# API

## iPXE Script

Serves a static iPXE boot script which gathers client machine attributes and chainloads to the iPXE endpoint. Configure your DHCP server or iPXE server to boot from this script (e.g. set `dhcp-boot:dhcp-boot=tag:ipxe,http://bootcfg.domain.com/ipxe/boot.ipxe` if using `dnsmasq`).

    GET http://bootcfg.foo/boot.ipxe
    GET http://bootcfg.foo/boot.ipxe.0   // for dnsmasq

**Response**

    #!ipxe
    chain ipxe?uuid=${uuid}&mac=${net0/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}

## iPXE

Finds the spec matching the attribute query parameters and renders the spec boot settings (kernel, options, initrd) as an iPXE script.

    GET http://bootcfg.foo/ipxe

**Query Parameters**

| Name | Type   | Description   |
|------|--------|---------------|
| uuid | string | Hardware UUID |
| mac  | string | MAC address   |

**Response**

    #!ipxe
    kernel /assets/coreos/899.6.0/coreos_production_pxe.vmlinuz cloud-config-url=http://bootcfg.foo:8080/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp} coreos.autologin
    initrd  /assets/coreos/899.6.0/coreos_production_pxe_image.cpio.gz
    boot

## Pixiecore

Finds the spec matching the attribute query parameters and renders the boot settings as JSON to implement the Pixiecore API [spec](https://github.com/danderson/pixiecore/blob/master/README.api.md). Currently, Pixiecore only provides the machine's MAC address for matching.

    GET http://bootcfg.foo/pixiecore/v1/boot/:MAC

**URL Parameters**

| Name | Type   | Description |
|------|--------|-------------|
| mac  | string | MAC address |

**Response**

    {
      "kernel":"/assets/coreos/899.6.0/coreos_production_pxe.vmlinuz",
      "initrd":["/assets/coreos/899.6.0/coreos_production_pxe_image.cpio.gz"],
      "cmdline":{
        "cloud-config-url":"http://bootcfg.foo/cloud?mac=ADDRESS",
        "coreos.autologin":""
      }
    }

## Cloud Config

Finds the spec matching the attribute query parameters and renders the corresponding cloud-config file.

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

Finds the spec matching the attribute query parameters and renders the corresponding Ignition config as JSON.

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

## OpenPGP Signatures

OpenPGPG signature endpoints serve ASCII armored detached signatures of rendered configs when signing is enabled. See [OpenPGP Signing](openpgp.md).

| Endpoint   | ASCII Signature Endpoint |
|------------|-----------------|
| Ignition   | `http://bootcfg.foo/ignition.asc` |
| Cloud-init | `http://bootcfg.foo/cloud.asc` |
| iPXE       | `http://bootcfg.foo/boot.ipxe.asc` |
| iPXE       | `http://bootcfg.foo/ipxe.asc` |
| Pixiecore  | `http://bootcfg.foo/pixiecore/v1/boot.asc/:MAC` |

Get an Ignition config and its detached signature.

    GET http://bootcfg.foo/ipxe?attribute=value

**Response**

    #!ipxe
    kernel /assets/coreos/899.6.0/coreos_production_pxe.vmlinuz cloud-config-url=http://bootcfg.foo:8080/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp} coreos.autologin
    initrd  /assets/coreos/899.6.0/coreos_production_pxe_image.cpio.gz
    boot

    GET http://bootcfg.foo/ipxe.asc?attribute=value

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

## API Resources

### Specs

Get a `Spec` definition by id (UUID, MAC).

    http://bootcfg.foo/spec/:id

**URL Parameters**

| Name | Type   | Description |
|------|--------|-------------|
| id   | string | spec identifier |

**Response**

```json
{
  "id": "etcd",
  "boot": {
    "kernel": "/assets/coreos/899.6.0/coreos_production_pxe.vmlinuz",
    "initrd": [
      "/assets/coreos/899.6.0/coreos_production_pxe_image.cpio.gz"
    ],
    "cmdline": {
      "coreos.autologin": "",
      "coreos.config.url": "http://bootcfg.foo:8080/ignition?uuid=${uuid}&mac=${net0/mac:hexhyp}",
      "coreos.first_boot": ""
    }
  },
  "cloud_id": "",
  "ignition_id": "etcd.yaml"
}
```

## Assets

If you need to serve static assets (e.g. kernel, initrd), `bootcfg` can serve arbitrary assets from `-assets-path` at `/assets/`.

    assets/
    └── coreos
        └── 835.9.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz
        └── 899.6.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

