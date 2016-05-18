
BIN_DIR=/usr/local/bin
DATA_DIR=/var/lib/bootcfg
ENV_FILE=/etc/bootcfg.env

all: build

build:
	./build

test:
	./test

install:
	touch ${ENV_FILE}
	cp bin/bootcfg $(BIN_DIR)
	cp bin/bootcmd $(BIN_DIR)
	@echo "**************"
	@echo "INSTALL SUCESS"
	@echo "**************"
	@echo "bootcfg was installed to /usr/local/bin/bootcfg"
	@echo "bootcmd was installed to /usr/local/bin/bootcmd"

uninstall:
	rm $(BIN_DIR)/bootcfg
	rm $(BIN_DIR)/bootcmd

.PHONY: build test install
