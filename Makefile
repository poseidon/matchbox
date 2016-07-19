
export CGO_ENABLED:=0
LD_FLAGS="-w -X github.com/mikeynap/coreos-baremetal/bootcfg/version.Version=$(shell ./git-version)"
LOCAL_BIN=/usr/local/bin

all: build
build: clean bin/bootcfg bin/bootcmd

bin/bootcfg:
	go build -o bin/bootcfg -ldflags $(LD_FLAGS) -a github.com/mikeynap/coreos-baremetal/cmd/bootcfg

bin/bootcmd:
	go build -o bin/bootcmd -ldflags $(LD_FLAGS) -a github.com/mikeynap/coreos-baremetal/cmd/bootcmd

test:
	./test

install:
	cp bin/bootcfg $(LOCAL_BIN)
	cp bin/bootcmd $(LOCAL_BIN)

release: clean _output/coreos-baremetal-linux-amd64.tar.gz _output/coreos-baremetal-darwin-amd64.tar.gz

bin/%/bootcfg:
	GOOS=$* go build -o bin/$*/bootcfg -ldflags $(LD_FLAGS) -a github.com/mikeynap/coreos-baremetal/cmd/bootcfg

bin/%/bootcmd:
	GOOS=$* go build -o bin/$*/bootcmd -ldflags $(LD_FLAGS) -a github.com/mikeynap/coreos-baremetal/cmd/bootcmd

_output/coreos-baremetal-%-amd64.tar.gz: NAME=coreos-baremetal-$(VERSION)-$*-amd64
_output/coreos-baremetal-%-amd64.tar.gz: DEST=_output/$(NAME)
_output/coreos-baremetal-%-amd64.tar.gz: bin/%/bootcfg bin/%/bootcmd
	mkdir -p $(DEST)
	cp bin/$*/bootcfg $(DEST)
	cp bin/$*/bootcmd $(DEST)
	./scripts/release-files $(DEST)
	tar zcvf $(DEST).tar.gz -C _output $(NAME)

clean:
	rm -rf bin
	rm -rf _output

.PHONY: all build test install release clean
.SECONDARY: _output/coreos-baremetal-linux-amd64 _output/coreos-baremetal-darwin-amd64

