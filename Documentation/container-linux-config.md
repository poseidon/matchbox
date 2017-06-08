# Container Linux Configs

A Container Linux Config is a YAML document which declares how Container Linux instances' disks should be provisioned on network boot and first-boot from disk. Configs can declare disk paritions, write files (regular files, systemd units, networkd units, etc.), and configure users. See the Container Linux Config [spec](https://coreos.com/os/docs/latest/configuration.html).

### Ignition

Container Linux Configs are validated and converted to *machine-friendly* Ignition configs (JSON) by matchbox when serving to booting machines. [Ignition](https://coreos.com/ignition/docs/latest/), the provisioning utility shipped in Container Linux, will parse and execute the Ignition config to realize the desired configuration. Matchbox users usually only need to write Container Linux Configs.

*Note: Container Linux directory names are still named "ignition" for historical reasons as outlined below. A future breaking change will rename to "container-linux-config".*

## Adding Container Linux Configs

Container Linux Config templates can be added to the `/var/lib/matchbox/ignition` directory or in an `ignition` subdirectory of a custom `-data-path`. Template files may contain [Go template](https://golang.org/pkg/text/template/) elements which will be evaluated with group metadata, selectors, and query params.

```
/var/lib/matchbox
 ├── cloud
 ├── ignition
 │   └── k8s-controller.yaml
 │   └── etcd.yaml
 │   └── k8s-worker.yaml
 │   └── raw.ign
 └── profiles
```

## Referencing in Profiles

Profiles can include a Container Linux Config for provisioning machines. Specify the Container Linux Config in a [Profile](matchbox.md#profiles) with `ignition_id`. When PXE booting, use the kernel option `coreos.first_boot=1` and `coreos.config.url` to point to the `matchbox` [Ignition endpoint](api.md#ignition-config).

## Examples

Here is an example Container Linux Config template. Variables will be interpreted using group metadata, selectors, and query params. Matchbox will convert the config to Ignition to serve Container Linux machines.

ignition/format-disk.yaml.tmpl:

<!-- {% raw %} -->
```yaml

---
storage:
  disks:
    - device: /dev/sda
      wipe_table: true
      partitions:
        - label: ROOT
  filesystems:
    - name: root
      mount:
        device: "/dev/sda1"
        format: "ext4"
        create:
          force: true
          options:
            - "-LROOT"
  files:
    - filesystem: root
      path: /home/core/foo
      mode: 0644
      user:
        id: 500
      group:
        id: 500
      contents:
        inline: |
          {{.example_contents}}
{{ if index . "ssh_authorized_keys" }}
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        {{ range $element := .ssh_authorized_keys }}
        - {{$element}}
        {{end}}
{{end}}
```
<!-- {% endraw %} -->

The Ignition config response (formatted) to a query `/ignition?label=value` for a Container Linux instance supporting Ignition 2.0.0 would be:

```json
{
  "ignition": {
    "version": "2.0.0",
    "config": {}
  },
  "storage": {
    "disks": [
      {
        "device": "/dev/sda",
        "wipeTable": true,
        "partitions": [
          {
            "label": "ROOT",
            "number": 0,
            "size": 0,
            "start": 0
          }
        ]
      }
    ],
    "filesystems": [
      {
        "name": "root",
        "mount": {
          "device": "/dev/sda1",
          "format": "ext4",
          "create": {
            "force": true,
            "options": [
              "-LROOT"
            ]
          }
        }
      }
    ],
    "files": [
      {
        "filesystem": "root",
        "path": "/home/core/foo",
        "contents": {
          "source": "data:,Example%20file%20contents%0A",
          "verification": {}
        },
        "mode": 420,
        "user": {
          "id": 500
        },
        "group": {
          "id": 500
        }
      }
    ]
  },
  "systemd": {},
  "networkd": {},
  "passwd": {}
}
```

See [examples/ignition](../examples/ignition) for numerous Container Linux Config template examples.

### Raw Ignition

If you prefer to design your own templating solution, raw Ignition files (suffixed with `.ign` or `.ignition`) are served directly.
