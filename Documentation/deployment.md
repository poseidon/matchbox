
# Installation

This guide walks through deploying the `bootcfg` service on a Linux host (via binary, rkt, or docker) or on a Kubernetes cluster.

## Provisoner

`bootcfg` is a service for network booting and provisioning machines to create CoreOS clusters. If the gRPC API is enabled, TLS client-cert authenticated clients (such as `bootcmd` or graphical tools) can be used to programmatically manage machine profiles as well.

`bootcfg` should be deployed to a provisioner node (CoreOS or any Linux distribution) or cluster (Kubernetes) which can serve configurations to client machines in a lab or datacenter. Choose a deployment target and a DNS name (or IP) at which it will be accessed (e.g. bootcfg.example.com).

## Binary Release

Download the latest coreos-baremetal [release](https://github.com/coreos/coreos-baremetal/releases) to the provisioner host.

    wget https://github.com/coreos/coreos-baremetal/releases/download/v0.4.0/coreos-baremetal-v0.4.0-linux-amd64.tar.gz
    wget https://github.com/coreos/coreos-baremetal/releases/download/v0.4.0/coreos-baremetal-v0.4.0-linux-amd64.tar.gz.asc

Verify the release has been signed by the [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/).

    gpg --keyserver pgp.mit.edu --recv-key 18AD5014C99EF7E3BA5F6CE950BDD3E0FC8A365E
    gpg --verify coreos-baremetal-v0.4.0-linux-amd64.tar.gz.asc coreos-baremetal-v0.4.0-linux-amd64.tar.gz
    # gpg: Good signature from "CoreOS Application Signing Key <security@coreos.com>"

Untar the release.

    tar xzvf coreos-baremetal-v0.4.0-linux-amd64.tar.gz
    cd coreos-baremetal-v0.4.0-linux-amd64

### Copy Binary

Install the `bootcfg` static binary at an appropriate location on the host.

```sh
sudo cp bootcfg /usr/local/bin
```

### Set Up User/Group

The `bootcfg` service should be run by a non-root user with access to the `bootcfg` data directory (`/var/lib/bootcfg`). Create a `bootcfg` user and group.

```sh
$ sudo useradd -U bootcfg
$ sudo mkdir -p /var/lib/bootcfg/assets
$ sudo chown -R bootcfg:bootcfg /var/lib/bootcfg
```

### Create systemd Service

Copy the provided `bootcfg` systemd unit file. The example systemd unit exposes the `bootcfg` HTTP machine config endpoints on port 8080 and the (optional) gRPC API on port 8081.

```sh
sudo cp contrib/systemd/bootcfg.service /etc/systemd/system/
sudo systemctl daemon-reload
```

Customize the port settings to suit your preferences and be sure to allow your choices within the host's firewall so client machines can access the services. The gRPC API can be disabled by removing the `-rpc-address` flag if you don't need higher level tools (CLI, graphical tools) for managing machine profiles and you may skip the next section.

### TLS Credentials

*Skip unless gRPC API enabled (optional)*

`bootcfg` provisions machines so, if the gRPC API is enabled, it is important for client-server communications to be TLS encrypted and for clients to authenticate with client certificates. If your organization manages public key infrastructure and a certificate authority, create a server certificate and key for the `bootcfg` service and a client certificate and key for each client tool.

Otherwise, you can generate a self-signed `ca.crt`, a server certificate  (`server.crt`, `server.key`), and client credentials (`client.crt`, `client.key`) with `scripts/tls/cert-gen`. Export the DNS name (or IP) of the provisioner host where `bootcfg` will be accessed.


```sh
$ cd scripts/tls
$ export SAN=DNS.1:bootcfg.example.com,IP.1:192.168.1.100
$ ./cert-gen
```

Place the TLS credentials in the default locations:

```sh
$ sudo mkdir /etc/bootcfg
$ sudo cp ca.crt server.crt server.key /etc/bootcfg/
```

Save `client.crt`, `client.key`, and `ca.crt` to use with a client tool later.

### Start bootcfg

Start the `bootcfg` service and enable it if you'd like it to start on every boot.

```sh
$ sudo systemctl enable bootcfg.service
$ sudo systemctl start bootcfg.service
```

### Verify

Check the bootcfg service can be reached by client machines.


```sh
$ dig bootcfg.example.com
```

Verify you receive a response from the HTTP and API endpoints. All of the following responses are expected:

```sh
$ curl http://bootcfg.example.com:8080
bootcfg
```

If you enabled the gRPC API,

```sh
$ openssl s_client -connect bootcfg.example.com:8081 -CAfile /etc/bootcfg/ca.crt -cert client.key -key client.key
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

### Download CoreOS (optional)

`bootcfg` can serve CoreOS image assets in development or lab environments to reduce bandwidth usage and increase the speed of CoreOS PXE boots and installs to disk.

```sh
$ cd scripts
$ ./get-coreos alpha 1109.1.0 .     # note the "." 3rd argument
```

Move the images to `/var/lib/bootcfg/assets`,

```sh
$ sudo cp -r coreos /var/lib/bootcfg/assets
```

```
/var/lib/bootcfg/assets/
├── coreos
│   └── 1109.1.0
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
$ curl http://bootcfg.example.com:8080/assets/coreos/1109.1.0/
<pre>...
```

For large production environments, use a cache proxy or mirror suitable for your environment to serve CoreOS images.

## Next Steps

### Network

Review [network setup](https://github.com/coreos/coreos-baremetal/blob/master/Documentation/network-setup.md) with your network administrator to set up DHCP, TFTP, and DNS services on your network, if it has not already been done. At a high level, your goals are to:

* Chainload PXE firmwares to iPXE
* Point iPXE client machines to the `bootcfg` iPXE HTTP endpoint `http://bootcfg.example.com:8080/boot.ipxe`
* Ensure `bootcfg.example.com` resolves to your `bootcfg` deployment

CoreOS provides [dnsmasq](https://github.com/coreos/coreos-baremetal/tree/master/contrib/dnsmasq) as `quay.io/coreos/dnsmasq`, if you wish to use rkt or Docker.

### Documentation

View the [documentation](https://github.com/coreos/coreos-baremetal#coreos-on-baremetal) for `bootcfg` service docs, tutorials, example clusters and Ignition configs, PXE booting guides, or machine lifecycle guides.

## rkt

Run the most recent tagged and signed `bootcfg` [release](https://github.com/coreos/coreos-baremetal/releases) ACI. Trust the [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/) for image signature verification.

    sudo rkt trust --prefix coreos.com/bootcfg
    # gpg key fingerprint is: 18AD 5014 C99E F7E3 BA5F  6CE9 50BD D3E0 FC8A 365E
    sudo rkt run --net=host --mount volume=assets,target=/var/lib/bootcfg/assets --volume assets,kind=host,source=$PWD/examples/assets quay.io/coreos/bootcfg:v0.4.0 -- -address=0.0.0.0:8080 -log-level=debug

Create machine profiles, groups, or Ignition configs at runtime with `bootcmd` or by using your own `/var/lib/bootcfg` volume mounts.

## Docker

Run the latest or the most recently tagged `bootcfg` [release](https://github.com/coreos/coreos-baremetal/releases) Docker image.

    sudo docker run --net=host --rm -v $PWD/examples/assets:/var/lib/bootcfg/assets:Z quay.io/coreos/bootcfg:v0.4.0 -address=0.0.0.0:8080 -log-level=debug

Create machine profiles, groups, or Ignition configs at runtime with `bootcmd` or by using your own `/var/lib/bootcfg` volume mounts.

## Kubernetes

*Note: Enhancements to the CLI, and `EtcdStore` backend will improve this deployment strategy in the future.*

Create a `bootcfg` Kubernetes `Deployment` and `Service` based on the example manifests provided in [contrib/k8s](../contrib/k8s).

    kubectl apply -f contrib/k8s/bootcfg-deployment.yaml
    kubectl apply -f contrib/k8s/bootcfg-service.yaml

The `bootcfg` HTTP server should be exposed on NodePort `tcp:31488` on each node in the cluster. `BOOTCFG_LOG_LEVEL` is set to debug.

    kubectl get deployments
    kubectl get services
    kubectl get pods
    kubectl logs POD-NAME

The example manifests use Kubernetes `emptyDir` volumes to back the `bootcfg` FileStore (`/var/lib/bootcfg`). This doesn't provide long-term persistent storage so you may wish to mount your machine groups, profiles, and Ignition configs with a [gitRepo](http://kubernetes.io/docs/user-guide/volumes/#gitrepo) and host image assets on a file server.
