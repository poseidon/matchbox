
# Ignition

Ignition is a system for declaratively provisioning disks during the initramfs, before systemd starts. It runs only on the first boot and handles formatting partitions, writing files (systemd units, networkd units, dropins, regular files), and configuring users. See the Ignition [docs](https://coreos.com/ignition/docs/latest/) for details.

Ignition template files can be added in an `ignition` subdirectory of the `bootcfg` data directory. The files may contain [Go template](https://golang.org/pkg/text/template/) elements which should evaluate, with `metadata`, to Ignition JSON or to Ignition YAML (which will be rendered as JSON).

    data
     ├── cloud
     ├── ignition
     │   └── simple.json
     │   └── etcd.yaml
     │   └── etcd_proxy.yaml
     │   └── networking.yaml
     └── specs

Add an Ignition config to a `Spec` by adding the `ignition_id` field. When PXE booting, use the kernel option `coreos.first_boot=1` and `coreos.config.url` to point to the `bootcfg`ignition endpoint.

spec.json:

     {
         "id": "etcd_profile",
         "boot": {
             "kernel": "/assets/coreos/899.6.0/coreos_production_pxe.vmlinuz",
             "initrd": ["/assets/coreos/899.6.0/coreos_production_pxe_image.cpio.gz"],
             "cmdline": {
                 "coreos.config.url": "http://bootcfg.foo/ignition?uuid=${uuid}&mac=${net0/mac:hexhyp}",
                 "coreos.first_boot": "1"
             }
         },
         "cloud_id": "",
         "ignition_id": "etcd.yaml"
     }

## Configs

Here is an example Ignition config for static networking, which will be evaluated with metadata into YAML and tranformed into machine-friendly JSON.

ignition/network.yaml:

    ---
    ignition_version: 1
    networkd:
      units:
        - name: 00-{{.networkd_name}}.network
          contents: |
            [Match]
            Name={{.networkd_name}}
            [Network]
            Gateway={{.networkd_gateway}}
            DNS={{.networkd_dns}}
            DNS=8.8.8.8
            Address={{.networkd_address}}
    {{ if .ssh_authorized_keys }}
    passwd:
      users:
        - name: core
          ssh_authorized_keys:
            {{ range $element := .ssh_authorized_keys }}
            - {{$element}}
            {{end}}
    {{end}}

Response from `/ignition?mac=address` for a particular machine.

    {
      "ignitionVersion": 1,
      "storage": {},
      "systemd": {},
      "networkd": {
        "units": [
          {
            "name": "00-ens3.network",
            "contents": "[Match]\nName=ens3\n[Network]\nGateway=172.15.0.1\nDNS=172.15.0.3\nDNS=8.8.8.8\nAddress=172.15.0.21/16\n"
          }
        ]
      },
      "passwd": {}
    }

Note that Ignition does **not** allow variables - the response has been fully rendered with `metadata` for the requesting machine.

Ignition configs can be provided directly as JSON as well. This is useful for simple cases or if you prefer to use your own templating solution to generate Ignition configs.

ignition/run-hello.json:

    {
        "ignitionVersion": 1,
        "systemd": {
            "units": [
                {
                    "name": "hello.service",
                    "enable": true,
                    "contents": "[Service]\nType=oneshot\nExecStart=/usr/bin/echo Hello World\n\n[Install]\nWantedBy=multi-user.target"
                }
            ]
        }
    }

### Examples

See [examples/ignition](../examples/ignition) for example Ignition configs which setup networking, install CoreOS to disk, or start etcd.

## Endpoint

The `bootcfg` [Ignition endpoint](api.md#ignition-config) `/ignition?param=val` endpoint matches parameters to a machine `Spec` and renders the corresponding Ignition config with `metadata`, transforming YAML to JSON if needed.
