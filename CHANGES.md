# Matchbox

Notable changes between releases.

## Latest

* Build multi-arch container images (amd64, arm64) ([#823](https://github.com/poseidon/matchbox/pull/823))

## v0.9.0

* Refresh docs and examples for Fedora CoreOS and Flatcar Linux ([#815](https://github.com/poseidon/matchbox/pull/815), [#816](https://github.com/poseidon/matchbox/pull/816))
* Update Kubernetes manifest examples ([#791](https://github.com/poseidon/matchbox/pull/791), [#817](https://github.com/poseidon/matchbox/pull/817))
* Update Matchbox container image publishing ([#795](https://github.com/poseidon/matchbox/pull/795))
  * Publish Matchbox images from internal infra to Quay (`quay.io/poseidon/matchbox`)
  * Update Go version from v1.13.4 to v1.14.9
  * Update base image from `alpine:3.10` to `alpine:3.12` ([#784](https://github.com/poseidon/matchbox/pull/784))
* Include `contrib/k8s` in release tarballs ([#788](https://github.com/poseidon/matchbox/pull/788))
* Remove outdated systemd units ([#817](https://github.com/poseidon/matchbox/pull/817))
* Remove RPM spec file (Copr publishing stopped in v0.6)

## v0.8.3

* Publish docs to [https://matchbox.psdn.io](https://matchbox.psdn.io/) ([#769](https://github.com/poseidon/matchbox/pull/769))
* Update Go version from v1.11.7 to v1.13.4 ([#766](https://github.com/poseidon/matchbox/pull/766), [#770](https://github.com/poseidon/matchbox/pull/770))
* Update container image base from `alpine:3.9` to `alpine:3.10` ([#761](https://github.com/poseidon/matchbox/pull/761))
* Include `get-fedora-coreos` convenience script ([#763](https://github.com/poseidon/matchbox/pull/763))
* Remove Kubernetes provisioning examples ([#759](https://github.com/poseidon/matchbox/pull/759))
* Remove rkt tutorials and docs ([#765](https://github.com/poseidon/matchbox/pull/765))

## v0.8.1 - v0.8.2

Releases `v0.8.1` and `v0.8.2` were not built cleanly

* Release tags and container images have been removed
* Caused by go get golint (module-aware) mutating `go.mod` on Travis (see [#775](https://github.com/poseidon/matchbox/pull/775), [#777](https://github.com/poseidon/matchbox/pull/777))

## v0.8.0

* Transfer Matchbox repo from coreos to poseidon GitHub Org
* Publish container images at [quay.io/poseidon/matchbox](https://quay.io/repository/poseidon/matchbox)
* Build Matchbox with Go v1.11.7 for images and binaries
* Update container image base from alpine:3.6 to alpine:3.9
* Render Container Linux Configs as Ignition v2.2.0
* Validate raw Ignition configs with the v2.2 spec (warn-only)
  * Fix warnings that v2.2 configs are too new

Note: Release signing key [has changed](https://github.com/poseidon/matchbox/blob/v0.8.0/Documentation/deployment.md) with the project move.

### Examples

* Update Kubernetes example clusters to v1.14.1 (Terraform-based)

## v0.7.1 (2018-11-01)

* Add `kernel_args` variable to the terraform bootkube-install cluster definition
* Add `get-flatcar` helper script
* Add optional TLS support to read-only HTTP API
* Build Matchbox with Go 1.11.1 for images and binaries

### Examples

* Upgrade Kubernetes example clusters to v1.10.0 (Terraform-based)
* Upgrade Kubernetes example clusters to v1.8.5

## v0.7.0 (2017-12-12)

* Add gRPC API endpoints for managing generic (experimental) templates
* Update Container Linux config transpiler to v0.5.0
* Update Ignition to v0.19.0, render v2.1.0 Ignition configs
* Drop support for Container Linux versions below 1465.0.0 (breaking)
* Build Matchbox with Go 1.8.5 for images and binaries
* Remove Profile `Cmdline` map (deprecated in v0.5.0), use `Args` slice instead
* Remove pixiecore support (deprecated in v0.5.0)
* Remove `ContextHandler`, `ContextHandlerFunc`, and `NewHandler` from the `matchbox/http` package.

### Examples / Modules

* Upgrade Kubernetes example clusters to v1.8.4
* Kubernetes examples clusters enable etcd TLS
* Deploy the Container Linux Update Operator (CLUO) to coordinate reboots of Container Linux nodes in Kubernetes clusters. See the cluster [addon docs](Documentation/cluster-addons.md).
* Kubernetes examples (terraform and non-terraform) mask locksmithd
* Terraform modules `bootkube` and `profiles` (Kubernetes) mask locksmithd

## v0.6.1 (2017-05-25)

* Improve the installation documentation
* Move examples/etc/matchbox/cert-gen to scripts/tls
* Build Matchbox with Go 1.8.3 for images and binaries

### Examples

* Upgrade self-hosted Kubernetes cluster examples to v1.6.4
* Add NoSchedule taint to self-hosted Kubernetes controllers
* Remove static Kubernetes and rktnetes cluster examples

## v0.6.0 (2017-04-25)

* New [terraform-provider-matchbox](https://github.com/coreos/terraform-provider-matchbox) plugin for Terraform users!
* New hosted [documentation](https://coreos.com/matchbox/docs/latest) on coreos.com
* Add `ProfileDelete`, `GroupDelete`, `IgnitionGet` and `IgnitionDelete` gRPC endpoints
* Build matchbox with Go 1.8 for container images and binaries
* Generate code with gRPC v1.2.1 and matching Go protoc-gen-go plugin
* Update Ignition to v0.14.0 and coreos-cloudinit to v1.13.0
* Update "fuze" docs to the new name [Container Linux Configs](https://coreos.com/os/docs/latest/configuration.html)
* Remove `bootcmd` binary from release tarballs

### Examples

* Upgrade Kubernetes v1.5.5 (static) example clusters
* Upgrade Kubernetes v1.6.1 (self-hosted) example cluster
* Use etcd3 by default in all clusters (remove etcd2 clusters)
* Add Terraform examples for etcd3 and self-hosted Kubernetes 1.6.1

## v0.5.0 (2017-01-23)

* Rename project to CoreOS `matchbox`!
* Add Profile `args` field to list kernel args
* Update [Fuze](https://github.com/coreos/container-linux-config-transpiler) and [Ignition](https://github.com/coreos/ignition) to v0.11.2
* Switch from `golang.org/x/net/context` to `context`
* Deprecate Profile `cmd` field map of kernel args
* Deprecate Pixiecore support
* Drop build support for Go 1.6

#### Rename

* Move repo from github.com/coreos/coreos-baremetal to github.com/coreos/matchbox
* Rename `bootcfg` binary to `matchbox`
* Rename `bootcfg` packages to `matchbox`
* Publish a `quay.io/coreos/matchbox` container image. The `quay.io/coreos/bootcfg` image will no longer be updated.
* Rename environment variable prefix from `BOOTCFG*` to `MATCHBOX*`
* Change config directory to `/etc/matchbox`
* Change default `-data-path` to `/var/lib/matchbox`
* Change default `-assets-path` to `/var/lib/matchbox/assets`

#### Examples

* Upgrade Kubernetes v1.5.1 (static) example clusters
* Upgrade Kubernetes v1.5.1 (self-hosted) example cluster
* Switch Kubernetes (self-hosted) to run flannel as pods
* Combine rktnetes Ignition into Kubernetes static cluster

#### Migration

* binary users should install the `matchbox` binary (see [installation](Documentation/deployment.md))
* rkt/docker users should start using `quay.io/coreos/matchbox` (see [installation](Documentation/deployment.md))
* RPM users should uninstall bootcfg and install matchbox (see [installation](Documentation/deployment.md))
* Move `/etc/bootcfg` configs and certificates to `/etc/matchbox`
* Move `/var/lib/bootcfg` data to `/var/lib/matchbox`
* See the new [contrib/systemd](contrib/systemd) service examples
* Remove the old `bootcfg` user if you created one

## v0.4.2 (2016-12-7)

#### Improvements

* Add RPM packages to Copr
* Fix packaged `contrib/systemd` units
* Update Go version to 1.7.4

#### Examples

* Upgrade Kubernetes v1.4.6 (static manifest) example clusters
* Upgrade Kubernetes v1.4.6 (rktnetes) example clusters
* Upgrade Kubernetes v1.4.6 (self-hosted) example cluster

## v0.4.1 (2016-10-17)

#### Improvements

* Add ARM and ARM64 release architectures (#309)
* Add guide for installing bootcfg on CoreOS (#306)
* Improvements to the bootcfg cert-gen script (#310)

#### Examples

* Add Kubernetes example with rkt container runtime (i.e. rktnetes)
* Upgrade Kubernetes v1.4.1 (static manifest) example clusters
* Upgrade Kubernetes v1.4.1 (rktnetes) example clusters
* Upgrade Kubernetes v1.4.1 (self-hosted) example cluster
* Add etcd3 example cluster (PXE in-RAM or install to disk)
* Use DNS names (instead of IPs) in example clusters (except bootkube)

## v0.4.0 (2016-07-21)

#### Features

* Add/improve rkt, Docker, Kubernetes, and binary/systemd deployment docs
* TLS Client Authentication:
    * Add gRPC API TLS and TLS client-to-server authentication (#140)
    * Enable gRPC API by providing a TLS server `-cert-file` and `-key-file`, and a `-ca-file` to authenticate client certificates
    * Provide the `bootcmd` tool a TLS client `-cert-file` and `-key-file`, and a `-ca-file` to verify the server identity.
* Improvements to Ignition Support:
    * Allow Fuze YAML template files for Ignition 2.0.0 (#141)
    * Stop requiring Ignition templates to use file extensions (#176)
* Logging Improvements:
    * Add structured logging with Logrus (#254, #268)
    * Log requests for bootcfg assets (#214)
    * Show `bootcfg` message at the home path `/`
    * Fix http package log messages (#173)
* Templating:
    * Allow query parameters to be used as template variables as `{{.request.query.foo}}` (#182)
    * Support nested maps in responses from the "env file" metadata endpoint (#84)
    * Error when a template is rendered with variables which are missing a referenced key. Previously, missing lookups defaulted to "no value" (#210)
* gRPC API
    * Add DialTimeout to gRPC client config (#273)
    * Add IgnitionPut and Close to the client (#160,#193)

#### Changes

* gRPC API requires TLS client authentication
* Replace Ignition YAML templates with Fuze templates
    - Fuze formalizes the transform from Fuze configs (YAML) to Ignition 2.0.0 (JSON)
    - [Migrate templates from v0.3.0](Documentation/ignition.md#migration-from-v030)
    - Require CoreOS 1010.1.0 or newer
    - Drop support for Ignition v1 format
* Replace template variable `{{.query}}` with `{{.request.raw_query}}`

#### Examples

* Kubernetes
    * Upgrade Kubernetes v1.3.0 (static manifest) example clusters
    * Add Kubernetes v1.3.0-beta.2 (self-hosted) example cluster
    * Mount /etc/resolv.conf into host kubelet for skydns and pod DNS lookups (#237,#260)
    * Fix a bug in the k8s example k8s-certs@.service file check (#156)
    * Avoid systemd dependency failures by restarting components (#257,#274)
    * Verify Kubernetes v1.2.4 and v1.3.0 clusters pass conformance tests (#71,#265)
* Add Torus distributed storage cluster example (PXE boot)
* Add `create-uefi` subcommand to `scripts/libvirt` for UEFI/GRUB testing
* Install CoreOS to disk from a cached copy via bootcfg baseurl (#228)
* Remove 8.8.8.8 from networkd example Ignition configs (#184)
* Match machines by MAC address in examples to simplify networkd device matching (#209)
* With rkt 1.8+, you can use `rkt gc --grace-period=0` to cleanup rkt IP assignments in examples. The `rkt-gc-force` script has been removed.

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
