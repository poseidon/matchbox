
# bootcfg

`bootcfg` is an HTTP service that renders signed [Ignition configs](https://coreos.com/ignition/docs/latest/what-is-ignition.html), [cloud-configs](https://coreos.com/os/docs/latest/cloud-config.html), network boot configs, and metadata to machines to create clusters of CoreOS machines. The service maintains a list of **Specs** which are named sets of configuration data (e.g. Ignition config, cloud-config, kernel, initrd). When started, `bootcfg` loads a list of **Group** definitions, which match machines to Specs and metadata based on attributes (e.g. UUID, MAC address, stage) passed as query arguments.

The aim is to use CoreOS Linux's early-boot capabilities to boot machines into functional cluster members with end to end [Distributed Trusted Computing](https://coreos.com/blog/coreos-trusted-computing.html). PXE, iPXE, and [Pixiecore](https://github.com/danderson/pixiecore/blob/master/README.api.md) endpoints provide support for network booting. The `bootcfg` service can be run as an [application container](https://github.com/appc/spec) with rkt, as a Docker container, or as a binary.

## Usage

Fetch the application container image (ACI) from [Quay](https://quay.io/repository/coreos/bootcfg?tab=tags).

    sudo rkt --insecure-options=image fetch docker://quay.io/coreos/bootcfg

Alternately, pull the Docker image.

    sudo docker pull quay.io/coreos/bootcfg

The `latest` image corresponds to the most recent commit on master, so choose a tagged [release](https://github.com/coreos/coreos-baremetal/releases) if you require more stability.

Get started running `bootcfg` with rkt or Docker to network boot libvirt VMs on your laptop into CoreOS clusters.

* [Getting Started with rkt](getting-started-rkt.md)
* [Getting Started with Docker](getting-started-docker.md)

Once you've tried those examples, you're ready to write your own configs.

## Data

A `Store` stores Ignition configs, cloud-configs, and named Specs. By default, `bootcfg` uses a `FileStore` to search a data directory (`-data-path`) for these resources.

Prepare a data directory similar to the [examples](../examples) directory, with `ignition`, `cloud`, and `specs` subdirectories. You might keep this directory under version control since it will define the early boot behavior of your machines.

     data
     ├── config.yaml
     ├── cloud
     │   ├── cloud.yaml
     │   └── worker.sh
     ├── ignition
     │   └── hello.json
     │   └── etcd.yaml
     │   └── simple_networking.yaml
     └── specs
         └── etcd
             └── spec.json
         └── worker
             └── spec.json

Ignition files can be JSON files or Ignition YAML. Cloud-Configs can be YAML or scripts. Both may contain may contain [Go template](https://golang.org/pkg/text/template/) elements which will be evaluated with [metadata](#groups-and-metadata). For details and examples:

* [Ignition Config](ignition.md)
* [Cloud-Config](cloud-config.md)

#### Spec

Specs specify the Ignition config, cloud-config, and PXE boot settings (kernel, options, initrd) of a matched machine.

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

The `"boot"` settings will be used to render configs to the network boot programs used in PXE, iPXE, or Pixiecore setups. You may reference remote kernel and initrd assets or [local assets](#assets).

To use cloud-config, set the `cloud-config-url` kernel option to the `bootcfg` [Cloud-Config endpoint](api.md#cloud-config) `/cloud?param=val`, which will render the `cloud_id` file.

To use Ignition, set the `coreos.config.url` kernel option to the `bootcfg` [Ignition endpoint](api.md#ignition-config) `/ignition?param=val`, which will render the `ignition_id` file. Be sure to add the `coreos.first_boot` option as well.

## Groups and Metadata

Groups define a set of required tags which match zero or more machines. Machines matching a group will boot and provision themselves according to the group's `spec` and metadata. Currently, `bootcfg` loads group definitions from a YAML config file specified by the `-config` flag. When running `bootcfg` as a container, it is easiest to keep the config file in the [data](#data) directory so it is mounted and versioned.

Define a list of named groups, name the `Spec` that should be applied, add the tags required to match the group, and add your own `metadata` needed to render your Ignition or Cloud configs.

Here is an example `bootcfg` config.yaml:

    ---
    api_version: v1alpha1
    groups:
      - name: default
        spec: discovery
      - name: Worker Node
        spec: worker
        require:
          region: us-central1-a
          zone: a
      - name: etcd Node 1
        spec: etcd
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
        spec: etcd_proxy
        require:
          mac: 52:54:00:89:d8:10
        metadata:
          etcd_initial_cluster: "node1=http://172.15.0.21:2380"

Requirements are AND'd together and evaluated from most constraints to least, in a deterministic order. For most endpoints, "tags" correspond to query arguments in machine requests. Machines are free to query `bootcfg` with additional information (query arguments) about themselves, but they must supply the required set of tags to match a group.

For example, a request to `/cloud?mac=52:54:00:89:d8:10` would render the cloud-config named in "etcd_proxy" `Spec` with the etcd proxy metadata. A request to `/cloud` would match the default group (which has no requirements) and serve the cloud-config from the "discovery" `Spec`. Avoid defining multiple default groups as resolution will not be deterministic.

### Reserved Attributes

The following attributes/tags have reserved semantic purpose. Do not use these tags for other purposes as they may be normalized or parsed specially.

* `uuid` - machine UUID
* `mac` - network interface physical address (MAC address)
* `hostname`
* `serial`

Client's booted with the `/ipxe.boot` endpoint will introspect and make a request to `/ipxe` with the `uuid`, `mac`, `hostname`, and `serial` value as query arguments. Pixiecore which can only detect MAC addresss and cannot substitute it into later config requests ([issue](https://github.com/coreos/coreos-baremetal/issues/36)).

## Assets

`bootcfg` can serve static assets from the `-assets-path` at `/assets`. This is helpful for reducing bandwidth usage when serving the kernel and initrd to network booted machines.

    assets/
    └── coreos
        └── VERSION
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

For example, a `Spec` might refer to a local asset `/assets/coreos/VERSION/coreos_production_pxe.vmlinuz` instead of `http://stable.release.core-os.net/amd64-usr/VERSION/coreos_production_pxe.vmlinuz`.

See the [get-coreos](../scripts/README.md#get-coreos) script to quickly download, verify, and move CoreOS assets to `assets`.

## Endpoints

The [API](api.md) documents the available endpoints.

## Network

`bootcfg` does not implement a DHCP/TFTP server or monitor running instances. If you need a quick DHCP, proxyDHCP, TFTP, or DNS setup, the [coreos/dnsmasq](../contrib/dnsmasq) image can create a suitable network boot environment on a virtual or physical network. Use `--net` to specify a network bridge and `--dhcp-boot` to point clients to `bootcfg`.

## Virtual and Physical Machine Guides

Next, setup a network of virtual machines with libvirt or boot a cluster of physical hardware. Follow the [libvirt guide](virtual-hardware.md) or [physical hardware guide](physical-hardware.md).
