
# Cloud config

**Note:** We recommend migrating to [Ignition](ignition.md) for hardware provisioning.

CoreOS Cloud-Config is a system for configuring machines with a Cloud-Config file or executable script from user-data. Cloud-Config runs in userspace on each boot and implements a subset of the [cloud-init spec](http://cloudinit.readthedocs.org/en/latest/topics/format.html#cloud-config-data). See the cloud-config [docs](https://coreos.com/os/docs/latest/cloud-config.html) for details.

Cloud-Config template files can be added in `/var/lib/matchbox/cloud` or in a `cloud` subdirectory of a custom `-data-path`. Template files may contain [Go template](https://golang.org/pkg/text/template/) elements which will be evaluated with group metadata, selectors, and query params.

```
/var/lib/matchbox
├── cloud
│   ├── cloud.yaml
│   └── script.sh
├── ignition
└── profiles
```

## Reference

Reference a Cloud-Config in a [Profile](matchbox.md#profiles) with `cloud_id`. When PXE booting, use the kernel option `cloud-config-url` to point to `matchbox` [cloud-config endpoint](api.md#cloud-config).

## Examples

Here is an example Cloud-Config which starts some units and writes a file.

<!-- {% raw %} -->
```yaml
#cloud-config
coreos:
  units:
    - name: etcd2.service
      command: start
    - name: fleet.service
      command: start
write_files:
  - path: "/home/core/welcome"
    owner: "core"
    permissions: "0644"
    content: |
      {{.greeting}}
```
<!-- {% endraw %} -->

The Cloud-Config [Validator](https://coreos.com/validate/) is also useful for checking your Cloud-Config files for errors.

## Comparison with Ignition

Cloud-Config starts after userspace has started, on every boot. Ignition starts before PID 1 and only runs on the first boot. Ignition favors immutable infrastructure.

Ignition is favored as the replacement for CoreOS Cloud-Config. Tasks often only need to be run once and can be performed more easily before systemd has started (e.g. configuring networking). Ignition can write service units for tasks that need to be run on each boot. Instead of depending on Cloud-Config variable substitution, Ignition favors using systemd's EnvironmentFile expansion to start units with a metadata file from a metadata source.
