
BIN_DIR=/usr/local/bin

all: build

bin/bootcfg:
	./build

bin/bootcmd:
	./build

test:
	./test

install:
	cp bin/bootcfg $(BIN_DIR)
	cp bin/bootcmd $(BIN_DIR)
	@echo "**************"
	@echo "INSTALL SUCCESS"
	@echo "**************"
	@echo "bootcfg was installed to /usr/local/bin/bootcfg"
	@echo "bootcmd was installed to /usr/local/bin/bootcmd"

uninstall:
	rm $(BIN_DIR)/bootcfg
	rm $(BIN_DIR)/bootcmd

release: clean _output/coreos-baremetal-linux-amd64.tar.gz

_output/coreos-baremetal-%-amd64:
	mkdir -p $@

_output/coreos-baremetal-%-amd64.tar.gz: bin/bootcfg bin/bootcmd | _output/coreos-baremetal-%-amd64
	./scripts/release-files $|
	tar zcvf $@ -C _output coreos-baremetal-$*-amd64

clean:
	rm bin/bootcfg
	rm bin/bootcmd
	rm -rf _output

.PHONY: build clean install test

.SECONDARY: _output/coreos-baremetal-linux-amd64
