
# Cloud Config

CoreOS Cloud-Config is a system for configuring machines with a Cloud-Config file or executable script from user-data. Cloud-Config runs in userspace on each boot and implements a subset of the [cloud-init spec](http://cloudinit.readthedocs.org/en/latest/topics/format.html#cloud-config-data). See the cloud-config [docs](https://coreos.com/os/docs/latest/cloud-config.html) for details.

Cloud-Config template files can be added in the `/etc/bootcfg/cloud` directory or in a `cloud` subdirectory of a custom `-data-path`. Template files may contain [Go template](https://golang.org/pkg/text/template/) elements which will be evaluated with `metadata` when served.

    data/
    ├── cloud
    │   ├── cloud.yaml
    │   ├── kubernetes-master.sh
    │   └── kubernetes-worker.sh
    ├── ignition
    └── profiles

Reference a Cloud-Config in a [Profile](bootcfg.md#profiles). When PXE booting, use the kernel option `cloud-config-url` to point to `bootcfg` [cloud-config endpoint](api.md#cloud-config).

profile.json:

    {
        "id": "worker_profile",
        "cloud_id": "worker.yaml",
        "ignition_id": "",
        "boot": {
            "kernel": "/assets/coreos/899.6.0/coreos_production_pxe.vmlinuz",
            "initrd": ["/assets/coreos/899.6.0/coreos_production_pxe_image.cpio.gz"],
            "cmdline": {
                "cloud-config-url": "http://bootcfg.foo/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp}"
            }
        }
    }

## Configs

Here is an example Cloud-Config which starts some units and writes a file.

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

### Examples

See [examples/cloud](../examples/cloud) for example Cloud-Config files.

### Validator

The Cloud-Config [Validator](https://coreos.com/validate/) is useful for checking your Cloud-Config files for errors.

## Comparison with Ignition

Cloud-Config starts after userspace has started, on every boot.Ignition starts before PID 1 and only runs on the first boot. Ignition favors immutable infrastructure.

Ignition is favored as the eventual replacement for CoreOS Cloud-Config. Tasks often only need to be run once and can be performed more easily before systemd has started (e.g. configuring networking). Ignition can write service units for tasks that need to be run on each boot. Instead of depending on Cloud-Config variable substitution, leverage systemd's EnvironmentFile expansion to start units with a metadata file from a source of truth.
