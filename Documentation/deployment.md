
# Installation

This guide walks through deploying the `matchbox` service on a Linux host (via RPM, rkt, docker, or binary) or on a Kubernetes cluster.

## Provisoner

`matchbox` is a service for network booting and provisioning machines to create CoreOS clusters. `matchbox` should be installed on a provisioner machine (CoreOS or any Linux distribution) or cluster (Kubernetes) which can serve configs to client machines in a lab or datacenter.

Choose one of the supported installation options:

* [CoreOS (rkt)](#coreos)
* [RPM-based](#rpm-based-distro)
* [General Linux (binary)](#general-linux)
* [With rkt](#rkt)
* [With docker](#docker)
* [Kubernetes Service](#kubernetes)

## Download

Download the latest matchbox [release](https://github.com/coreos/matchbox/releases) to the provisioner host.

```sh
$ wget https://github.com/coreos/matchbox/releases/download/v0.5.0/matchbox-v0.5.0-linux-amd64.tar.gz
$ wget https://github.com/coreos/matchbox/releases/download/v0.5.0/matchbox-v0.5.0-linux-amd64.tar.gz.asc
```

Verify the release has been signed by the [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/).

```
$ gpg --keyserver pgp.mit.edu --recv-key 18AD5014C99EF7E3BA5F6CE950BDD3E0FC8A365E
$ gpg --verify matchbox-v0.5.0-linux-amd64.tar.gz.asc matchbox-v0.5.0-linux-amd64.tar.gz
# gpg: Good signature from "CoreOS Application Signing Key <security@coreos.com>"
```

Untar the release.

```sh
$ tar xzvf matchbox-v0.5.0-linux-amd64.tar.gz
$ cd matchbox-v0.5.0-linux-amd64
```

## Install

### RPM-based Distro

On an RPM-based provisioner, install the `matchbox` RPM from the Copr [repository](https://copr.fedorainfracloud.org/coprs/g/CoreOS/matchbox/) using `dnf` or `yum`.

```sh
dnf copr enable @CoreOS/matchbox
dnf install matchbox
```

### CoreOS

On a CoreOS provisioner, rkt run `matchbox` image with the provided systemd unit.

```sh
$ sudo cp contrib/systemd/matchbox-on-coreos.service /etc/systemd/system/matchbox.service
```

### General Linux

Pre-built binaries are available for general Linux distributions. Copy the `matchbox` static binary to an appropriate location on the host.

```sh
$ sudo cp matchbox /usr/local/bin
```

#### Set Up User/Group

The `matchbox` service should be run by a non-root user with access to the `matchbox` data directory (`/var/lib/matchbox`). Create a `matchbox` user and group.

```sh
$ sudo useradd -U matchbox
$ sudo mkdir -p /var/lib/matchbox/assets
$ sudo chown -R matchbox:matchbox /var/lib/matchbox
```

#### Create systemd Service

Copy the provided `matchbox` systemd unit file.

```sh
$ sudo cp contrib/systemd/matchbox-local.service /etc/systemd/system/
```

## Customization

Customize matchbox by editing the systemd unit or adding a systemd dropin. Find the complete set of `matchbox` flags and environment variables at [config](config.md).

    sudo systemctl edit matchbox

By default, the read-only HTTP machine endpoint will be exposed on port **8080**.

```ini
# /etc/systemd/system/matchbox.service.d/override.conf
[Service]
Environment="MATCHBOX_ADDRESS=0.0.0.0:8080"
Environment="MATCHBOX_LOG_LEVEL=debug"
```

A common customization is enabling the gRPC API to allow clients with a TLS client certificate to change machine configs.

```ini
# /etc/systemd/system/matchbox.service.d/override.conf
[Service]
Environment="MATCHBOX_ADDRESS=0.0.0.0:8080"
Environment="MATCHBOX_RPC_ADDRESS=0.0.0.0:8081"
```

The Tectonic [Installer](https://tectonic.com/enterprise/docs/latest/install/bare-metal/index.html) uses this API. Tectonic users with a CoreOS provisioner can start with an example that enables it.

```sh
$ sudo cp contrib/systemd/matchbox-for-tectonic.service /etc/systemd/system/matchbox.service
```

Customize `matchbox` to suit your preferences.

## Firewall

Allow your port choices on the provisioner's firewall so the clients can access the service. Here are the commands for those using `firewalld`:

```sh
$ sudo firewall-cmd --zone=MYZONE --add-port=8080/tcp --permanent
$ sudo firewall-cmd --zone=MYZONE --add-port=8081/tcp --permanent
```

## Generate TLS Credentials

*Skip this unless you need to enable the gRPC API*

The `matchbox` gRPC API allows client apps (`bootcmd` CLI, Tectonic Installer, etc.) to update how machines are provisioned. TLS credentials are needed for client authentication and to establish a secure communication channel. Client machines (those PXE booting) read from the HTTP endpoints and do not require this setup.

If your organization manages public key infrastructure and a certificate authority, create a server certificate and key for the `matchbox` service and a client certificate and key for each client tool.

Otherwise, generate a self-signed `ca.crt`, a server certificate  (`server.crt`, `server.key`), and client credentials (`client.crt`, `client.key`) with the `examples/etc/matchbox/cert-gen` script. Export the DNS name or IP (discouraged) of the provisioner host.

```sh
$ cd examples/etc/matchbox
# DNS or IP Subject Alt Names where matchbox can be reached
$ export SAN=DNS.1:matchbox.example.com,IP.1:192.168.1.42
$ ./cert-gen
```

Place the TLS credentials in the default location:

```sh
$ sudo mkdir -p /etc/matchbox
$ sudo cp ca.crt server.crt server.key /etc/matchbox/
```

Save `client.crt`, `client.key`, and `ca.crt` to use with a client tool later.

## Start matchbox

Start the `matchbox` service and enable it if you'd like it to start on every boot.

```sh
$ sudo systemctl daemon-reload
$ sudo systemctl start matchbox
$ sudo systemctl enable matchbox
```

## Verify

Verify the matchbox service is running and can be reached by client machines (those being provisioned).

```sh
$ systemctl status matchbox
$ dig matchbox.example.com
```

Verify you receive a response from the HTTP and API endpoints.

```sh
$ curl http://matchbox.example.com:8080
matchbox
```

If you enabled the gRPC API,

```sh
$ openssl s_client -connect matchbox.example.com:8081 -CAfile /etc/matchbox/ca.crt -cert examples/etc/matchbox/client.crt -key examples/etc/matchbox/client.key
CONNECTED(00000003)
depth=1 CN = fake-ca
verify return:1
depth=0 CN = fake-server
verify return:1
---
Certificate chain
 0 s:/CN=fake-server
   i:/CN=fake-ca
---
....
```

## Download CoreOS (optional)

`matchbox` can serve CoreOS images in development or lab environments to reduce bandwidth usage and increase the speed of CoreOS PXE boots and installs to disk.

Download a recent CoreOS [release](https://coreos.com/releases/) with signatures.

```sh
$ ./scripts/get-coreos stable 1235.9.0 .     # note the "." 3rd argument
```

Move the images to `/var/lib/matchbox/assets`,

```sh
$ sudo cp -r coreos /var/lib/matchbox/assets
```

```
/var/lib/matchbox/assets/
├── coreos
│   └── 1235.9.0
│       ├── CoreOS_Image_Signing_Key.asc
│       ├── coreos_production_image.bin.bz2
│       ├── coreos_production_image.bin.bz2.sig
│       ├── coreos_production_pxe_image.cpio.gz
│       ├── coreos_production_pxe_image.cpio.gz.sig
│       ├── coreos_production_pxe.vmlinuz
│       └── coreos_production_pxe.vmlinuz.sig
```

and verify the images are acessible.

```
$ curl http://matchbox.example.com:8080/assets/coreos/1235.9.0/
<pre>...
```

For large production environments, use a cache proxy or mirror suitable for your environment to serve CoreOS images.

## Network

Review [network setup](https://github.com/coreos/matchbox/blob/master/Documentation/network-setup.md) with your network administrator to set up DHCP, TFTP, and DNS services on your network. At a high level, your goals are to:

* Chainload PXE firmwares to iPXE
* Point iPXE client machines to the `matchbox` iPXE HTTP endpoint `http://matchbox.example.com:8080/boot.ipxe`
* Ensure `matchbox.example.com` resolves to your `matchbox` deployment

CoreOS provides [dnsmasq](https://github.com/coreos/matchbox/tree/master/contrib/dnsmasq) as `quay.io/coreos/dnsmasq`, if you wish to use rkt or Docker.

## rkt

Run the container image with rkt.

latest or most recent tagged `matchbox` [release](https://github.com/coreos/matchbox/releases) ACI. Trust the [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/) for image signature verification.

```sh
$ sudo rkt run --net=host --mount volume=data,target=/var/lib/matchbox --volume data,kind=host,source=/var/lib/matchbox quay.io/coreos/matchbox:latest --mount volume=config,target=/etc/matchbox --volume config,kind=host,source=/etc/matchbox,readOnly=true -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
```

Create machine profiles, groups, or Ignition configs at runtime with `bootcmd` or by using your own `/var/lib/matchbox` volume mounts.

## Docker

Run the container image with docker.

```sh
sudo docker run --net=host --rm -v /var/lib/matchbox:/var/lib/matchbox:Z -v /etc/matchbox:/etc/matchbox:Z,ro quay.io/coreos/matchbox:latest -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
```

Create machine profiles, groups, or Ignition configs at runtime with `bootcmd` or by using your own `/var/lib/matchbox` volume mounts.

## Kubernetes

Create a `matchbox` Kubernetes `Deployment` and `Service` based on the example manifests provided in [contrib/k8s](../contrib/k8s).

```
$ kubectl apply -f contrib/k8s/matchbox-deployment.yaml
$ kubectl apply -f contrib/k8s/matchbox-service.yaml
```

This runs the `matchbox` service exposed on NodePort `tcp:31488` on each node in the cluster. `MATCHBOX_LOG_LEVEL` is set to debug.

```sh
$ kubectl get deployments
$ kubectl get services
$ kubectl get pods
$ kubectl logs POD-NAME
```

The example manifests use Kubernetes `emptyDir` volumes to back the `matchbox` FileStore (`/var/lib/matchbox`). This doesn't provide long-term persistent storage so you may wish to mount your machine groups, profiles, and Ignition configs with a [gitRepo](http://kubernetes.io/docs/user-guide/volumes/#gitrepo) and host image assets on a file server.

### Documentation

View the [documentation](https://github.com/coreos/matchbox#coreos-on-baremetal) for `matchbox` service docs, tutorials, example clusters and Ignition configs, PXE booting guides, or machine lifecycle guides.
