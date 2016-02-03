# coreos-baremetal Config Service

## Latest

* Render Ignition config and Cloud configs as Go templates (favor systemd EnvironmentFile where possible).
* Add a `metadata` endpoint so machine instances can fetch their metadata.
* Allow `metadata` to be added to group definitions in config.yaml
* Add detached OpenPGP signature endpoints (suffix `.asc`) for all configs.
    - Enable signing by providing a `-key-ring-path` with a signing key and setting `BOOTCFG_PASSPHRASE` if needed.
* Require the `-config` flag if the default file path doesn't exist
* Normalize user-defined MAC address tags
* Rename flag `-images-path` to `-assets-path`
* Rename endpoint `/images` to `/assets`
* Example TLS-authenticated Kubernetes cluster with rkt
* Example TLS-authenticated Kubernetes cluster with Docker
* Example custom metadata agent with Ignition, fetches on boot

## v0.1.0 (2015-01-08)

Initial release of the coreos-baremetal Config Service.

### Features

* Support for PXE, iPXE, and Pixiecore network boot environments
* Match machines based on hardware attributes or free-form tag matchers
* Render boot configs (kernel, initrd), [Ignition](https://coreos.com/ignition/docs/latest/what-is-ignition.html) configs, and [Cloud-Init](https://github.com/coreos/coreos-cloudinit) configs
