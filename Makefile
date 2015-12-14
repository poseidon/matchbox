
.PHONY: build

build:
	./build

build-docker:
	./docker-build

run-docker:
	./docker-run

run-pixiecore:
	./scripts/pixiecore

run-dhcp:
	./scripts/vethdhcp

