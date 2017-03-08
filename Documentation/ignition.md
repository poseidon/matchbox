# Ignition

Ignition is a system for declaratively provisioning disks during the initramfs, before systemd starts. It runs only on the first boot and handles partitioning disks, formatting partitions, writing files (regular files, systemd units, networkd units, etc.), and configuring users. See the Ignition [docs](https://coreos.com/ignition/docs/latest/) for details.

## Fuze configs

Ignition 2.0.0+ configs are versioned, *machine-friendly* JSON documents (which contain encoded file contents). Operators should write and maintain configs in a *human-friendly* format, such as CoreOS [fuze](https://github.com/coreos/fuze) configs. As of `matchbox` v0.4.0, Fuze configs are the primary way to use CoreOS Ignition.

The [Fuze schema](https://github.com/coreos/fuze/blob/master/doc/configuration.md) formalizes and improves upon the YAML to Ignition JSON transform. Fuze provides better support for Ignition 2.0.0+, handles file content encoding, patches Ignition bugs, performs better validations, and lets services (like `matchbox`) negotiate the Ignition version required by a CoreOS client.

### Adding Fuze configs

Fuze template files can be added in the `/var/lib/matchbox/ignition` directory or in an `ignition` subdirectory of a custom `-data-path`. Template files may contain [Go template](https://golang.org/pkg/text/template/) elements which will be evaluated with group metadata, selectors, and query params.

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

### Reference

Reference an Fuze config in a [Profile](matchbox.md#profiles) with `ignition_id`. When PXE booting, use the kernel option `coreos.first_boot=1` and `coreos.config.url` to point to the `matchbox` [Ignition endpoint](api.md#ignition-config).

### Migration from v0.3.0

In v0.4.0, `matchbox` switched to using the CoreOS [fuze](https://github.com/coreos/fuze) library, which formalizes and improves upon the YAML to Ignition JSON transform. Fuze provides better support for Ignition 2.0.0+, handles file content encoding, patches Ignition bugs, and performs better validations.

Upgrade your Ignition YAML templates to match the [Fuze config schema](https://github.com/coreos/fuze/blob/master/doc/configuration.md). Typically, you'll need to do the following:

* Remove `ignition_version: 1`, Fuze configs are version-less
* Update `filesystems` section and set the `name`
* Update `files` section to use `inline` as shown below
* Replace `uid` and `gid` with `user` and `group` objects as shown above

Maintain readable inline file contents in Fuze:

```
...
files:
  - path: /etc/foo.conf
    filesystem: root
    contents:
      inline: |
        foo bar
```

Support for the older Ignition v1 format has been dropped, so CoreOS machines must be **1010.1.0 or newer**. Read the upstream Ignition v1 to 2.0.0 [migration guide](https://coreos.com/ignition/docs/latest/migrating-configs.html) to understand the reasons behind schema changes.

## Examples

Here is an example Fuze template. This template will be rendered into a Fuze config (YAML), using group metadata, selectors, and query params as template variables. Finally, the Fuze config is served to client machines as Ignition JSON.

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

The Ignition config response (formatted) to a query `/ignition?label=value` for a CoreOS instance supporting Ignition 2.0.0 would be:

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

See [examples/ignition](../examples/ignition) for numerous Fuze template examples.

### Raw Ignition

If you prefer to design your own templating solution, raw Ignition files (suffixed with `.ign` or `.ignition`) are served directly.
