
.PHONY: build

build:
	./build

docker-build:
	docker build --rm=true -t dghubble.io/metapxe .

aci-build:
	./acifile

docker-run:
	docker run -p 8080:8080 -v $(shell echo $$PWD)/static:/static dghubble.io/metapxe:latest

rkt-run:
	# Fedora 23 issue https://github.com/coreos/rkt/issues/1727
	rkt --insecure-options=image run --no-overlay bin/metapxe-0.0.1-linux-amd64.aci
