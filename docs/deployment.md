# Installation

This guide walks through deploying the `matchbox` service on a Linux host (as a binary or container image) or on a Kubernetes cluster.

## Provisoner

Matchbox is a service for network booting and provisioning machines to create Fedora CoreOS or Flatcar Linux clusters. Matchbox may installed on a host server or Kubernetes cluster that can serve configs to client machines in a lab or datacenter.

Choose one of the supported installation options:

* [Matchbox binary](#matchbox-binary)
* [Container image](#container-image)
* [Kubernetes manifests](#kubernetes)

## Download

Download the latest Matchbox [release](https://github.com/poseidon/matchbox/releases).

```sh
$ wget https://github.com/poseidon/matchbox/releases/download/v0.10.0/matchbox-v0.10.0-linux-amd64.tar.gz
$ wget https://github.com/poseidon/matchbox/releases/download/v0.10.0/matchbox-v0.10.0-linux-amd64.tar.gz.asc
```

Verify the release has been signed by Dalton Hubble's GPG [Key](https://keyserver.ubuntu.com/pks/lookup?search=0x8F515AD1602065C8&op=vindex)'s signing subkey.

```sh
$ gpg --keyserver keyserver.ubuntu.com --recv-key 2E3D92BF07D9DDCCB3BAE4A48F515AD1602065C8
$ gpg --verify matchbox-v0.10.0-linux-amd64.tar.gz.asc matchbox-v0.10.0-linux-amd64.tar.gz
gpg: Good signature from "Dalton Hubble <dghubble@gmail.com>"
```

Untar the release.

```sh
$ tar xzvf matchbox-v0.10.0-linux-amd64.tar.gz
$ cd matchbox-v0.10.0-linux-amd64
```

## Install

Run Matchbox as a binary, a container image, or on Kubernetes.

### Matchbox Binary

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
$ sudo cp contrib/systemd/matchbox.service /etc/systemd/system/matchbox.service
```

#### systemd dropins

Customize Matchbox by editing the systemd unit or adding a systemd dropin. Find the complete set of `matchbox` flags and environment variables at [config](config.md).

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

Customize `matchbox` to suit your preferences.

#### Start

Start the Matchbox service and enable it if you'd like it to start on every boot.

```
$ sudo systemctl daemon-reload
$ sudo systemctl start matchbox
$ sudo systemctl enable matchbox
```

### Container Image

Run the container image with Podman,

```
mkdir -p /var/lib/matchbox/assets
podman run --net=host --rm -v /var/lib/matchbox:/var/lib/matchbox:Z -v /etc/matchbox:/etc/matchbox:Z,ro quay.io/poseidon/matchbox:v0.10.0 -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
```

Or with Docker,

```
mkdir -p /var/lib/matchbox/assets
sudo docker run --net=host --rm -v /var/lib/matchbox:/var/lib/matchbox:Z -v /etc/matchbox:/etc/matchbox:Z,ro quay.io/poseidon/matchbox:v0.10.0 -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
```

Create machine profiles, groups, or Ignition configs by adding files to `/var/lib/matchbox`.

### Kubernetes

Install Matchbox on a Kubernetes cluster with the example manifests.

```sh
$ kubectl apply -R -f contrib/k8s
$ kubectl get services
NAME                 CLUSTER-IP   EXTERNAL-IP   PORT(S)             AGE
matchbox             10.3.0.145   <none>        8080/TCP,8081/TCP   46m
```

Example manifests in [contrib/k8s](../contrib/k8s) enable the gRPC API to allow client apps to update matchbox objects. Generate TLS server certificates for `matchbox-rpc.example.com` [as shown](#generate-tls-certificates) and create a Kubernetes secret. Alternately, edit the example manifests if you don't need the gRPC API enabled.

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
$ export SAN=DNS.1:matchbox.example.com,IP.1:172.17.0.2
```

Generate a `ca.crt`, `server.crt`, `server.key`, `client.crt`, and `client.key`.

```sh
$ ./cert-gen
```

Move TLS credentials to the matchbox server's default location.

```sh
$ sudo mkdir -p /etc/matchbox
$ sudo cp ca.crt server.crt server.key /etc/matchbox
$ sudo chown -R matchbox:matchbox /etc/matchbox
```

Save `client.crt`, `client.key`, and `ca.crt` for later use (e.g. `~/.matchbox`).

```sh
$ mkdir -p ~/.matchbox
$ cp client.crt client.key ca.crt ~/.matchbox/
```

## Verify

Verify the matchbox service is running and can be reached by client machines (those being provisioned).

```sh
$ systemctl status matchbox   # Matchbox binary method
$ dig matchbox.example.com
```

Verify you receive a response from the HTTP and API endpoints.

```sh
$ curl http://matchbox.example.com:8080
matchbox
```

If you enabled the gRPC API,

```sh
$ openssl s_client -connect matchbox.example.com:8081 -CAfile scripts/tls/ca.crt -cert scripts/tls/client.crt -key scripts/tls/client.key
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

## Download Images (optional)

Matchbox can serve OS images in development or lab environments to reduce bandwidth usage and increase the speed of PXE boots and installs to disk.

Download a recent Fedora CoreOS or Flatcar Linux release.

```
$ ./scripts/get-fedora-coreos stable 36.20220906.3.2 .
$ ./scripts/get-flatcar stable 3227.2.0 .
```

Move the images to `/var/lib/matchbox/assets`,

```
/var/lib/matchbox/assets/fedora-coreos/
├── fedora-coreos-36.20220906.3.2-live-initramfs.x86_64.img
├── fedora-coreos-36.20220906.3.2-live-kernel-x86_64
├── fedora-coreos-36.20220906.3.2-live-rootfs.x86_64.img

/var/lib/matchbox/assets/flatcar/
└── 3227.2.0
    ├── Flatcar_Image_Signing_Key.asc
    ├── flatcar_production_image.bin.bz2
    ├── flatcar_production_image.bin.bz2.sig
    ├── flatcar_production_pxe_image.cpio.gz
    ├── flatcar_production_pxe_image.cpio.gz.sig
    ├── flatcar_production_pxe.vmlinuz
    ├── flatcar_production_pxe.vmlinuz.sig
    └── version.txt
```

and verify the images are accessible.

```sh
$ curl http://matchbox.example.com:8080/assets/fedora-coreos/
<pre>...
```

For large production environments, use a cache proxy or mirror suitable for your environment to serve images.

## Network

Review [network setup](https://github.com/poseidon/matchbox/blob/master/docs/network-setup.md) with your network administrator to set up DHCP, TFTP, and DNS services on your network. At a high level, your goals are to:

* Chainload PXE firmwares to iPXE
* Point iPXE client machines to the `matchbox` iPXE HTTP endpoint `http://matchbox.example.com:8080/boot.ipxe`
* Ensure `matchbox.example.com` resolves to your `matchbox` deployment

Poseidon provides [dnsmasq](https://github.com/poseidon/matchbox/tree/master/contrib/dnsmasq) as `quay.io/poseidon/dnsmasq`.

# TLS

Matchbox can serve the read-only HTTP API with TLS.

| Name           | Type   | Description |
|----------------|--------|-------------|
| -web-ssl       | bool   | true/false  |
| -web-cert-file | string | Path to the server TLS certificate file |
| -web-key-file  | string | Path to the server TLS key file |

However, it is more common to use an Ingress Controller (Kubernetes) to terminate TLS.

### Operational notes

* Secrets: Matchbox **can** be run as a public facing service. However, you **must** follow best practices and avoid writing secret material into machine user-data. Instead, load secret materials from an internal secret store.
* Storage: Example manifests use Kubernetes `emptyDir` volumes to store `matchbox` data. Swap those out for a Kubernetes persistent volume if available.

