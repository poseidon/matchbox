
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

## OpenPGP Signatures

OpenPGP signature endpoints serve ASCII armored signatures of configs. Signatures are available if the config service is provided with a `-key-ring-path` to a private keyring containing a single signing key. If the key has a passphrase, set the `BOOTCFG_PASSPHRASE` environment variable

* `http://bootcfg.example.com/boot.ipxe.asc`
* `http://bootcfg.example.com/boot.ipxe.0.asc`
* `http://bootcfg.example.com/ipxe.asc`
* `http://bootcfg.example.com/pixiecore/v1/boot.asc/:MAC`
* `http://bootcfg.example.com/cloud.asc`
* `http://bootcfg.example.com/ignition.asc`

Signature endpoints mirror the config endpoints, but provide detached signatures and are suffixed with `.asc`. For example, an iPXE config endpoint like the following:

    GET http://bootcfg.example.com/ipxe?attribute=value

**Response**

    #!ipxe
    kernel /assets/coreos/835.9.0/coreos_production_pxe.vmlinuz cloud-config-url=http://172.17.0.2:8080/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp} coreos.autologin
    initrd  /assets/coreos/835.9.0/coreos_production_pxe_image.cpio.gz
    boot

Provides a sibling OpenPGP signature endpoint.

    GET http://bootcfg.example.com/ipxe.asc?attribute=value

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

If you need to host static assets (e.g. kernel, initrd) within your network, bootcfg server's `/assets/` route serves free-form static assets. Set the `-assets-path` when starting the bootcfg server. Here is an example:

    assets/
    └── coreos
        └── 835.9.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz
        └── 877.1.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

