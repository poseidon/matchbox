
# Registry

Clusters can source Docker images from a local insecure Docker registry or rkt ACI images from a simple HTTP file server. This can be useful for development bandwidth reduction and for offline demos.

## Usage

Run the Docker registry with rkt on port 21504, with a mounted data directory as a filesystem backend.

    sudo rkt run --port=5000-tcp:21504 --net=metal0:IP=172.15.0.4 --insecure-options=image docker://registry --mount volume=data,target=/var/lib/registry --volume data,kind=host,source=$PWD/contrib/registry/data

On the `metal0` virtual bridge, dnsmasq will ensure `registry.example.com` resolves to the registry pod.

## Mirror

Pull required images, re-tag them, and push them to the local insecure registry via localhost:21504.

    sudo docker pull quay.io/coreos/hyperkube:v1.3.6_coreos.0
    sudo docker tag quay.io/coreos/hyperkube:v1.3.6_coreos.0 localhost:21504/hyperkube:v1.3.6_coreos.0
    sudo docker push localhost:5000/hyperkube:v1.3.6_coreos.0

Quickly mirror Docker images for Kubernetes by running,

    sudo ./scripts/mirror-k8s-images

## Client Machines

Provision QEMU/KVM VMs with `docker_registry` and `insecure_registry` set in `bootcfg` group metadata:

    "docker_registry": "docker-registry.example.com:",
    "insecure_registry": "true"

A systemd dropin will allow insecure registries and switch image names.

    - name: docker.service
      enable: true
      dropins:
        - name: 40-insecure-registry.conf
          contents: |
            [Service]
            Environment=INSECURE_REGISTRY='--insecure-registry=docker-registry.example.com:5000'
            ExecStart=
            ExecStart=/usr/lib/coreos/dockerd --host=fd:// $DOCKER_OPTS $INSECURE_REGISTRY $DOCKER_CGROUPS $DOCKER_OPT_BIP $DOCKER_OPT_MTU $DOCKER_OPT_IPMASQ

## rkt ACIs

rkt ACI images can be served from an HTTP file server, such as the `bootcfg` `assets` folder.

Use [docker2aci](https://github.com/appc/docker2aci) to download and convert images to ACIs, which are just files.

    docker2aci docker://quay.io/coreos/hyperkube:v1.3.6_coreos.0
    docker2aci docker://quay.io/coreos/flannel:0.5.5

Add the output ACI files to the `bootcfg` mounted `assets`.

    mv coreos-hyperkube-v1.3.6_coreos.0.aci coreos-baremetal/examples/assets/coreos-hyperkube:v1.3.6_coreos.0.aci
    mv coreos-flannel-0.5.5.aci coreos-baremetal/examples/assets/coreos-flannel:v0.5.5.aci

The rename to use colons is used to work with scripts which assume the Docker-style registry/image:tag scheme. rkt can refer to images by URL (e.g. `http://bootcfg.foo:8080/assets/coreos-flannelv0.5.5.aci`).

You may also GPG sign ACI files and configure machines to trust those signatures for verified images.
