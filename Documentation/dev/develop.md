
# bootcfg Development

Develop `bootcfg` locally.

## Binary

Build the static binary.

    ./build

Test with vendored dependencies.

    ./test

## Container Image

Build an ACI `bootcfg.aci`.

    ./build-aci

Alternately, build a Docker image `coreos/bootcfg:latest`.

    sudo ./build-docker

## Version

    ./bin/bootcfg -version
    sudo rkt --insecure-options=image run bootcfg.aci -- -version
    sudo docker run coreos/bootcfg:latest -version

## Run

Run the binary.

    ./bin/bootcfg -address=0.0.0.0:8080 -log-level=debug -data-path examples -assets-path examples/assets

Run the ACI with rkt on `metal0`.

    sudo rkt --insecure-options=image run --net=metal0:IP=172.18.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=config,target=/etc/bootcfg --volume config,kind=host,source=$PWD/examples/etc/bootcfg --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd bootcfg.aci -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug

Alternately, run the Docker image on `docker0`.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd-docker:/var/lib/bootcfg/groups:Z coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

### bootcmd

Run `bootcmd` against the gRPC API of the service running via rkt.

    ./bin/bootcmd profile list --endpoints 172.18.0.2:8081 --cacert examples/etc/bootcfg/ca.crt

## Dependencies

Dependencies are commited to the `vendor` directory so builds and tests can be isolated from the Go workspace environment.

[Glide](https://github.com/Masterminds/glide/releases) 0.12 is used to download, update, and manage dependencies listed in `glide.yaml` recursively.

    glide update --strip-vendor --skip-test

You can regenerate or repair your `vendor` directory from scratch by installing the exact versions pinned in `glide.lock`.

    rm -rf vendor/
    glide install --strip-vendor --skip-test
