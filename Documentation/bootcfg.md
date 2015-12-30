
# Boot Config Service

The `bootcfg` HTTP service provides configs to PXE, iPXE, and Pixiecore network boot clients based on their hardware attributes to boot and configure virtual or physical machines.

Boot configs (i.e. kernel, initrd, kernel options) and cloud configs can be declared for machines by UUID, MAC address, or as the default for machines. The service renders boot configs as iPXE scripts and as JSON responses to implement the Pixiecore [API spec](https://github.com/danderson/pixiecore/blob/master/README.api.md).

Currently, `bootcfg` is a proof of concept, but it can make it easier to declare the desired state of network booted machines and get started with clusters of virtual or physical machines.

## Usage

The `bootcfg` service can be run as a container to boot libvirt VMs or on a provisioner host to boot baremetal machines.

Build the binary and docker image from source

    ./build
    ./docker-build

Or pull a published container image from [quay.io/repository/coreoso/bootcfg](https://quay.io/repository/coreos/bootcfg?tab=tags).

    docker pull quay.io/coreos/bootcfg:latest
    docker tag quay.io/coreos/bootcfg:latest coreos/bootcfg:latest

The latest image corresponds to the most recent `coreos-baremetal` master commit.

Prepare a directory with machine [configs](#configs) and download CoreOS kernel and initrd images that `bootcfg` should serve (optional).

    ./scripts/get-coreos   # download CoreOS 835.9.0 to images/coreos/835.9.0
    ./scripts/get-coreos beta 877.1.0

Run the container and mount the configs and images directories as volumes.

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/data:/data:Z -v $PWD/images:/images:Z coreos/bootcfg -address=0.0.0.0:8080

## Endpoints

The [API](api.md) documents the iPXE and Pixiecore boot config endpoints, the cloud config endpoint, and image assets.

Map container port 8080 to host port 8080 to quickly check endpoints:

* iPXE Scripts: `/ipxe?uuid=val`
* Pixiecore JSON: `/pixiecore/v1/boot/:mac`
* Cloud Config: `/cloud?uuid=val`
* Images: `/images`

## Configs

A `Store` maintains associations between machine attributes and different types of bootstrapping configs. Currently, `bootcfg` includes a `FileStore` which can search a filesystem directory for `boot` and `cloud` config files.

Prepare a directory of config data or use the example provided in [data](../data). The `FileStore` expects `boot` and `cloud` files in subdirectories, with files in nested `uuid` and `mac` subdirectories.

    data
    ├── boot
    │   └── default
    └── cloud
        ├── default
        └── uuid
            └── 1cff2cd8-f00a-42c8-9426-f55e6a1847f6
        └── mac
            └── 52:54:00:c7:b6:64

To find boot configs and cloud configs, the `FileStore` searches the `uuid` directory for a file matching a client machine's UUID, then searches `mac` for file matching the client's MAC address, and finally falls back to using the `default` file if present.

You may keep the config data directory in a separate location to keep it under version control with other declarative configs.

### Boot Config

Boot config files contain JSON referencing a kernel image, init RAM fileystems, and kernel options for booting a machine.

    {
        "kernel": "/images/coreos/835.9.0/coreos_production_pxe.vmlinuz",
        "initrd": ["/images/coreos/835.9.0/coreos_production_pxe_image.cpio.gz"],
        "cmdline": {
            "cloud-config-url": "http://172.17.0.2:8080/cloud?uuid=${uuid}",
            "coreos.autologin": ""
        }
    }

Point kernel and initrd to image URIs or to local [assets](#assets). If the kernel supports the `cloud-config-url` option, you can set the URL to refer to the name or IP where `bootcfg` is running so cloud configs can also be served based on hardware attributes.

In this example, `bootcfg` was expected to run as the first container within Docker's virtual bridge subnet 172.17.0.0/16 under libvirt.

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

## Image Assets

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
