export CGO_ENABLED:=0
export GO111MODULE=on
export GOFLAGS=-mod=vendor

VERSION=$(shell git describe --tags --match=v* --always --dirty)
LD_FLAGS="-w -X github.com/poseidon/matchbox/matchbox/version.Version=$(VERSION)"

REPO=github.com/poseidon/matchbox
LOCAL_REPO=poseidon/matchbox
IMAGE_REPO=quay.io/poseidon/matchbox

.PHONY: all
all: build test vet lint fmt

.PHONY: build
build: clean bin/matchbox

bin/%:
	git describe --tags --match=v* --always --dirty
	git status
	git diff
	@go build -o bin/$* -ldflags $(LD_FLAGS) $(REPO)/cmd/$*

.PHONY: test
test:
	@go test ./... -cover

.PHONY: vet
vet:
	@go vet -all ./...

.PHONY: lint
lint:
	@golint -set_exit_status `go list ./... | grep -v pb`

.PHONY: fmt
fmt:
	@test -z $$(go fmt ./...)

.PHONY: docker-image
docker-image:
	@sudo docker build --rm=true -t $(LOCAL_REPO):$(VERSION) .
	@sudo docker tag $(LOCAL_REPO):$(VERSION) $(LOCAL_REPO):latest

.PHONY: docker-push
docker-push: docker-image
	@sudo docker tag $(LOCAL_REPO):$(VERSION) $(IMAGE_REPO):latest
	@sudo docker tag $(LOCAL_REPO):$(VERSION) $(IMAGE_REPO):$(VERSION)
	@sudo docker push $(IMAGE_REPO):latest
	@sudo docker push $(IMAGE_REPO):$(VERSION)

.PHONY: update
update:
	@GOFLAGS="" go get -u
	@go mod tidy

.PHONY: vendor
vendor:
	@go mod vendor

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

