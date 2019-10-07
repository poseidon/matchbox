# Development

To develop `matchbox` locally, compile the binary and build the container image.

## Static binary

Build the static binary.

```sh
$ make build
```

Test with vendored dependencies.

```sh
$ make test
```

## Container image

Build a container image `coreos/matchbox:latest`.

```sh
$ make docker-image
```

## Version

```sh
$ ./bin/matchbox -version
$ sudo docker run coreos/matchbox:latest -version
```
## Run

Run the binary.

```sh
$ ./bin/matchbox -address=0.0.0.0:8080 -log-level=debug -data-path examples -assets-path examples/assets
```

Run the Docker image on `docker0`.

```sh
$ sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/matchbox:Z -v $PWD/examples/groups/etcd:/var/lib/matchbox/groups:Z coreos/matchbox:latest -address=0.0.0.0:8080 -log-level=debug
```

## bootcmd

Run `bootcmd` against the gRPC API of the service.

```sh
$ ./bin/bootcmd profile list --endpoints 172.18.0.2:8081 --cacert examples/etc/matchbox/ca.crt
```

## Vendor

Add or update dependencies in `go.mod` and vendor.

```
make update
make vendor
```

## Codegen

Generate code from *proto* definitions using `protoc` and the `protoc-gen-go` plugin.

```sh
$ make codegen
```
