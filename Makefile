
.PHONY: build

build:
	./build

docker-build:
	docker build --rm=true -t dghubble.io/metapxe .

aci-build:
	./acifile

run-docker:
	docker run -p 8081:8081 -v $(shell echo $$PWD)/static:/static dghubble.io/metapxe:latest

run-rkt:
	rkt --insecure-options=image run --no-overlay bin/metapxe-0.0.1-linux-amd64.aci

run-pixiecore:
	docker run -v $(shell echo $$PWD)/static:/static danderson/pixiecore -api http://172.17.0.2:8081/

run-dhcp:
	./scripts/vethdhcp




