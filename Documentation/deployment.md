# Installation

This guide walks through deploying the `matchbox` service on a Linux host (via RPM, rkt, docker, or binary) or on a Kubernetes cluster.

## Provisoner

`matchbox` is a service for network booting and provisioning machines to create CoreOS Container Linux clusters. `matchbox` should be installed on a provisioner machine (Container Linux or any Linux distribution) or cluster (Kubernetes) which can serve configs to client machines in a lab or datacenter.

Choose one of the supported installation options:

* [CoreOS Container Linux (rkt)](#coreos-container-linux)
* [RPM-based](#rpm-based-distro)
* [Generic Linux (binary)](#generic-linux)
* [With rkt](#rkt)
* [With docker](#docker)
* [Kubernetes Service](#kubernetes)

## Download

Download the latest matchbox [release](https://github.com/coreos/matchbox/releases) to the provisioner host.

```sh
$ wget https://github.com/coreos/matchbox/releases/download/v0.6.1/matchbox-v0.6.1-linux-amd64.tar.gz
$ wget https://github.com/coreos/matchbox/releases/download/v0.6.1/matchbox-v0.6.1-linux-amd64.tar.gz.asc
```

Verify the release has been signed by the [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/).

```sh
$ gpg --keyserver pgp.mit.edu --recv-key 18AD5014C99EF7E3BA5F6CE950BDD3E0FC8A365E
$ gpg --verify matchbox-v0.6.1-linux-amd64.tar.gz.asc matchbox-v0.6.1-linux-amd64.tar.gz
# gpg: Good signature from "CoreOS Application Signing Key <security@coreos.com>"
```

Untar the release.

```sh
$ tar xzvf matchbox-v0.6.1-linux-amd64.tar.gz
$ cd matchbox-v0.6.1-linux-amd64
```

## Install

### RPM-based distro

On an RPM-based provisioner (Fedora 24+), install the `matchbox` RPM from the Copr [repository](https://copr.fedorainfracloud.org/coprs/g/CoreOS/matchbox/) using `dnf`.

```sh
dnf copr enable @CoreOS/matchbox
dnf install matchbox
```

RPMs are not currently available for CentOS and RHEL (due to Go version). CentOS and RHEL users should follow the Generic Linux section below.

### CoreOS Container Linux

On a Container Linux provisioner, rkt run `matchbox` image with the provided systemd unit.

```sh
$ sudo cp contrib/systemd/matchbox-on-coreos.service /etc/systemd/system/matchbox.service
```

### Generic Linux

Pre-built binaries are available for generic Linux distributions. Copy the `matchbox` static binary to an appropriate location on the host.

```sh
$ sudo cp matchbox /usr/local/bin
```

#### Set up User/Group

The `matchbox` service should be run by a non-root user with access to the `matchbox` data directory (`/var/lib/matchbox`). Create a `matchbox` user and group.

```sh
$ sudo useradd -U matchbox
$ sudo mkdir -p /var/lib/matchbox/assets
$ sudo chown -R matchbox:matchbox /var/lib/matchbox
```

#### Create systemd service

Copy the provided `matchbox` systemd unit file.

```sh
$ sudo cp contrib/systemd/matchbox-local.service /etc/systemd/system/matchbox.service
```

## Customization

Customize matchbox by editing the systemd unit or adding a systemd dropin. Find the complete set of `matchbox` flags and environment variables at [config](config.md).

```sh
$ sudo systemctl edit matchbox
```

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

The Tectonic [Installer](https://tectonic.com/enterprise/docs/latest/install/bare-metal/index.html) uses this API. Tectonic users with a Container Linux provisioner can start with an example that enables it.

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

## Generate TLS Certificates

The Matchbox gRPC API allows clients (terraform-provider-matchbox) to create and update Matchbox resources. TLS credentials are needed for client authentication and to establish a secure communication channel. Client machines (those PXE booting) read from the HTTP endpoints and do not require this setup.

The `cert-gen` helper script generates a self-signed CA, server certificate, and client certificate. **Prefer your organization's PKI, if possible**

Navigate to the `scripts/tls` directory.

```sh
$ cd scripts/tls
```

Export `SAN` to set the Subject Alt Names which should be used in certificates. Provide the fully qualified domain name or IP (discouraged) where Matchbox will be installed.

```sh
# DNS or IP Subject Alt Names where matchbox runs
$ export SAN=DNS.1:matchbox.example.com,IP.1:172.18.0.2
```

Generate a `ca.crt`, `server.crt`, `server.key`, `client.crt`, and `client.key`.

```sh
$ ./cert-gen
```

Move TLS credentials to the matchbox server's default location.

```sh
$ sudo mkdir -p /etc/matchbox
$ sudo cp ca.crt server.crt server.key /etc/matchbox
```

Save `client.crt`, `client.key`, and `ca.crt` for later use (e.g. `~/.matchbox`).

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
$ openssl s_client -connect matchbox.example.com:8081 -CAfile /etc/matchbox/ca.crt -cert scripts/tls/client.crt -key scripts/tls/client.key
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

## Download Container Linux (optional)

`matchbox` can serve Container Linux images in development or lab environments to reduce bandwidth usage and increase the speed of Container Linux PXE boots and installs to disk.

Download a recent Container Linux [release](https://coreos.com/releases/) with signatures.

```sh
$ ./scripts/get-coreos stable 1465.7.0 .     # note the "." 3rd argument
```

Move the images to `/var/lib/matchbox/assets`,

```sh
$ sudo cp -r coreos /var/lib/matchbox/assets
```

```
/var/lib/matchbox/assets/
├── coreos
│   └── 1465.7.0
│       ├── CoreOS_Image_Signing_Key.asc
│       ├── coreos_production_image.bin.bz2
│       ├── coreos_production_image.bin.bz2.sig
│       ├── coreos_production_pxe_image.cpio.gz
│       ├── coreos_production_pxe_image.cpio.gz.sig
│       ├── coreos_production_pxe.vmlinuz
│       └── coreos_production_pxe.vmlinuz.sig
```

and verify the images are acessible.

```sh
$ curl http://matchbox.example.com:8080/assets/coreos/1465.7.0/
<pre>...
```

For large production environments, use a cache proxy or mirror suitable for your environment to serve Container Linux images. See [contrib/squid](../contrib/squid/README.md) for details.

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
$ mkdir -p /var/lib/matchbox/assets
$ sudo rkt run --net=host --mount volume=data,target=/var/lib/matchbox --volume data,kind=host,source=/var/lib/matchbox quay.io/coreos/matchbox:latest --mount volume=config,target=/etc/matchbox --volume config,kind=host,source=/etc/matchbox,readOnly=true -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
```

Create machine profiles, groups, or Ignition configs by adding files to `/var/lib/matchbox`.

## Docker

Run the container image with docker.

```sh
$ mkdir -p /var/lib/matchbox/assets
$ sudo docker run --net=host --rm -v /var/lib/matchbox:/var/lib/matchbox:Z -v /etc/matchbox:/etc/matchbox:Z,ro quay.io/coreos/matchbox:latest -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
```

Create machine profiles, groups, or Ignition configs by adding files to `/var/lib/matchbox`.

## Kubernetes

Install `matchbox` on a Kubernetes cluster by creating a deployment and service.

```sh
$ kubectl apply -f contrib/k8s/matchbox-deployment.yaml
$ kubectl apply -f contrib/k8s/matchbox-service.yaml
$ kubectl get services
NAME                 CLUSTER-IP   EXTERNAL-IP   PORT(S)             AGE
matchbox             10.3.0.145   <none>        8080/TCP,8081/TCP   46m
```

Example manifests in [contrib/k8s](../contrib/k8s) enable the gRPC API to allow client apps to update matchbox objects. Generate TLS server credentials for `matchbox-rpc.example.com` [as shown](#generate-tls-credentials) and create a Kubernetes secret. Alternately, edit the example manifests if you don't need the gRPC API enabled.

```sh
$ kubectl create secret generic matchbox-rpc --from-file=ca.crt --from-file=server.crt --from-file=server.key
```

Create an Ingress resource to expose the HTTP read-only and gRPC API endpoints. The Ingress example requires the cluster to have a functioning [Nginx Ingress Controller](https://github.com/kubernetes/ingress).

```sh
$ kubectl create -f contrib/k8s/matchbox-ingress.yaml
$ kubectl get ingress
NAME      HOSTS                                          ADDRESS            PORTS     AGE
matchbox       matchbox.example.com                      10.128.0.3,10...   80        29m
matchbox-rpc   matchbox-rpc.example.com                  10.128.0.3,10...   80, 443   29m
```

Add DNS records `matchbox.example.com` and `matchbox-rpc.example.com` to route traffic to the Ingress Controller.

Verify `http://matchbox.example.com` responds with the text "matchbox" and verify gRPC clients can connect to `matchbox-rpc.example.com:443`.

```sh
$ curl http://matchbox.example.com
$ openssl s_client -connect matchbox-rpc.example.com:443 -CAfile ca.crt -cert client.crt -key client.key
```

### Operational notes

* Secrets: Matchbox **can** be run as a public facing service. However, you **must** follow best practices and avoid writing secret material into machine user-data. Instead, load secret materials from an internal secret store.
* Storage: Example manifests use Kubernetes `emptyDir` volumes to store `matchbox` data. Swap those out for a Kubernetes persistent volume if available.

