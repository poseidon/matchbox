# coreos-baremetal bootcfg

## Latest

## v0.2.0 (2016-02-09)

#### Features

* Render Ignition config and cloud-configs as Go templates
* Allow writing Ignition configs as YAML configs. Render as JSON for machines.
* Add detached OpenPGP signature endpoints (`.asc`) for all configs.
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
