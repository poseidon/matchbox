# coreos-baremetal bootcfg

## Latest

* Add initial gRPC client and server packages
* Add initial Grub net boot support and an example
* Add initial command line client tool
* Add detached OpenPGP signature endpoints (`.sig`)

#### Changes

* Profiles
    - Move Profiles to JSON files under `/var/lib/bootcfg/profiles`
    - Rename `Spec` to `Profile`
* Groups
    - Move Groups to JSON files under `/var/lib/bootcfg/groups`
    - Require Group metadata to be valid JSON
    - Rename groups field `spec` to `profile`
* Stop parsing Groups from the `-config` YAML file. Remove the flag.
* Change default `-data-path` to `/var/lib/bootcfg`
* Change default `-assets-path` to `/var/lib/bootcfg/assets`
* Change the default assets download location to `examples/assets`
* Remove HTTP `/spec/id` JSON endpoint

#### Examples

* Convert all Cloud-Configs to Ignition
* Kubernetes
    * Upgraded Kubernetes examples to v1.2.0
    * Add example to install Kubernetes to disk
    * Run Heapster service by default
* Examples which PXE boot with or without a root partition
* Example etcd cluster installed to disk
* Setup fleet in multi-node example clusters


## v0.2.0 (2016-02-09)

#### Features

* Render Ignition config and cloud-configs as Go templates
* Allow writing Ignition configs as YAML configs. Render as JSON for machines.
* Add ASCII armored detached OpenPGP signature endpoints (`.asc`)
    - Enable signing by providing a `-key-ring-path` with a signing key and setting `BOOTCFG_PASSPHRASE` if needed
* Add `metadata` endpoint which matches machines to custom metadata
* Add `metadata` to group definitions in `config.yaml`

#### Changes

* Require the `-config` flag if the default file path doesn't exist
* Normalize user-defined MAC address tags
* Rename flag `-images-path` to `-assets-path`
* Rename endpoint `/images` to `/assets`

#### New Examples

* Example TLS-authenticated Kubernetes cluster with rkt and CNI
* Example TLS-authenticated Kubernetes cluster with Docker
* Example custom metadata agent with Ignition, fetches metadata on boot and writes it to `/run/metadata/bootcfg`
* Example CoreOS install to disk with Ignition
* Update etcd cluster examples to use Ignition, rather than cloud-config.

## v0.1.0 (2016-01-08)

Initial release of the coreos-baremetal Config Service.

#### Features

* Match machines based on hardware attributes or free-form tag matchers
* Render boot configs (kernel, initrd), [Ignition](https://coreos.com/ignition/docs/latest/what-is-ignition.html) configs, and [Cloud-Init](https://github.com/coreos/coreos-cloudinit) configs
* Support for PXE, iPXE, and Pixiecore network boot environments
