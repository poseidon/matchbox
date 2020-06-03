DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
VERSION=$(shell git rev-parse HEAD)

LOCAL_REPO=poseidon/dnsmasq
IMAGE_REPO=quay.io/poseidon/dnsmasq

.PHONY: all
all: image

.PHONY: tftp
tftp:
	@$(DIR)/get-tftp-files

.PHONY: image
image: tftp
	@buildah bud -t $(LOCAL_REPO):$(VERSION) .
	@buildah tag $(LOCAL_REPO):$(VERSION) $(LOCAL_REPO):latest

.PHONY: push
push:
	@buildah tag $(LOCAL_REPO):$(VERSION) $(IMAGE_REPO):$(VERSION)
	@buildah tag $(LOCAL_REPO):$(VERSION) $(IMAGE_REPO):latest
	@buildah push docker://$(IMAGE_REPO):$(VERSION)
	@buildah push docker://$(IMAGE_REPO):latest
