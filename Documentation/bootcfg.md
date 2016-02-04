
# Config Service

The bare metal config service is an HTTP service run inside a data center which provides (optionally signed) boot configs, ignition configs, and cloud configs to PXE, iPXE, and Pixiecore network-boot client machines.

The service maintains **Spec** resources which define a named set of boot settings (kernel, options, initrd) and configuration settings (ignition config, cloud config). **Groups** match zero or more machines to a `Spec` based on tags such as machine attributes (e.g. UUID, MAC) or arbitrary key/value pairs (e.g. zone, region, etc.).

The aim is to declare the desired boot and kernel/userspace provisoning behavior of machines so they come online as functioning clusters, while supporting multiple network boot environments as entrypoints.

Currently, iPXE and [Pixiecore](https://github.com/danderson/pixiecore/blob/master/README.api.md) network boot environments are supported. End to end [Distributed Trusted Computing](https://coreos.com/blog/coreos-trusted-computing.html) is a goal.

## Usage

The config service (`bootcfg`) can be run as a container to boot libvirt VMs or on a provisioner host to boot baremetal machines.

Build the binary and docker image from source

    ./build
    ./build-docker

Or pull a published container image from [quay.io/repository/coreoso/bootcfg](https://quay.io/repository/coreos/bootcfg?tab=tags).

    docker pull quay.io/coreos/bootcfg:latest
    docker tag quay.io/coreos/bootcfg:latest coreos/bootcfg:latest

The latest image corresponds to the most recent `coreos-baremetal` master commit.

[Prepare a data volume](#data) with `Spec` and ignition/cloud configs. Optionally, prepare a volume of downloaded CoreOS kernel and initrd image assets that `bootcfg` should serve.

    ./scripts/get-coreos               # CoreOS Beta 899.6.0
    ./scripts/get-coreos alpha 942.0.0

Run the container and mount the data and assets directories as volumes.

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/examples/dev:/data:Z -v $PWD/assets:/assets:Z coreos/bootcfg -address=0.0.0.0:8080 [-log-level=debug]

## Endpoints

The [API](api.md) documents the iPXE and Pixiecore endpoints, the cloud config endpoint, and assets.

Map container port 8080 to host port 8080 to quickly check endpoints:

* iPXE Scripts: `/ipxe?uuid=val`
* Pixiecore JSON: `/pixiecore/v1/boot/:mac`
* Cloud Config: `/cloud?uuid=val`
* Ignition Config: `/ignition?uuid=val`
* Spec: `/spec/:id`
* Assets: `/assets`

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

### Groups

Groups define a set of tag requirements which match zero or more machines to a `Spec`. Groups have a human readable name, a `Spec` id, and a free-form map of key/value tag requirements.

Several tags have reserved semantic purpose. You cannot use these tags for other purposes.

* `uuid` - machine UUID
* `mac` - network interface physical address (MAC address) in normalized form (e.g. `01:ab:23:cd:67:89`)
* `hostname`
* `serial`

Client's booted with the Config service include `uuid`, `mac`, `hostname`, and `serial` arguments in their requests. The exception is with Pixiecore which can only detect MAC addresss and cannot substitute it into later config requests.

Currently, group definitions are loaded from a YAML config file specified by the `-config` flag. With containers, it is easiest to keep the file in the data path they gets mounted.

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

Machines are matched to a `Spec` by evaluating group tag requirements from most constraints to least, in a deterministic order. Machines may supply extra arguments, but every tag requirement must be satisfied to match a group (i.e. AND operation). With the groups defined above, a request to `/cloud?mac=52:54:00:89:d8:10` would serve the cloud config from the "etcd2" `Spec`.

A default group can be defined by omitting the `require` field. Avoid defining multiple default groups as resolution will not be deterministic.

### Spec

Specs can have any unique identifier you choose and specify boot settings (kernel, options, initrd) and configuration settings (cloud config) for one or more machines.

Boot config files contain JSON referencing a kernel image, init RAM fileystems, and kernel options for booting a machine.

    {
        "id": "etcd2",
        "boot": {
            "kernel": "/assets/coreos/835.9.0/coreos_production_pxe.vmlinuz",
            "initrd": ["/assets/coreos/835.9.0/coreos_production_pxe_image.cpio.gz"],
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

The `"boot"` section references the kernel image, init RAM filesystem, and kernel options to use. Point kernel and initrd to remote images or to local [assets](#assets).

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

## OpenPGP Signatures

OpenPGP signature endpoints serve ASCII armored signatures of configs. Signatures are available if the config service is provided with a `-key-ring-path` to a private keyring containing a single signing key. If the key has a passphrase, set the `BOOTCFG_PASSPHRASE` environment variable.

    docker run -p 8080:8080 -e BOOTCFG_PASSPHRASE=phrase --rm -v $PWD/examples/dev:/data:Z -v $PWD/assets:/assets:Z coreos/bootcfg -address=0.0.0.0:8080 -key-ring-path /data/secring.gpg [-log-level=debug]

It is recommended that a subkey be used and exported to a key ring which is solely used for config signing and can be revoked by a master if needed. If running the config service on a Kubernetes cluster, Kubernetes secrets provide a better way to mount the key ring and source a passphrase variable.

Signature endpoints mirror the config endpoints, but provide detached signatures and are suffixed with `.asc`.

* `http://bootcfg.example.com/boot.ipxe.asc`
* `http://bootcfg.example.com/boot.ipxe.0.asc`
* `http://bootcfg.example.com/ipxe.asc`
* `http://bootcfg.example.com/pixiecore/v1/boot.asc/:MAC`
* `http://bootcfg.example.com/cloud.asc`
* `http://bootcfg.example.com/ignition.asc`

## Assets

Optionally, `bootcfg` can host free-form static assets if an `-assets-path` argument to a directory is provided. This is a quick way to serve kernel and initrd assets and reduce bandwidth usage.

    assets/
    └── coreos
        └── 835.9.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

Run the `get-coreos` script to quickly download kernel and initrd image assets.

    ./scripts/get-coreos                 # beta, 899.6.0
    ./scripts/get-coreos alpha 942.0.0

To reference local assets, change the `kernel` and `initrd` in a boot config file. For example, change `http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz` to `/assets/coreos/835.9.0/coreos_production_pxe.vmlinuz`.

## Virtual and Physical Machine Guides

Next, setup a virtual machine network within libvirt or a baremetal machine network. Follow the [libvirt guide](virtual-hardware.md) or [baremetal guide](physical-hardware.md).
