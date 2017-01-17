
export CGO_ENABLED:=0
LD_FLAGS="-w -X github.com/coreos/matchbox/matchbox/version.Version=$(shell ./git-version)"
LOCAL_BIN=/usr/local/bin

all: build
build: clean bin/matchbox bin/bootcmd

tools:
	./scripts/gentools

codegen: tools
	./scripts/codegen

bin/matchbox:
	go build -o bin/matchbox -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/matchbox

bin/bootcmd:
	go build -o bin/bootcmd -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/bootcmd

test:
	./test

install:
	cp bin/matchbox $(LOCAL_BIN)
	cp bin/bootcmd $(LOCAL_BIN)

release: \
	clean \
	_output/matchbox-linux-amd64.tar.gz \
	_output/matchbox-linux-arm.tar.gz \
	_output/matchbox-linux-arm64.tar.gz \
	_output/matchbox-darwin-amd64.tar.gz \

# matchbox

bin/linux-amd64/matchbox:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/matchbox -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/matchbox

bin/linux-arm/matchbox:
	GOOS=linux GOARCH=arm go build -o bin/linux-arm/matchbox -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/matchbox

bin/linux-arm64/matchbox:
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/matchbox -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/matchbox

bin/darwin-amd64/matchbox:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/matchbox -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/matchbox

# bootcmd

bin/linux-amd64/bootcmd:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/bootcmd -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/bootcmd

bin/linux-arm/bootcmd:
	GOOS=linux GOARCH=arm go build -o bin/linux-arm/bootcmd -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/bootcmd

bin/linux-arm64/bootcmd:
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/bootcmd -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/bootcmd

bin/darwin-amd64/bootcmd:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/bootcmd -ldflags $(LD_FLAGS) -a github.com/coreos/matchbox/cmd/bootcmd

_output/matchbox-%.tar.gz: NAME=matchbox-$(VERSION)-$*
_output/matchbox-%.tar.gz: DEST=_output/$(NAME)
_output/matchbox-%.tar.gz: bin/%/matchbox bin/%/bootcmd
	mkdir -p $(DEST)
	cp bin/$*/matchbox $(DEST)
	cp bin/$*/bootcmd $(DEST)
	./scripts/release-files $(DEST)
	tar zcvf $(DEST).tar.gz -C _output $(NAME)

clean:
	rm -rf tools
	rm -rf bin
	rm -rf _output

.PHONY: all build tools test install release clean
.SECONDARY: _output/matchbox-linux-amd64 _output/matchbox-darwin-amd64

