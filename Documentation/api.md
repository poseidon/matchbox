
# HTTP API

## iPXE script

Serves a static iPXE boot script which gathers client machine attributes and chainloads to the iPXE endpoint. Use DHCP/TFTP to point iPXE clients to this endpoint as the next-server.

```
GET http://matchbox.foo/boot.ipxe
GET http://matchbox.foo/boot.ipxe.0   // for dnsmasq
```

**Response**

```
#!ipxe
chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}
```

Client's booted with the `/ipxe.boot` endpoint will introspect and make a request to `/ipxe` with the `uuid`, `mac`, `hostname`, and `serial` value as query arguments.

## iPXE

Finds the profile for the machine and renders the network boot config (kernel, options, initrd) as an iPXE script.

```
GET http://matchbox.foo/ipxe?label=value
```

**Query parameters**

| Name | Type   | Description     |
|------|--------|-----------------|
| uuid | string | Hardware UUID   |
| mac  | string | MAC address     |
| *    | string | Arbitrary label |

**Response**

```
#!ipxe
kernel /assets/coreos/1465.7.0/coreos_production_pxe.vmlinuz coreos.config.url=http://matchbox.foo:8080/ignition?uuid=${uuid}&mac=${mac:hexhyp} coreos.first_boot=1 coreos.autologin
initrd  /assets/coreos/1465.7.0/coreos_production_pxe_image.cpio.gz
boot
```

## GRUB2

Finds the profile for the machine and renders the network boot config as a GRUB config. Use DHCP/TFTP to point GRUB clients to this endpoint as the next-server.

```
GET http://matchbox.foo/grub?label=value
```

**Query parameters**

| Name | Type   | Description     |
|------|--------|-----------------|
| uuid | string | Hardware UUID   |
| mac  | string | MAC address     |
| *    | string | Arbitrary label |

**Response**

```
default=0
timeout=1
menuentry "CoreOS" {
echo "Loading kernel"
linuxefi "(http;matchbox.foo:8080)/assets/coreos/1465.7.0/coreos_production_pxe.vmlinuz" "coreos.autologin" "coreos.config.url=http://matchbox.foo:8080/ignition" "coreos.first_boot"
echo "Loading initrd"
initrdefi "(http;matchbox.foo:8080)/assets/coreos/1465.7.0/coreos_production_pxe_image.cpio.gz"
}
```

## Cloud config

DEPRECATED: Finds the profile matching the machine and renders the corresponding Cloud-Config with group metadata, selectors, and query params.

```
GET http://matchbox.foo/cloud?label=value
```

**Query Parameters**

| Name | Type   | Description     |
|------|--------|-----------------|
| uuid | string | Hardware UUID   |
| mac  | string | MAC address     |
| *    | string | Arbitrary label |

**Response**

```yaml
#cloud-config
coreos:
  units:
    - name: etcd2.service
      command: start
    - name: fleet.service
      command: start
```

## Container Linux Config / Ignition Config

Finds the profile matching the machine and renders the corresponding Ignition Config with group metadata, selectors, and query params.

```
GET http://matchbox.foo/ignition?label=value
```

**Query parameters**

| Name | Type   | Description     |
|------|--------|-----------------|
| uuid | string | Hardware UUID   |
| mac  | string | MAC address     |
| *    | string | Arbitrary label |

**Response**

```json
{
  "ignition": { "version": "2.0.0" },
  "systemd": {
    "units": [{
      "name": "example.service",
      "enable": true,
      "contents": "[Service]\nType=oneshot\nExecStart=/usr/bin/echo Hello World\n\n[Install]\nWantedBy=multi-user.target"
    }]
  }
}
```

## Generic config

Finds the profile matching the machine and renders the corresponding generic config with group metadata, selectors, and query params.

```
GET http://matchbox.foo/generic?label=value
```

**Query parameters**

| Name | Type   | Description     |
|------|--------|-----------------|
| uuid | string | Hardware UUID   |
| mac  | string | MAC address     |
| *    | string | Arbitrary label |

**Response**

```
{
  “uuid”: “”,
  “mac”: “52:54:00:a1:9c:ae”,
  “osInstalled”: true,
  “rawQuery”: “mac=52:54:00:a1:9c:ae&os=installed”
}
```

## Metadata

Finds the matching machine group and renders the group metadata, selectors, and query params in an "env file" style response.

```
GET http://matchbox.foo/metadata?mac=52-54-00-a1-9c-ae&foo=bar&count=3&gate=true
```

**Query Parameters**

| Name | Type   | Description     |
|------|--------|-----------------|
| uuid | string | Hardware UUID   |
| mac  | string | MAC address     |
| *    | string | Arbitrary label |

**Response**

```
META=data
ETCD_NAME=node1
SOME_NESTED_DATA=some-value
MAC=52:54:00:a1:9c:ae
REQUEST_QUERY_MAC=52:54:00:a1:9c:ae
REQUEST_QUERY_FOO=bar
REQUEST_QUERY_COUNT=3
REQUEST_QUERY_GATE=true
REQUEST_RAW_QUERY=mac=52-54-00-a1-9c-ae&foo=bar&count=3&gate=true
```

## OpenPGP signatures

OpenPGPG signature endpoints serve detached binary and ASCII armored signatures of rendered configs, if enabled. See [OpenPGP Signing](openpgp.md).

| Endpoint   | Signature Endpoint | ASCII Signature Endpoint |
|------------|--------------------|-------------------------|
| iPXE       | `http://matchbox.foo/ipxe.sig` | `http://matchbox.foo/ipxe.asc` |
| GRUB2      | `http://bootcf.foo/grub.sig` | `http://matchbox.foo/grub.asc` |
| Ignition   | `http://matchbox.foo/ignition.sig` | `http://matchbox.foo/ignition.asc` |
| Cloud-Config | `http://matchbox.foo/cloud.sig` | `http://matchbox.foo/cloud.asc` |
| Generic    | `http://matchbox.foo/generic.sig` | `http://matchbox.foo/generic.asc` |
| Metadata   | `http://matchbox.foo/metadata.sig` | `http://matchbox.foo/metadata.asc` |

Get a config and its detached ASCII armored signature.

```
GET http://matchbox.foo/ipxe?label=value
GET http://matchbox.foo/ipxe.asc?label=value
```

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

If you need to serve static assets (e.g. kernel, initrd), `matchbox` can serve arbitrary assets from the `-assets-path`.

```
matchbox.foo/assets/
└── coreos
    └── 1465.7.0
        ├── coreos_production_pxe.vmlinuz
        └── coreos_production_pxe_image.cpio.gz
    └── 1153.0.0
        ├── coreos_production_pxe.vmlinuz
        └── coreos_production_pxe_image.cpio.gz
```
