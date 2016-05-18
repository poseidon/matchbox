
# Deployment

## rkt

Run the most recent tagged and signed `bootcfg` [release](https://github.com/coreos/coreos-baremetal/releases) ACI. Trust the [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/) for image signature verification.

    sudo rkt trust --prefix coreos.com/bootcfg
    # gpg key fingerprint is: 18AD 5014 C99E F7E3 BA5F  6CE9 50BD D3E0 FC8A 365E
    sudo rkt run --net=host --mount volume=assets,target=/var/lib/bootcfg/assets --volume assets,kind=host,source=$PWD/examples/assets quay.io/coreos/bootcfg:v0.3.0 -- -address=0.0.0.0:8080 -log-level=debug

Create machine profiles, groups, or Ignition configs at runtime with `bootcmd` or by using your own `/var/lib/bootcfg` volume mounts.

## Docker

Run the latest or the most recently tagged `bootcfg` [release](https://github.com/coreos/coreos-baremetal/releases) Docker image.

    sudo docker run --net=host --rm -v $PWD/examples/assets:/var/lib/bootcfg/assets:Z quay.io/coreos/bootcfg:v0.3.0 -address=0.0.0.0:8080 -log-level=debug

Create machine profiles, groups, or Ignition configs at runtime with `bootcmd` or by using your own `/var/lib/bootcfg` volume mounts.

## Kubernetes

*Note: Enhancements to the gRPC API, CLI, and `EtcdStore` backend will improve this deployment strategy in the future.*

Create a `bootcfg` Kubernetes `Deployment` and `Service` based on the example manifests provided in [contrib/k8s](../contrib/k8s).

    kubectl apply -f contrib/k8s/bootcfg-deployment.yaml
    kubectl apply -f contrib/k8s/bootcfg-service.yaml

The `bootcfg` HTTP server should be exposed on NodePort `tcp:31488` on each node in the cluster. `BOOTCFG_LOG_LEVEL` is set to debug.

    kubectl get deployments
    kubectl get services
    kubectl get pods
    kubectl logs POD-NAME

The example manifests use Kubernetes `emptyDir` volumes to back the `bootcfg` FileStore (`/var/lib/bootcfg`). This doesn't provide long-term persistent storage so you may wish to mount your machine groups, profiles, and Ignition configs with a [gitRepo](http://kubernetes.io/docs/user-guide/volumes/#gitrepo) and host image assets on a file server.

## Binary

### Prebuilt

Download a prebuilt binary from the Github [releases](https://github.com/coreos/coreos-baremetal/releases).

    wget https://github.com/coreos/coreos-baremetal/releases/download/VERSION/bootcfg-VERSION-linux-amd64.tar.gz
    wget https://github.com/coreos/coreos-baremetal/releases/download/VERSION/bootcfg-VERSION-linux-amd64.tar.gz.asc

Verify the signature from the [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/).

    gpg --keyserver pgp.mit.edu --recv-key 18AD5014C99EF7E3BA5F6CE950BDD3E0FC8A365E
    gpg --verify bootcfg-VERSION-linux-amd64.tar.gz.asc bootcfg-VERSION-linux-amd64.tar.gz
    # gpg: Good signature from "CoreOS Application Signing Key <security@coreos.com>"

Install the `bootcfg` static binary to `/usr/local/bin`.

    tar xzvf bootcfg-VERSION-linux-amd64.tar.gz
    sudo cp bootcfg-VERSION/bootcfg /usr/local/bin

### Source

Clone the coreos-baremetal project into your $GOPATH.

    go get github.com/coreos/coreos-baremetal/cmd/bootcfg
    cd $GOPATH/src/github.com/coreos/coreos-baremetal

Build `bootcfg` from source.

    make

Install the `bootcfg` static binary to `/usr/local/bin`.

    $ sudo make install

### User/Group

The `bootcfg` service should be run by a non-root user with access to the `bootcfg` data directory (e.g. `/var/lib/bootcfg`). Create a `bootcfg` user and group.

    sudo useradd -U bootcfg

Run the provided script to setup the `bootcfg` data directory.

    sudo ./scripts/setup-data-dir

Add yourself to the `bootcfg` group if you'd like to data by modifying files rather than through the `bootcmd` client.

    SELF=$(whoami)
    sudo gpasswd --add $SELF bootcfg

### Run

Run the `bootcfg` server.

    $ bootcfg -version
    $ bootcfg -address 0.0.0.0:8080
    main: starting bootcfg HTTP server on 0.0.0.0:8080

See [flags and variables](config.md).

## systemd

First, install the `bootcfg` binary from a pre-built binary or from source. Then add and start bootcfg's example systemd unit.

    sudo cp contrib/systemd/bootcfg.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl start bootcfg.service

Check the status and logs.

    systemctl status bootcfg.service
    journalctl -u bootcfg.service

Enable the `bootcfg` service if you'd like it to start at boot time.

    sudo systemctl enable bootcfg.service

### Uninstall

    sudo systemctl stop bootcfg.service
    sudo make uninstall


