export CGO_ENABLED:=0

VERSION=$(shell ./scripts/dev/git-version)
LD_FLAGS="-w -X github.com/coreos/matchbox/matchbox/version.Version=$(VERSION)"

REPO=github.com/coreos/matchbox
IMAGE_REPO=coreos/matchbox
QUAY_REPO=quay.io/coreos/matchbox

all: build

build: clean bin/matchbox

bin/%:
	@go build -o bin/$* -v -ldflags $(LD_FLAGS) $(REPO)/cmd/$*

test:
	@./scripts/dev/test

.PHONY: aci
aci: clean build
	@sudo ./scripts/dev/build-aci

.PHONY: docker-image
docker-image:
	@sudo docker build --rm=true -t $(IMAGE_REPO):$(VERSION) .
	@sudo docker tag $(IMAGE_REPO):$(VERSION) $(IMAGE_REPO):latest

.PHONY: docker-push
docker-push: docker-image
	@sudo docker tag $(IMAGE_REPO):$(VERSION) $(QUAY_REPO):latest
	@sudo docker tag $(IMAGE_REPO):$(VERSION) $(QUAY_REPO):$(VERSION)
	@sudo docker push $(QUAY_REPO):latest
	@sudo docker push $(QUAY_REPO):$(VERSION)

.PHONY: vendor
vendor:
	@glide update --strip-vendor
	@glide-vc --use-lock-file --no-tests --only-code

.PHONY: codegen
codegen: tools
	@./scripts/dev/codegen

.PHONY: tools
tools: bin/protoc bin/protoc-gen-go

bin/protoc:
	@./scripts/dev/get-protoc

bin/protoc-gen-go:
	@go build -o bin/protoc-gen-go $(REPO)/vendor/github.com/golang/protobuf/protoc-gen-go

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

