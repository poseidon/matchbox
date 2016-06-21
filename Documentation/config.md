
# Flags and Variables

Configuration arguments can be provided as flags or as environment variables.

| flag | variable | default | example |
|------|----------|---------|---------|
| -address | BOOTCFG_ADDRESS | 127.0.0.1:8080 | 0.0.0.0:8080 |
| -log-level | BOOTCFG_LOG_LEVEL | info | critical, error, warning, notice, info, debug |
| -data-path | BOOTCFG_DATA_PATH | /var/lib/bootcfg | ./examples |
| -assets-path | BOOTCFG_ASSETS_PATH | /var/lib/bootcfg/assets | ./examples/assets |
| -rpc-address | BOOTCFG_RPC_ADDRESS | (gRPC API disabled) | 127.0.0.1:8081 |
| -cert-file | BOOTCFG_CERT_FILE | /etc/bootcfg/server.crt | ./examples/etc/bootcfg/server.crt |
| -key-file | BOOTCFG_KEY_FILE | /etc/bootcfg/server.key | ./examples/etc/bootcfg/server.key
| -ca-file | BOOTCFG_CA_FILE | /etc/bootcfg/ca.crt | ./examples/etc/bootcfg/ca.crt |
| -key-ring-path | BOOTCFG_KEY_RING_PATH | (no key ring) | ~/.secrets/vault/bootcfg/secring.gpg |
| (no flag) | BOOTCFG_PASSPHRASE | (no passphrase) | "secret passphrase" |

## Files and Directories

| Data | Default Location                                  |
|:---------|:--------------------------------------------------|
| data     | /var/lib/bootcfg/{profiles,groups,ignition,cloud} |
| assets   | /var/lib/bootcfg/assets                           |
| environment variables | /etc/bootcfg.env                     |

| gRPC API TLS Credentials | Default Location                  |
|:---------|:--------------------------------------------------|
| CA certificate | /etc/bootcfg/ca.crt                         |
| Server certificate | /etc/bootcfg/server.crt                 |
| Server private key | /etc/bootcfg/server.key                 |
| Client certificate | /etc/bootcfg/client.crt                 |
| Client private key | /etc/bootcfg/client.key                 |

## Version

    ./bin/bootcfg -version
    sudo rkt run quay.io/coreos/bootcfg:latest -- -version
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

Run the ACI with rkt. Mounts are used to add the provided examples.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -log-level=debug

Run the Docker image. Mounts are used to add the provided examples.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd:/var/lib/bootcfg/groups:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

### gRPC API

The gRPC API can be enabled with the `-rpc-address` flag and by providing a TLS server certificate and key with `-cert-file` and `-key-file` and a CA certificate for authenticating clients with `-ca-file`. gRPC clients (such as `bootcmd`) must verify the server's certificate with a CA bundle passed via `-ca-file` and present a client certificate and key via `-cert-file` and `-key-file`.

Run the binary with TLS credentials from `examples/etc/bootcfg`.

    ./bin/bootcfg -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug -data-path=examples -assets-path=examples/assets -cert-file examples/etc/bootcfg/server.crt -key-file examples/etc/bootcfg/server.key -ca-file examples/etc/bootcfg/ca.crt

A `bootcmd` client can call the gRPC API.

    ./bin/bootcmd profile list --endpoints 127.0.0.1:8081 --ca-file examples/etc/bootcfg/ca.crt --cert-file examples/etc/bootcfg/client.crt --key-file examples/etc/bootcfg/client.key

Run the ACI with rkt and TLS credentials from `examples/etc/bootcfg`.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples,readOnly=true --mount volume=config,target=/etc/bootcfg --volume config,kind=host,source=$PWD/examples/etc/bootcfg --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug

A `bootcmd` client can call the gRPC API running at the IP used in the rkt example.

    ./bin/bootcmd profile list --endpoints 172.15.0.2:8081 --ca-file examples/etc/bootcfg/ca.crt --cert-file examples/etc/bootcfg/client.crt --key-file examples/etc/bootcfg/client.key

Run the Docker image with TLS credentials from `examples/etc/bootcfg`.

    sudo docker run -p 8080:8080 -p 8081:8081 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/etc/bootcfg:/etc/bootcfg:Z,ro -v $PWD/examples/groups/etcd:/var/lib/bootcfg/groups:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug

A `bootcmd` client can call the gRPC API running at the IP used in the Docker example.

    ./bin/bootcmd profile list --endpoints 127.0.0.1:8081 --ca-file examples/etc/bootcfg/ca.crt --cert-file examples/etc/bootcfg/client.crt --key-file examples/etc/bootcfg/client.key

#### With [OpenPGP Signing](openpgp.md)

Run with the binary with a test key.

    export BOOTCFG_PASSPHRASE=test
    ./bin/bootcfg -address=0.0.0.0:8080 -key-ring-path bootcfg/sign/fixtures/secring.gpg -data-path=examples -assets-path=examples/assets

Run the ACI with a test key.

    sudo rkt run --net=metal0:IP=172.15.0.2 --set-env=BOOTCFG_PASSPHRASE=test --mount volume=secrets,target=/secrets --volume secrets,kind=host,source=$PWD/bootcfg/sign/fixtures --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -key-ring-path secrets/secring.gpg

Run the Docker image with a test key.

    sudo docker run -p 8080:8080 --rm --env BOOTCFG_PASSPHRASE=test -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd:/var/lib/bootcfg/groups:Z -v $PWD/bootcfg/sign/fixtures:/secrets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -key-ring-path secrets/secring.gpg
