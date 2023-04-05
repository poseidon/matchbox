# Ignition Configs

[Ignition](https://coreos.github.io/ignition/) configs define how disks should be provisioned (on network boot and first-boot from disk) to partition disks, write files (regular files, systemd units, networkd units, etc.), and configure users. Ignition is used by:

* Fedora CoreOS
* RHEL CoreOS
* Flatcar Linux

See the Ignition Config v3.x [specs](https://coreos.github.io/ignition/specs/) for details.

## Usage

Ignition configs can be added to the `/var/lib/matchbox/ignition` directory or in an `ignition` subdirectory of a custom `-data-path`. Ignition configs must end in `.ign` or `ignition`.

```
/var/lib/matchbox
 ├── ignition
 │   └── k8s-controller.ign
 │   └── k8s-worker.ign
 └── profiles
```

Matchbox Profiles can set an Ignition config for provisioning machines. Specify the Ignition config in a [Profile](matchbox.md#profiles) with `ignition_id`.

```json
{
  "id": "worker",
  "name": "My Profile",
  "boot": {
    ...
  },
  "ignition_id": "my-ignition.ign"
}
```

When PXE booting, set kernel arguments depending on the OS (e.g. `ignition.firstboot` on FCOS, `flatcar.first_boot=yes` on Flatcar).

* [Fedora CoreOS](https://github.com/poseidon/matchbox/blob/main/examples/profiles/fedora-coreos.json)
* [Flatcar Linux](https://github.com/poseidon/matchbox/blob/main/examples/profiles/flatcar.json)

Point the `ignition.config.url` or `flatcar.config.url` to point to the `matchbox` [Ignition endpoint](api-http.md#ignition-config).

Matchbox parses Ignition configs (e.g. `.ign` or `.ignition`) at spec v3.3 or below and renders to the current supported version (v3.3). This relies on Ignition's [forward compatibility](https://github.com/coreos/ignition/blob/main/config/v3_3/config.go#L61).

## Writing Configs

Ignition configs can be prepared externally and loaded via the gRPC API, rather than writing Ignition by hand.

### Terraform

Terraform can be used to prepare Ignition configs, while providing integrations with external systems and rich templating. Using tools like [poseidon/terraform-provider-ct](https://github.com/poseidon/terraform-provider-ct), you can write Butane config (an easier YAML format), validate configs, and load Ignition into Matchbox ([examples](https://github.com/poseidon/matchbox/tree/main/examples/terraform)).

Define a Butane config for Fedora CoreOS or Flatcar Linux:

```yaml
variant: fcos
version: 1.4.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

```yaml
variant: flatcar
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

Define a `ct_config` data source with strict validation. Optionally use Terraform [templating](https://github.com/poseidon/terraform-provider-ct).

```tf
data "ct_config" "worker" {
  content      = file("worker.yaml")
  strict       = true
  pretty_print = false

  snippets = [
    file("units.yaml"),
    file("storage.yaml"),
  ]
}
```

Then render the Butane config to Ignition and use it in a Matchbox Profile.

```tf
resource "matchbox_profile" "fedora-coreos-install" {
  name   = "worker"
  kernel = "/assets/fedora-coreos/fedora-coreos-${var.os_version}-live-kernel-x86_64"
  initrd = [
    "--name main /assets/fedora-coreos/fedora-coreos-${var.os_version}-live-initramfs.x86_64.img"
  ]

  args = [
    "initrd=main",
    "coreos.live.rootfs_url=${var.matchbox_http_endpoint}/assets/fedora-coreos/fedora-coreos-${var.os_version}-live-rootfs.x86_64.img",
    "coreos.inst.install_dev=/dev/vda",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?uuid=${uuid}&mac=${mac:hexhyp}",
  ]

  raw_ignition = data.ct_config.worker.rendered
}
```

See the Terraform [examples](https://github.com/poseidon/matchbox/tree/main/examples#terraform-examples) for details.

### Butane

The [Butane](https://coreos.github.io/butane/) command line tool can be used to convert Butane configs (an easier YAML format) to Ignition. Then you can use the Matchbox gRPC API to upload the rendered Ignition to Matchbox for serving to machines on boot.

See [examples/ignition](../examples/ignition) for Butane config examples.

### Matchbox Rendering

While Matchbox recommends preparing Ignition configs externally (e.g. using Terraform's rich templating), Matchbox does still support limited templating and translation features with a builtin Butane converter.

Specify a Butane config in a [Profile](matchbox.md#profiles) with `ignition_id` (file must not end in `.ign` or `.ignition`).

```json
{
  "id": "worker",
  "name": "My Profile",
  "boot": {
    ...
  },
  "ignition_id": "butane.yaml"
}
```

Here is an example Butane config with Matchbox template elements. Template files may contain [Go template](https://golang.org/pkg/text/template/) elements which will be interpreted using group metadata, selectors, and query params.

```yaml
variant: flatcar
version: 1.0.0
storage:
  files:
    - path: /var/home/core/foo
      mode: 0644
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

Matchbox will use the Butane library to config to the current supported Ignition version. This relies on Ignition's [forward compatibility](https://github.com/coreos/ignition/blob/main/config/v3_3/config.go#L61).
