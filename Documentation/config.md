
# Flags and Variables

Configuration arguments can be provided as flags or as environment variables.

| flag | variable | default | example |
|------|----------|---------|---------|
| -address | MATCHBOX_ADDRESS | 127.0.0.1:8080 | 0.0.0.0:8080 |
| -log-level | MATCHBOX_LOG_LEVEL | info | critical, error, warning, notice, info, debug |
| -data-path | MATCHBOX_DATA_PATH | /var/lib/matchbox | ./examples |
| -assets-path | MATCHBOX_ASSETS_PATH | /var/lib/matchbox/assets | ./examples/assets |
| -rpc-address | MATCHBOX_RPC_ADDRESS | (gRPC API disabled) | 0.0.0.0:8081 |
| -cert-file | MATCHBOX_CERT_FILE | /etc/matchbox/server.crt | ./examples/etc/matchbox/server.crt |
| -key-file | MATCHBOX_KEY_FILE | /etc/matchbox/server.key | ./examples/etc/matchbox/server.key
| -ca-file | MATCHBOX_CA_FILE | /etc/matchbox/ca.crt | ./examples/etc/matchbox/ca.crt |
| -key-ring-path | MATCHBOX_KEY_RING_PATH | (no key ring) | ~/.secrets/vault/matchbox/secring.gpg |
| (no flag) | MATCHBOX_PASSPHRASE | (no passphrase) | "secret passphrase" |

## Files and Directories

| Data | Default Location                                  |
|:---------|:--------------------------------------------------|
| data     | /var/lib/matchbox/{profiles,groups,ignition,cloud,generic} |
| assets   | /var/lib/matchbox/assets                           |

| gRPC API TLS Credentials | Default Location                  |
|:---------|:--------------------------------------------------|
| CA certificate | /etc/matchbox/ca.crt                         |
| Server certificate | /etc/matchbox/server.crt                 |
| Server private key | /etc/matchbox/server.key                 |
| Client certificate | /etc/matchbox/client.crt                 |
| Client private key | /etc/matchbox/client.key                 |

## Version

```sh
$ ./bin/matchbox -version
$ sudo rkt run quay.io/coreos/matchbox:latest -- -version
$ sudo docker run quay.io/coreos/matchbox:latest -version
```

## Usage

Run the binary.

```sh
$ ./bin/matchbox -address=0.0.0.0:8080 -log-level=debug -data-path=examples -assets-path=examples/assets
```

Run the latest ACI with rkt.

```sh
$ sudo rkt run --mount volume=assets,target=/var/lib/matchbox/assets --volume assets,kind=host,source=$PWD/examples/assets quay.io/coreos/matchbox:latest -- -address=0.0.0.0:8080 -log-level=debug
```

Run the latest Docker image.

```sh
$ sudo docker run -p 8080:8080 --rm -v $PWD/examples/assets:/var/lib/matchbox/assets:Z quay.io/coreos/matchbox:latest -address=0.0.0.0:8080 -log-level=debug
```

#### With Examples

Mount `examples` to pre-load the [example](../examples/README.md) machine groups and profiles. Run the container with rkt,

```sh
$ sudo rkt run --net=metal0:IP=172.18.0.2 --mount volume=data,target=/var/lib/matchbox --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/matchbox/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/matchbox:latest -- -address=0.0.0.0:8080 -log-level=debug
```

or with Docker.

```sh
$ sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/matchbox:Z -v $PWD/examples/groups/etcd:/var/lib/matchbox/groups:Z quay.io/coreos/matchbox:latest -address=0.0.0.0:8080 -log-level=debug
```

### gRPC API

The gRPC API allows clients with a TLS client certificate and key to make RPC requests to programmatically create or update `matchbox` resources. The API can be enabled with the `-rpc-address` flag and by providing a TLS server certificate and key with `-cert-file` and `-key-file` and a CA certificate for authenticating clients with `-ca-file`.

Run the binary with TLS credentials from `examples/etc/matchbox`.

```sh
$ ./bin/matchbox -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug -data-path=examples -assets-path=examples/assets -cert-file examples/etc/matchbox/server.crt -key-file examples/etc/matchbox/server.key -ca-file examples/etc/matchbox/ca.crt
```

Clients, such as `bootcmd`, verify the server's certificate with a CA bundle passed via `-ca-file` and present a client certificate and key via `-cert-file` and `-key-file` to cal the gRPC API.

```sh
$ ./bin/bootcmd profile list --endpoints 127.0.0.1:8081 --ca-file examples/etc/matchbox/ca.crt --cert-file examples/etc/matchbox/client.crt --key-file examples/etc/matchbox/client.key
```

#### With rkt

Run the ACI with rkt and TLS credentials from `examples/etc/matchbox`.

```sh
$ sudo rkt run --net=metal0:IP=172.18.0.2 --mount volume=data,target=/var/lib/matchbox --volume data,kind=host,source=$PWD/examples,readOnly=true --mount volume=config,target=/etc/matchbox --volume config,kind=host,source=$PWD/examples/etc/matchbox --mount volume=groups,target=/var/lib/matchbox/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/matchbox:latest -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
```

A `bootcmd` client can call the gRPC API running at the IP used in the rkt example.

```sh
$ ./bin/bootcmd profile list --endpoints 172.18.0.2:8081 --ca-file examples/etc/matchbox/ca.crt --cert-file examples/etc/matchbox/client.crt --key-file examples/etc/matchbox/client.key
```

#### With docker

Run the Docker image with TLS credentials from `examples/etc/matchbox`.

```sh
$ sudo docker run -p 8080:8080 -p 8081:8081 --rm -v $PWD/examples:/var/lib/matchbox:Z -v $PWD/examples/etc/matchbox:/etc/matchbox:Z,ro -v $PWD/examples/groups/etcd:/var/lib/matchbox/groups:Z quay.io/coreos/matchbox:latest -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
```

A `bootcmd` client can call the gRPC API running at the IP used in the Docker example.

```sh
$ ./bin/bootcmd profile list --endpoints 127.0.0.1:8081 --ca-file examples/etc/matchbox/ca.crt --cert-file examples/etc/matchbox/client.crt --key-file examples/etc/matchbox/client.key
```

### OpenPGP [Signing](openpgp.md)

Run with the binary with a test key.

```sh
$ export MATCHBOX_PASSPHRASE=test
$ ./bin/matchbox -address=0.0.0.0:8080 -key-ring-path matchbox/sign/fixtures/secring.gpg -data-path=examples -assets-path=examples/assets
```

Run the ACI with a test key.

```sh
$ sudo rkt run --net=metal0:IP=172.18.0.2 --set-env=MATCHBOX_PASSPHRASE=test --mount volume=secrets,target=/secrets --volume secrets,kind=host,source=$PWD/matchbox/sign/fixtures --mount volume=data,target=/var/lib/matchbox --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/matchbox/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd quay.io/coreos/matchbox:latest -- -address=0.0.0.0:8080 -key-ring-path secrets/secring.gpg
```

Run the Docker image with a test key.

```sh
$ sudo docker run -p 8080:8080 --rm --env MATCHBOX_PASSPHRASE=test -v $PWD/examples:/var/lib/matchbox:Z -v $PWD/examples/groups/etcd:/var/lib/matchbox/groups:Z -v $PWD/matchbox/sign/fixtures:/secrets:Z quay.io/coreos/matchbox:latest -address=0.0.0.0:8080 -log-level=debug -key-ring-path secrets/secring.gpg
```
