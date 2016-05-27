
BIN_DIR=/usr/local/bin

all: build

build:
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

.PHONY: build test install
