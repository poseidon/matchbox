
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

## Check Version

    ./bin/bootcfg -version
    sudo rkt --insecure-options=image run bootcfg.aci -- -version
    sudo docker run coreos/bootcfg:latest -version

## Run

Run the binary.

    ./bin/bootcfg -address=0.0.0.0:8080 -log-level=debug -data-path examples -assets-path assets

Run the ACI with rkt on `metal0`.

    sudo rkt --insecure-options=image run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/var/bootcfg --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/etc/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/etc/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/etcd bootcfg.aci -- -address=0.0.0.0:8080 -log-level=debug

Alternately, run the Docker image on `docker0`.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/etc/bootcfg:Z -v $PWD/assets:/var/bootcfg:Z -v $PWD/examples/groups/etcd:/etc/bootcfg/groups:Z coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug