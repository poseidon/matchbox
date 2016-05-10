
# Flags and Variables

Configuration arguments can be provided as flags or as environment variables.

| flag | variable | example |
|------|----------|---------|
| -address | BOOTCFG_ADDRESS | 0.0.0.0:8080 |
| -rpc-address | BOOTCFG_RPC_ADDRESS | 127.0.0.1:8081
| -data-path | BOOTCFG_DATA_PATH | /var/lib/bootcfg |
| -assets-path | BOOTCFG_ASSETS_PATH | /var/lib/bootcfg/assets |
| -key-ring-path | BOOTCFG_KEY_RING_PATH | ~/.secrets/vault/bootcfg/secring.gpg |
| Disallowed | BOOTCFG_PASSPHRASE | secret passphrase |
| -log-level | BOOTCFG_LOG_LEVEL | critical, error, warning, notice, info, debug |

## Files and Directories

| Contents  | Default Location  |
|-----------|-------------------|
| data      | /var/lib/bootcfg/{profiles,groups,ignition,cloud} |
| assets    | /var/lib/bootcfg/assets |

## Version

    ./bin/bootcfg -version
    sudo rkt --insecure-options=image run quay.io/coreos/bootcfg:latest -- -version
    sudo docker run quay.io/coreos/bootcfg:latest -version

## Minimal

Start the latest ACI with rkt.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/var/lib/bootcfg/assets --volume assets,kind=host,source=$PWD/examples/assets quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -log-level=debug

Start the latest Docker image.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples/assets:/var/lib/bootcfg/assets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

To start containers with the example machine Groups and Profiles, see the commands below.

## Examples

Run the binary.

    ./bin/bootcfg -address=0.0.0.0:8080 -log-level=debug -data-path=examples -assets-path=examples/assets

Run the latest ACI with rkt. Mounts are used to add the provided examples.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -log-level=debug

Run the latest Docker image. Mounts are used to add the provided examples.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd:/var/lib/bootcfg/groups:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

#### With [OpenPGP Signing](openpgp.md)

Run with the binary with a test key.

    export BOOTCFG_PASSPHRASE=test
    ./bin/bootcfg -address=0.0.0.0:8080 -key-ring-path bootcfg/sign/fixtures/secring.gpg -data-path=examples -assets-path=examples/assets

Run the ACI with a test key.

    sudo rkt run --net=metal0:IP=172.15.0.2 --set-env=BOOTCFG_PASSPHRASE=test --mount volume=secrets,target=/secrets --volume secrets,kind=host,source=$PWD/bootcfg/sign/fixtures --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -key-ring-path secrets/secring.gpg

Run the Docker image with a test key.

    sudo docker run -p 8080:8080 --rm --env BOOTCFG_PASSPHRASE=test -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd:/var/lib/bootcfg/groups:Z -v $PWD/bootcfg/sign/fixtures:/secrets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -key-ring-path secrets/secring.gpg
