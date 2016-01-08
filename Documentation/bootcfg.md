
# Config Service

The `bootcfg` HTTP service provides virtual or physical hardware with PXE, iPXE, or Pixiecore boot settings and igntion/cloud configs based on their attributes.

The service maintains `Spec` definitions and ignition/cloud config resources and matches machines to `Spec`'s based on matcher groups you can declare. `Spec` resources define a named set of boot settings (kernel, options, initrd) and configuration settings (ignition config, cloud config). Group matchers associate zero or more machines to a `Spec` based on required attributes (e.g. UUID, MAC, region, free-form pairs).

Boot settings are presented as iPXE scripts or [Pixiecore JSON](https://github.com/danderson/pixiecore/blob/master/README.api.md) to support different network boot environments.

## Usage

The `bootcfg` service can be run as a container to boot libvirt VMs or on a provisioner host to boot baremetal machines.

Build the binary and docker image from source

    ./build
    ./docker-build

Or pull a published container image from [quay.io/repository/coreoso/bootcfg](https://quay.io/repository/coreos/bootcfg?tab=tags).

    docker pull quay.io/coreos/bootcfg:latest
    docker tag quay.io/coreos/bootcfg:latest coreos/bootcfg:latest

The latest image corresponds to the most recent `coreos-baremetal` master commit.

[Prepare a data volume](#data) with `Spec` and ignition/cloud configs. Optionally, prepare a volume of downloaded CoreOS kernel and initrd image assets that `bootcfg` should serve.

    ./scripts/get-coreos   # download CoreOS 835.9.0 to images/coreos/835.9.0
    ./scripts/get-coreos beta 877.1.0

Run the container and mount the data and images directories as volumes.

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/examples/dev:/data:Z -v $PWD/images:/images:Z coreos/bootcfg -address=0.0.0.0:8080

## Endpoints

The [API](api.md) documents the iPXE and Pixiecore endpoints, the cloud config endpoint, and image assets.

Map container port 8080 to host port 8080 to quickly check endpoints:

* iPXE Scripts: `/ipxe?uuid=val`
* Pixiecore JSON: `/pixiecore/v1/boot/:mac`
* Cloud Config: `/cloud?uuid=val`
* Ignition Config: `/ignition?uuid=val`
* Spec: `/spec/:id`
* Images: `/images`

## Data

A `Store` maintains `Spec` definitions, matcher groups, and ignition/cloud config resources. By default, `bootcfg` uses a `FileStore` to search a filesystem data directory for these resources.

Prepare a data directory or modify the [examples](../examples) provided. The `FileStore` expects `Spec` JSON files to be located at `specs/:id/spec.json` with any unique spec identifier.

You may wish to keep the data directory under version control with your other infrastructure configs, since it contains the declarative configuration of your hardware.

Ignition configs and cloud configs can be named whatever you like and dropped into `ignition` and `cloud`, respectively.

     data
     ├── config.yaml
     ├── cloud
     │   ├── etcd.yaml
     │   └── worker.yaml
     ├── ignition
     │   └── node1.json
     │   └── node2.json
     └── specs
         └── etcd
             └── spec.json
         └── worker
             └── spec.json

### Matcher Groups

Matcher groups define a set of requirements which match zero or more machines to a `Spec`. Groups have a human readable name, a `Spec` id, and a free-form map of key/value matcher requirements.

Baremetal clients network booted with `bootcfg` include hardware attributes in requests which make it simple to match baremetal instances.

* `uuid`
* `mac`
* `hostname`
* `serial`

Note that Pixiecore only provides MAC addresses and does not subsitute variables into later config requests.

Currently, matcher groups are loaded from a YAML config file specified by the `-config` flag. With containers, it is easiest to keep the file in the data path they gets mounted.

    ---
    api_version: v1alpha1
    groups:
      - name: node1
        spec: etcd1
        require:
          uuid: 16e7d8a7-bfa9-428b-9117-363341bb330b
      - name: node2
        spec: etcd2
        require:
          mac: 52:54:00:89:d8:10
      - name: workers
        spec: worker
        require:
          region: okla
          zone: a1
      - name: default
        spec: default

Machines are matched to a `Spec` by evaluating group matchers requirements by decreasing number of constraints, in deterministic order. In this example, a request to `/cloud?mac=52:54:00:89:d8:10` would serve the cloud config from the "etcd2" `Spec`.

A default group matcher can be defined by omitting the `require` field. Avoid defining multiple default groups as resolution will not be deterministic.

### Spec

Specs can have any unique identifier you choose and specify boot settings (kernel, options, initrd) and configuration settings (cloud config) for one or more machines.

Boot config files contain JSON referencing a kernel image, init RAM fileystems, and kernel options for booting a machine.

    {
        "id": "etcd2",
        "boot": {
            "kernel": "/images/coreos/835.9.0/coreos_production_pxe.vmlinuz",
            "initrd": ["/images/coreos/835.9.0/coreos_production_pxe_image.cpio.gz"],
            "cmdline": {
                "cloud-config-url": "http://bootcfg.foo/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp}",
                "coreos.autologin": "",
                "coreos.config.url": "http://bootcfg.foo/ignition?uuid=${uuid}",
                "coreos.first_boot": ""
            }
        },
        "cloud_id": "etcd.yaml",
        "ignition_id": "node2.json"
    }

The `"boot"` section references the kernel image, init RAM filesystem, and kernel options to use. Point kernel and initrd to remote images or to local [image assets](#images).

To use cloud-init, set the `cloud-config-url` kernel option to the `bootcfg` cloud endpoint to reference the cloud config named by `cloud_id`.

To use ignition, set the `coreos.config.url` kernel option to the `bootcfg` ignition endpoint to refernce the ignition config named by `ignition_id`. Be sure to add the `coreos.first_boot` kernel argument when network booting.

### Cloud Config

Cloud config files are declarative configurations for customizing early initialization of machine instances. CoreOS supports a subset of the [cloud-init project](http://cloudinit.readthedocs.org/en/latest/index.html) and supports a kernel option `cloud-config-url`. CoreOS downloads the HTTP config after kernel initialization.

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
          File added by the default cloud-config.

See the CoreOS cloud config [docs](https://coreos.com/os/docs/latest/cloud-config.html), config [validator](https://coreos.com/validate/), and [implementation](https://github.com/coreos/coreos-cloudinit) for more details or [data](../data) for examples.

## Ignition Config

Ignition is a configuration system for provisioning CoreOS instances before userspace startup. Here is an example `ignition-config.json`.

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


See the Ignition [docs](https://coreos.com/ignition/docs/latest/) and [github](https://github.com/coreos/ignition) for the latest details.

## Images

Optionally, `bootcfg` can host free-form static assets if an `-images-path` argument to a directory is provided. This is a quick way to serve kernel and init RAM filesystem images, GPG verify them at an origin, and lets client machines download images without using egress bandwidth.

    images/
    └── coreos
        └── 835.9.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

Run the `get-coreos` script to quickly download kernel and initrd images from a recent CoreOS release into an `/images` directory.

    ./scripts/get-coreos                 # stable, 835.9.0
    ./scripts/get-coreos beta 877.1.0

To reference these local images, change the `kernel` and `initrd` in a boot config file. For example, change `http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz` to `/images/coreos/835.9.0/coreos_production_pxe.vmlinuz`.

## Virtual and Physical Machine Guides

Next, setup a virtual machine network within libvirt or a baremetal machine network. Follow the [libvirt guide](virtual-hardware.md) or [baremetal guide](physical-hardware.md).
