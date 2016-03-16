
# Flags and Variables

Configuration arguments can be provided as flags or as environment variables.

| flag | variable | example |
|------|----------|---------|
| -address | BOOTCFG_ADDRESS | 127.0.0.1:8080 |
| -config | BOOTCFG_CONFIG | /etc/bootcfg.conf |
| -data-path | BOOTCFG_DATA_PATH | /etc/bootcfg |
| -assets-path | BOOTCFG_ASSETS_PATH | /var/bootcfg |
| -key-ring-path | BOOTCFG_KEY_RING_PATH | ~/.secrets/vault/bootcfg/secring.gpg |
| Disallowed | BOOTCFG_PASSPHRASE | secret passphrase |
| -log-level | BOOTCFG_LOG_LEVEL | critical, error, warning, notice, info, debug |

## Files and Directories

| Contents  | Default Location  |
|-----------|-------------------|
| conf file | /etc/bootcfg.conf |
| configs   | /etc/bootcfg/{profiles,ignition,cloud} |
| assets    | /var/bootcfg/     |

## Check Version

    ./bin/bootcfg -version
    sudo rkt --insecure-options=image run quay.io/coreos/bootcfg:latest -- -version
    sudo docker run quay.io/coreos/bootcfg:latest -version

## Examples

Run the binary.

    ./bin/bootcfg -address=0.0.0.0:8080 -log-level=debug -data-path examples/ -config examples/etcd-rkt.yaml

Run the latest ACI with rkt.

    sudo rkt --insecure-options=image run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/var/bootcfg --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/etc/bootcfg --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -log-level=debug -config /etc/bootcfg/etcd-rkt.yaml

Run the latest Docker image.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/etc/bootcfg:Z -v $PWD/assets:/var/bootcfg:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -config /etc/bootcfg/etcd-docker.yaml

#### With [OpenPGP Signing](openpgp.md)

Run with the binary with a test key.

    export BOOTCFG_PASSPHRASE=test
    ./bin/bootcfg -address=0.0.0.0:8080 -key-ring-path bootcfg/sign/fixtures/secring.gpg -data-path examples/ -config examples/etcd-rkt.yaml

Run the ACI with a test key.

    sudo rkt --insecure-options=image run --net=metal0:IP=172.15.0.2 --set-env=BOOTCFG_PASSPHRASE=test --mount volume=secrets,target=/secrets --volume secrets,kind=host,source=$PWD/bootcfg/sign/fixtures --mount volume=assets,target=/var/bootcfg --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/etc/bootcfg --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -config /etc/bootcfg/etcd-rkt.yaml -key-ring-path secrets/secring.gpg

Run the Docker image with a test key.

    sudo docker run -p 8080:8080 --rm --env BOOTCFG_PASSPHRASE=test -v $PWD/examples:/etc/bootcfg:Z -v $PWD/assets:/var/bootcfg:Z -v $PWD/bootcfg/sign/fixtures:/secrets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -config /etc/bootcfg/etcd-docker.yaml -key-ring-path secrets/secring.gpg
