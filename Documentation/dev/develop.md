
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

    sudo rkt --insecure-options=image run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=config,target=/etc/bootcfg --volume config,kind=host,source=$PWD/examples/etc/bootcfg --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd bootcfg.aci -- -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug

Alternately, run the Docker image on `docker0`.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd-docker:/var/lib/bootcfg/groups:Z coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

### bootcmd

Run `bootcmd` against the gRPC API of the service running via rkt.

    ./bin/bootcmd profile list --endpoints 172.15.0.2:8081 --cacert examples/etc/bootcfg/ca.crt

## Dependencies

Project dependencies are commited to the `vendor` directory, so Go 1.6+ users can clone to their `GOPATH` and build or test immediately. Go 1.5 users should set `GO15VENDOREXPERIMENT=1`.

Project developers should use [glide](https://github.com/Masterminds/glide) to manage commited dependencies under `vendor`. Configure `glide.yaml` as desired. Use `glide update` to download and update dependencies listed in `glide.yaml` into `/vendor` (do **not** use glide `get`).

    glide update --update-vendored --strip-vendor --strip-vcs

Recursive dependencies are also vendored. A `glide.lock` will be created to represent the exact versions of each dependency.

With an empty `vendor` directory, you can install the `glide.lock` dependencies.

    rm -rf vendor/
    glide install --strip-vendor --strip-vcs
