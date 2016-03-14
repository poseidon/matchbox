
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

    ./bin/bootcfg -address=0.0.0.0:8080 -log-level=debug -data-path examples/ -config examples/etcd-rkt.yaml

Run the ACI with rkt on `metal0`.

    sudo rkt --insecure-options=image run --net=metal0:IP=172.15.0.2 --mount volume=assets,target=/var/bootcfg --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/etc/bootcfg --volume data,kind=host,source=$PWD/examples bootcfg.aci -- -address=0.0.0.0:8080 -log-level=debug -config /etc/bootcfg/etcd-rkt.yaml

Alternately, run the Docker image on `docker0`.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/etc/bootcfg:Z -v $PWD/assets:/var/bootcfg:Z coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug -config /etc/bootcfg/etcd-docker.yaml
