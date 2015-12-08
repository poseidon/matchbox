
.PHONY: build

build:
	./build

build-docker:
	./docker-build

build-aci:
	./acifile

run-docker:
	docker run -p 8080:8080 --name=bcs -v $(shell echo $$PWD)/static:/static dghubble/bcs:latest

run-rkt:
	rkt --insecure-options=image run --no-overlay bin/bcs-0.0.1-linux-amd64.aci

run-pixiecore:
	./scripts/pixiecore

run-dhcp:
	./scripts/vethdhcp

