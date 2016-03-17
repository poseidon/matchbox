
# bootcfg

`bootcfg` is a HTTP and gRPC service that renders signed [Ignition configs](https://coreos.com/ignition/docs/latest/what-is-ignition.html), [cloud-configs](https://coreos.com/os/docs/latest/cloud-config.html), network boot configs, and metadata to machines to create clusters of CoreOS machines. `bootcfg` maintains a list of **Group** definitions which match machines to profiles based on their attributes (e.g. UUID, MAC address, stage, region). A **Profile** is a named set of config templates (e.g. iPXE, GRUB, Ignition config, Cloud-Config) and metadata.

The aim is to use CoreOS Linux's early-boot capabilities to network boot and provision CoreOS machines into cluster members. Network boot endpoints provide PXE, iPXE, GRUB, and [Pixiecore](https://github.com/danderson/pixiecore/blob/master/README.api.md) support. The `bootcfg` service can be run as binary, as an [application container](https://github.com/appc/spec) with rkt, or as a Docker container.

## Getting Started

Get started running `bootcfg` on your laptop, with rkt or Docker, to network boot libvirt VMs into CoreOS clusters.

* [Getting Started with rkt](getting-started-rkt.md)
* [Getting Started with Docker](getting-started-docker.md)

## Flags

See [flags and variables](config.md)

## API

See [API](api.md)

## Data

A `Store` stores Profiles, Ignition configs, cloud-configs. By default, `bootcfg` uses a `FileStore` to search a data directory (`-data-path`) for these resources.

Prepare `/etc/bootcfg` or a custom `-data-path` with `profile`, `ignition`, and `cloud` subdirectories. You may wish to keep these files under version control. The [examples](../examples) directory is a valid target with some pre-defined configs and templates.

     /etc/bootcfg
     ├── cloud
     │   ├── cloud.yaml
     │   └── worker.sh
     ├── ignition
     │   └── hello.json
     │   └── etcd.yaml
     │   └── simple_networking.yaml
     └── profiles
         └── etcd
             └── profile.json
         └── worker
             └── profile.json

Ignition templates can be JSON or YAML files. Cloud-Config templates can be a script or YAML file. Both may contain may contain [Go template](https://golang.org/pkg/text/template/) elements which will be executed machine group [metadata](#groups-and-metadata). For details and examples:

* [Ignition Config](ignition.md)
* [Cloud-Config](cloud-config.md)

### Profiles

Profiles specify the Ignition config, Cloud-Config, and network boot config to be used by machine(s).

    {
        "id": "etcd_profile",
        "cloud_id": "",
        "ignition_id": "etcd.yaml",
        "boot": {
            "kernel": "/assets/coreos/899.6.0/coreos_production_pxe.vmlinuz",
            "initrd": ["/assets/coreos/899.6.0/coreos_production_pxe_image.cpio.gz"],
            "cmdline": {
                "cloud-config-url": "http://bootcfg.foo/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp}",
                "coreos.autologin": "",
                "coreos.config.url": "http://bootcfg.foo/ignition?uuid=${uuid}&mac=${net0/mac:hexhyp}",
                "coreos.first_boot": "1"
            }
        }
    }

The `"boot"` settings will be used to render configs to network boot programs such as iPXE, GRUB, or Pixiecore. You may reference remote kernel and initrd assets or [local assets](#assets).

To use cloud-config, set the `cloud-config-url` kernel option to reference the `bootcfg` [Cloud-Config endpoint](api.md#cloud-config), which will render the `cloud_id` file.

To use Ignition, set the `coreos.config.url` kernel option to reference the `bootcfg` [Ignition endpoint](api.md#ignition-config), which will render the `ignition_id` file. Be sure to add the `coreos.first_boot` option as well.

## Groups and Metadata

Groups define tag selectors which match zero or more machines. Machine(s) matching a group will boot and provision according to the group's `Profile` and `metadata`.

Define a list of groups, define the required tags, name the `Profile` that should be applied, and add any `metadata` needed to render the templates in your Ignition or Cloud configs.

Here is an example `/etc/bootcfg.conf` YAML file:

    ---
    api_version: v1alpha1
    groups:
      - name: default
        profile: discovery
      - name: Worker Node
        profile: worker
        require:
          region: us-central1-a
          zone: a
      - name: etcd Node 1
        profile: etcd
        require:
          uuid: 16e7d8a7-bfa9-428b-9117-363341bb330b
        metadata:
          networkd_name: ens3
          networkd_gateway: 172.15.0.1
          networkd_dns: 172.15.0.3
          networkd_address: 172.15.0.21/16
          ipv4_address: 172.15.0.21
          etcd_name: node1
          etcd_initial_cluster: "node1=http://172.15.0.21:2380"
          ssh_authorized_keys:
            - "ssh-rsa pub-key-goes-here"
      - name: etcd Proxy
        profile: etcd_proxy
        require:
          mac: 52:54:00:89:d8:10
        metadata:
          etcd_initial_cluster: "node1=http://172.15.0.21:2380"

For example, a request to `/cloud?mac=52:54:00:89:d8:10` would render the Cloud-Config template in the "etcd_proxy" `Profile`, with the machine group's metadata. A request to `/cloud` would match the default group (which has no selectors) and render the Cloud-Config in the "discovery" Profile. Avoid defining multiple default groups as resolution will not be deterministic.

### Reserved Attributes

The following attributes/tags have reserved semantic purpose. Do not use these tags for other purposes as they may be normalized or parsed specially.

* `uuid` - machine UUID
* `mac` - network interface physical address (MAC address)
* `hostname`
* `serial`

Client's booted with the `/ipxe.boot` endpoint will introspect and make a request to `/ipxe` with the `uuid`, `mac`, `hostname`, and `serial` value as query arguments. Pixiecore which can only detect MAC addresss and cannot substitute it into later config requests ([issue](https://github.com/coreos/coreos-baremetal/issues/36)).

## Assets

`bootcfg` can serve arbitrary static assets from `-assets-path` at `/assets`. This is helpful for reducing bandwidth usage when serving the kernel and initrd to network booted machines.

    bootcfg.foo/assets/
    └── coreos
        └── VERSION
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

For example, a `Profile` might refer to a local asset `/assets/coreos/VERSION/coreos_production_pxe.vmlinuz` instead of `http://stable.release.core-os.net/amd64-usr/VERSION/coreos_production_pxe.vmlinuz`.

See the [get-coreos](../scripts/README.md#get-coreos) script to quickly download, verify, and place CoreOS assets.

## Network

`bootcfg` does not implement a DHCP/TFTP server. Its easy to use the [coreos/dnsmasq](../contrib/dnsmasq) image if you need a quick DHCP, proxyDHCP, TFTP, or DNS setup.

