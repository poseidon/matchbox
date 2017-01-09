
# Development

To develop `matchbox` locally, compile the binary and build the container image. The following tasks can also be performed using the `Makefile` if you prefer.

## Static Binary

Build the static binary.

    ./build

Test with vendored dependencies.

    ./test

## Vendor

Use `glide` to manage dependencies committed to the `vendor` directory.

    glide update --strip-vendor --skip-test

## Codegen

Generate code from *proto* definitions using `protoc` and the `protoc-gen-go` plugin.

    # install tools
    ./scripts/gentools
    # generate code from proto definitions in subpackages
    ./scripts/codegen

## Container Image

Build an ACI `matchbox.aci`.

    ./build-aci

Alternately, build a Docker image `coreos/matchbox:latest`.

    sudo ./build-docker

## Version

    ./bin/matchbox -version
    sudo rkt --insecure-options=image run matchbox.aci -- -version
    sudo docker run coreos/matchbox:latest -version

## Run

Run the binary.

    ./bin/matchbox -address=0.0.0.0:8080 -log-level=debug -data-path examples -assets-path examples/assets

Run the container image with rkt, on `metal0`.

    sudo rkt --insecure-options=image run --net=metal0:IP=172.18.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=config,target=/etc/bootcfg --volume config,kind=host,source=$PWD/examples/etc/bootcfg --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd matchbox.aci -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug

Alternately, run the Docker image on `docker0`.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd:/var/lib/bootcfg/groups:Z coreos/matchbox:latest -address=0.0.0.0:8080 -log-level=debug

### bootcmd

Run `bootcmd` against the gRPC API of the service running via rkt.

    ./bin/bootcmd profile list --endpoints 172.18.0.2:8081 --cacert examples/etc/bootcfg/ca.crt
