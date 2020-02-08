DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
VERSION=$(shell git rev-parse HEAD)

LOCAL_REPO=poseidon/dnsmasq
IMAGE_REPO=quay.io/poseidon/dnsmasq

.PHONY: all
all: docker-image

.PHONY: tftp
tftp:
	@$(DIR)/get-tftp-files

.PHONY: image
image:
	@buildah bud -t $(LOCAL_REPO):$(VERSION) .
	@buildah tag $(LOCAL_REPO):$(VERSION) $(LOCAL_REPO):latest

.PHONY: push
push:
	@buildah tag $(LOCAL_REPO):$(VERSION) $(IMAGE_REPO):$(VERSION)
	@buildah tag $(LOCAL_REPO):$(VERSION) $(IMAGE_REPO):latest
	@buildah push docker://$(IMAGE_REPO):$(VERSION)
	@buildah push docker://$(IMAGE_REPO):latest

# for travis-only

.PHONY: docker-image
docker-image: tftp
	@sudo docker build --rm=true -t $(LOCAL_REPO):$(VERSION) .
	@sudo docker tag $(LOCAL_REPO):$(VERSION) $(LOCAL_REPO):latest

.PHONY: docker-push
docker-push:
	@sudo docker tag $(LOCAL_REPO):$(VERSION) $(IMAGE_REPO):latest
	@sudo docker tag $(LOCAL_REPO):$(VERSION) $(IMAGE_REPO):$(VERSION)
	@sudo docker push $(IMAGE_REPO):latest
	@sudo docker push $(IMAGE_REPO):$(VERSION)
