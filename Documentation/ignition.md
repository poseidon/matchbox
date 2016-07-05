
# Ignition

Ignition is a system for declaratively provisioning disks during the initramfs, before systemd starts. It runs only on the first boot and handles partitioning disks, formatting partitions, writing files (regular files, systemd units, networkd units, dropins), and configuring users. See the Ignition [docs](https://coreos.com/ignition/docs/latest/) for details.

## Fuze Configs

Ignition 2.0.0+ configs are versioned, *machine-friendly* JSON documents (which contain encoded file contents). Operators should write and maintain configs in a *human-friendly* format, such as CoreOS [fuze](https://github.com/coreos/fuze) configs. As of `bootcfg` v0.4.0, Fuze configs are the primary way to use CoreOS Ignition.

Fuze formalizes the transform from a [Fuze config](https://github.com/coreos/fuze/blob/master/doc/configuration.md) (YAML) to Ignition. Fuze allows services (like `bootcfg`) to negotiate versions and serve Ignition configs to different CoreOS clients. That means, you can write a Fuze config in YAML and serve Ignition to different CoreOS instances, even as we update the Ignition version shipped in the OS.

Fuze automatically handles file content encoding so that you can continue to write and maintain readable inline file contents.

```
...
files:
  - path: /etc/foo.conf
    filesystem: rootfs
    contents:
      inline: |
        foo bar
```

Fuze can also patch some bugs in shipped Ignition versions and validate configs for errors.

## Adding Fuze Configs

Fuze template files can be added in the `/var/lib/bootcfg/ignition` directory or in an `ignition` subdirectory of a custom `-data-path`. Template files may contain [Go template](https://golang.org/pkg/text/template/) elements which will be evaluated with Group `metadata` and should render to a Fuze config.

    /var/lib/bootcfg
     ├── cloud
     ├── ignition
     │   └── raw.ign
     │   └── etcd.yaml
     │   └── etcd-proxy.yaml
     │   └── networking.yaml
     └── profiles

### Referencing

Reference an Fuze config in a [Profile](bootcfg.md#profiles) with `ignition_id`. When PXE booting, use the kernel option `coreos.first_boot=1` and `coreos.config.url` to point to the `bootcfg` [Ignition endpoint](api.md#ignition-config).

## Examples

See [examples/ignition](../examples/ignition) for numerous Fuze template examples.

Here is an example Fuze template. This template will be rendered into a Fuze config (YAML), using metadata (from a Group) to fill in the template values. At query time, `bootcfg` transforms the Fuze config to Ignition for clients.

ignition/format-disk.yaml.tmpl:

    ---
    storage:
      disks:
        - device: /dev/sda
          wipe_table: true
          partitions:
            - label: ROOT
      filesystems:
        - name: rootfs
          mount:
            device: "/dev/sda1"
            format: "ext4"
            create:
              force: true
              options:
                - "-LROOT"
      files:
        - filesystem: rootfs
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

Below is the Ignition config response from `/ignition?selector=value` for a CoreOS instance supporting Ignition 2.0.0. In this case, no `"ssh_authorized_keys"` list was provided in metadata.

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
            "name": "rootfs",
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
            "filesystem": "rootfs",
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

### Raw Ignition

If you prefer to design your own templating solution, raw Ignition files (suffixed with `.ign` or `.ignition`) are served directly.

ignition/run-hello.ign:

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

### Migration from v0.3.0

`bootcfg` v0.3.0 and earlier accepted YAML templates and rendered them to Ignition v1 JSON.

In v0.4.0, `bootcfg` switched to using the CoreOS [fuze](https://github.com/coreos/fuze) library, which formalizes and improves upong the YAML to Ignition JSON transform. Fuze provides better support for Ignition 2.0.0+, handles file content encoding, patches Ignition bugs, and performs better validations.

To upgrade to bootcfg `v0.4.0`, upgrade your Ignition YAML templates to match the [Fuze config schema](https://github.com/coreos/fuze/blob/master/doc/configuration.md). CoreOS machines **must be 1010.1.0 or newer**, support for the older Ignition v1 format has been dropped.

Typically, you'll need to do the following:

* Remove `ignition_version: 1`, Fuze configs are versionless
* Update `filesystems` section and set the `name`
* Update `files` section to use `inline` as shown above
* Replace `uid` and `gid` with `user` and `group` objects as shown above

Read the upstream Ignition v1 to 2.0.0 [migration guide](https://coreos.com/ignition/docs/latest/migrating-configs.html) to better understand the reasons behind Fuze's schema.

