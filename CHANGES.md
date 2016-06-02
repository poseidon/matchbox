# coreos-baremetal bootcfg

## Latest

* Allow Ignition 2.0.0 JSON and YAML template files
* Add/improve rkt, Docker, Kubernetes, and binary/systemd deployment docs
* Show `bootcfg` message at the home path `/`
* Fix http package log messages and increase request logging (#173)
* Log requests for bootcfg hosted assets (#214)
* Error when an Ignition/Cloud-config template is rendered with a machine Group which is missing a metadata value. Previously, missing values defaulted to "no value" (#210)
* Stop requiring Ignition templates to use file extensions (#176)

#### Examples

* Add self-hosted Kubernetes example (PXE boot or install to disk)
* Add CoreOS Torus distributed storage cluster example (PXE boot)
* Add `create-uefi` subcommand to `scripts/libvirt` for UEFI/GRUB testing
* Updated Kubernetes examples to v1.2.4
* Remove 8.8.8.8 from networkd example Ignition configs (#184)
* Fix a bug in the k8s example k8s-certs@.service file check (#156)
* Match machines by MAC address in examples to simplify networkd device matching (#209)

## v0.3.0 (2016-04-14)

#### Features

* Add server library package for implementing servers
* Add initial gRPC client/server and a CLI tool
    - Allow listing, viewing, and creating Groups and Profiles
* Add initial Grub net boot support examples
* Add detached OpenPGP signature endpoints (`.sig`)
* Document deployment as a binary with systemd
* Upgrade from Go 1.5.3 to Go 1.6.1 (#139)

#### Changes

* Profiles
    - Move Profiles to JSON files under `/var/lib/bootcfg/profiles`
    - Rename `Spec` to `Profile` (#104)
* Groups
    - Move Groups to JSON files under `/var/lib/bootcfg/groups`
    - Require Group metadata to be valid JSON
    - Rename Group field `spec` to `profile`
    - Rename Group field `require` to `selector` (#147)
* Allow asset serving to be disabled with `-assets-path=""` (#118)
* Allow `selector` key/value pairs to be used in Ignition and Cloud config templates (#64)
* Change default `-data-path` to `/var/lib/bootcfg` (#132)
* Change default `-assets-path` to `/var/lib/bootcfg/assets` (#132)
* Change the default assets download location to `examples/assets`
* Stop parsing Groups from the `-config` YAML file. Remove the flag.
* Remove HTTP `/spec/id` JSON endpoint

#### Examples

* Convert all Cloud-Configs to Ignition
* Kubernetes
    * Upgraded Kubernetes examples to v1.2.0 (#122)
    * Run Heapster service by default (#142)
    * Example multi-node Kubernetes cluster installed to disk
* Example multi-node etcd cluster installed to disk
* Example which PXE boots with or without a root partition
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
