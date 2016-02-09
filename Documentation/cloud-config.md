
# Cloud Config

CoreOS cloud-config is a system for configuring machines with a cloud-config file or executable script from user-data. Cloud-Config runs in userspace on each boot and implements a subset of the [cloud-init spec](http://cloudinit.readthedocs.org/en/latest/topics/format.html#cloud-config-data). See the cloud-config [docs](https://coreos.com/os/docs/latest/cloud-config.html) for details.

Cloud-Config files and scripts can be added in a `cloud` subdirectory of the `bootcfg` data directory. The files may contain [Go template](https://golang.org/pkg/text/template/) elements which will be rendered with `metadata` when served.

    data/
    ├── cloud
    │   ├── cloud.yaml
    │   ├── kubernetes-master.sh
    │   └── kubernetes-worker.sh
    ├── ignition
    └── specs

Add a cloud-config to a `Spec` by adding the `cloud_id` field. When PXE booting, use the kernel option `cloud-config-url` to point to `bootcfg` cloud config endpoint.

spec.json:

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

Here is an example cloud-config which starts some units and writes a file.

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

See [examples/cloud](../examples/cloud) for example cloud-config files.

### Validator

The cloud-config [validator](https://coreos.com/validate/) is useful for checking your cloud-config files for errors.

## Endpoint

The `bootcfg` [cloud-config endpoint](api.md#cloud-config) `/cloud?param=val` endpoint matches parameters to a machine `Spec` and renders the corresponding cloud-config with `metadata`.

## Comparison with Ignition

Cloud-Config starts after userspace has started and runs on every boot. Ignition starts earlier and only runs on the first boot to provision disk state. Often, tasks do not need to be repeated on each boot (e.g. writing systemd unit files) and can be performed more easily before systemd starts (e.g. configuring networking). Ignition is recommended unless a task requires re-execution on each boot.

If a service needs to be started with dynamic data, a good approach is to use Ignition to write static files which leverage systemd's environment file expansion and start a metadata service to fetch runtime data for services which require it.