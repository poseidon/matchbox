DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
VERSION=$(shell git describe --tags --match=v* --always --dirty)

IMAGE_REPO=coreos/dnsmasq
QUAY_REPO=quay.io/coreos/dnsmasq

.PHONY: all
all: docker-image

.PHONY: tftp
tftp:
	@$(DIR)/get-tftp-files

.PHONY: docker-image
docker-image: tftp
	@sudo docker build --rm=true -t $(IMAGE_REPO):$(VERSION) .
	@sudo docker tag $(IMAGE_REPO):$(VERSION) $(IMAGE_REPO):latest

.PHONY: docker-push
docker-push:
	@sudo docker tag $(IMAGE_REPO):$(VERSION) $(QUAY_REPO):latest
	@sudo docker tag $(IMAGE_REPO):$(VERSION) $(QUAY_REPO):$(VERSION)
	@sudo docker push $(QUAY_REPO):latest
	@sudo docker push $(QUAY_REPO):$(VERSION)
