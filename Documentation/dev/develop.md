
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

## Run

Run the ACI with rkt on `metal0`.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/assets --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/data --volume data,kind=host,source=$PWD/examples bootcfg.aci -- -address=0.0.0.0:8080 -log-level=debug -config /data/etcd-rkt.yaml

Alternately, run the Docker image on `docker0`.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/data:Z -v $PWD/assets:/assets:Z coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -config /data/etcd-docker.yaml

## Version

    ./bin/bootcfg -version
    sudo rkt --insecure-options=image run bootcfg.aci -- -version
    sudo docker run coreos/bootcfg:latest -version
