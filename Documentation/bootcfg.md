
# Boot Config Service

The `bootcfg` HTTP service provides virtual or physical machines with PXE, iPXE, or Pixiecore boot settings and igntion/cloud configs based on their hardware attributes.

The service maintains `Machine`, `Spec`, and ignition/cloud config resources and matches machines to `Spec`'s based on their attributes. `Spec` resources define a named set of boot settings (kernel, options, initrd) and configuration settings (ignition config, cloud config). `Machine` resources declare a machine by id (UUID, MAC) and the `Spec` that machine should use.

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

[Prepare a data volume](#data) with `Machine`, `Spec`, and ignition/cloud configs. Optionally, prepare a volume of downloaded CoreOS kernel and initrd image assets that `bootcfg` should serve.

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
* Machines: `/machine/:id`
* Spec: `/spec/:id`
* Images: `/images`

## Data

A `Store` maintains `Machine`, `Spec`, and ignition/cloud config resources. By default, `bootcfg` uses a `FileStore` to search a filesystem data directory for these resources.

Prepare a data directory or modify the example [data](../data) provided. The `FileStore` expects `Machine` JSON files to be located at `machines/:id/machine.json` where the id can be a UUID or MAC address or "default". `Spec` JSON files should be located at `specs/:id/spec.json` with any unique spec identifier.

You may wish to keep the data directory under version control with your other infrastructure configs, since it contains the declarative configuration of your hardware.

Ignition configs and cloud configs can be named whatever you like and dropped into `ignition` and `cloud`, respectively.

     data
     ├── cloud
     │   ├── node1.yml
     │   └── orion.yml
     ├── ignition
     │   └── node2.json
     ├── machines
     │   ├── 074fbe06-94a9-4336-9e8a-20b6f81efb6c
     │   │   └── machine.json
     │   ├── 2d9354a2-e8db-4021-bff5-20ffdf443d6f
     │   │   └── machine.json
     │   └── default
     │       └── machine.json
     └── specs
         └── orion
             └── spec.json

### Machine

A machine file contains JSON describing a machine and the `Spec` that machine should use. Machines may embed `Spec` data or reference an existing (shared) spec by `spec_id`.

Here is a `machine.json` which embeds `"spec"` properties.

    {
        "id": "2d9354a2-e8db-4021-bff5-20ffdf443d6f",
        "spec": {
          "boot": {
              "kernel": "/images/coreos/877.1.0/coreos_production_pxe.vmlinuz",
              "initrd": ["/images/coreos/877.1.0/coreos_production_pxe_image.cpio.gz"],
              "cmdline": {
                  "cloud-config-url": "http://bootcfg.foo/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp}",
                  "coreos.autologin": "",
                  "coreos.config.url": "http://bootcfg.foo/ignition?uuid=${uuid}",
                  "coreos.first_boot": ""
              }
          },
          "cloud_id": "cloud-config.yml",
          "ignition_id": "ignition.json"
        },
        "spec_id": ""
    }

and another `machine.json` which simplify references a `Spec` by id.

    {
        "id": "074fbe06-94a9-4336-9e8a-20b6f81efb6c",
        "spec_id": "orion"
    }

### Spec

Specs can have any unique identifier you choose and specify boot settings (kernel, options, initrd) and configuration settings (cloud config) for one or more machines.

Boot config files contain JSON referencing a kernel image, init RAM fileystems, and kernel options for booting a machine.

    {
        "id": "orion",
        "boot": {
            "kernel": "/images/coreos/835.9.0/coreos_production_pxe.vmlinuz",
            "initrd": ["/images/coreos/835.9.0/coreos_production_pxe_image.cpio.gz"],
            "cmdline": {
                "cloud-config-url": "http://172.17.0.2:8080/cloud?uuid=${uuid}&mac=${net0/mac:hexhyp}",
                "coreos.autologin": ""
            }
        },
        "cloud_id": "orion-cloud-config.yml",
        "ignition_id": "ignition.json"
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
