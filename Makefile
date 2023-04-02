export CGO_ENABLED:=0

DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
VERSION=$(shell git describe --tags --match=v* --always --dirty)
LD_FLAGS="-w -X github.com/poseidon/matchbox/matchbox/version.Version=$(VERSION)"

REPO=github.com/poseidon/matchbox
LOCAL_REPO=poseidon/matchbox
IMAGE_REPO=quay.io/poseidon/matchbox

.PHONY: all
all: build test vet fmt

.PHONY: build
build:
	@go build -o bin/matchbox -ldflags $(LD_FLAGS) $(REPO)/cmd/matchbox

.PHONY: test
test:
	@go test ./... -cover

.PHONY: vet
vet:
	@go vet -all ./...

.PHONY: fmt
fmt:
	@test -z $$(go fmt ./...)

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: image
image: \
	image-amd64 \
	image-arm64

image-%:
	buildah bud -f Dockerfile \
	-t $(LOCAL_REPO):$(VERSION)-$* \
	--arch $* --override-arch $* \
	--format=docker .

protoc/%:
	podman run --security-opt label=disable \
		-u root \
		--mount type=bind,src=$(DIR),target=/mnt/code \
		quay.io/dghubble/protoc:v3.10.1 \
		--go_out=plugins=grpc,paths=source_relative:. $*

codegen: \
	protoc/matchbox/storage/storagepb/*.proto \
	protoc/matchbox/server/serverpb/*.proto \
	protoc/matchbox/rpc/rpcpb/*.proto

clean:
	@rm -rf bin

clean-release:
	@rm -rf _output

release: \
	clean \
	clean-release \
	_output/matchbox-linux-amd64.tar.gz \
	_output/matchbox-linux-arm.tar.gz \
	_output/matchbox-linux-arm64.tar.gz \
	_output/matchbox-darwin-amd64.tar.gz

bin/linux-amd64/matchbox: GOARGS = GOOS=linux GOARCH=amd64
bin/linux-arm/matchbox: GOARGS = GOOS=linux GOARCH=arm GOARM=6
bin/linux-arm64/matchbox: GOARGS = GOOS=linux GOARCH=arm64
bin/darwin-amd64/matchbox: GOARGS = GOOS=darwin GOARCH=amd64
bin/linux-ppc64le/matchbox: GOARGS = GOOS=linux GOARCH=ppc64le

bin/%/matchbox:
	$(GOARGS) go build -o $@ -ldflags $(LD_FLAGS) -a $(REPO)/cmd/matchbox

_output/matchbox-%.tar.gz: NAME=matchbox-$(VERSION)-$*
_output/matchbox-%.tar.gz: DEST=_output/$(NAME)
_output/matchbox-%.tar.gz: bin/%/matchbox
	mkdir -p $(DEST)
	cp bin/$*/matchbox $(DEST)
	./scripts/dev/release-files $(DEST)
	tar zcvf $(DEST).tar.gz -C _output $(NAME)

.PHONY: all build clean test release
.SECONDARY: _output/matchbox-linux-amd64 _output/matchbox-darwin-amd64

release-sign:
	gpg2 --armor --detach-sign _output/matchbox-$(VERSION)-linux-amd64.tar.gz
	gpg2 --armor --detach-sign _output/matchbox-$(VERSION)-linux-arm.tar.gz
	gpg2 --armor --detach-sign _output/matchbox-$(VERSION)-linux-arm64.tar.gz
	gpg2 --armor --detach-sign _output/matchbox-$(VERSION)-darwin-amd64.tar.gz

release-verify: NAME=_output/matchbox
release-verify:
	gpg2 --verify $(NAME)-$(VERSION)-linux-amd64.tar.gz.asc $(NAME)-$(VERSION)-linux-amd64.tar.gz
	gpg2 --verify $(NAME)-$(VERSION)-linux-arm.tar.gz.asc $(NAME)-$(VERSION)-linux-arm.tar.gz
	gpg2 --verify $(NAME)-$(VERSION)-linux-arm64.tar.gz.asc $(NAME)-$(VERSION)-linux-arm64.tar.gz
	gpg2 --verify $(NAME)-$(VERSION)-darwin-amd64.tar.gz.asc $(NAME)-$(VERSION)-darwin-amd64.tar.gz
