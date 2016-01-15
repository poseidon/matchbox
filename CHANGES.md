# coreos-baremetal Config Service

## Latest

* Require the `-config` flag if the default file path doesn't exist
* Normalize user-defined MAC address tags
* Renamed flag `-images-path` to `-assets-path`
* Renamed endpoint `/images` to `/assets`

## v0.1.0 (2015-01-08)

Initial release of the coreos-baremetal Config Service.

### Features

* Support for PXE, iPXE, and Pixiecore network boot environments
* Match machines based on hardware attributes or free-form tag matchers
* Render boot configs (kernel, initrd), [Ignition](https://coreos.com/ignition/docs/latest/what-is-ignition.html) configs, and [Cloud-Init](https://github.com/coreos/coreos-cloudinit) configs
