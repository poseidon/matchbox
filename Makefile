LOCAL_OS := $(shell uname)

ifeq ($(LOCAL_OS),Linux)
    BIN_DIR=/usr/local/bin
    DATA_DIR=/var/lib/bootcfg
    ENV_FILE=/etc/bootcfg.env
else
    ifeq ($(LOCAL_OS),Darwin)
        BIN_DIR=/usr/local/bin
        DATA_DIR=~/.bootcfg
        ENV_FILE=~/.bootcfg/env
    else
        $(error Local OS not supported)
    endif
endif

all: build

build:
	./build

test:
	./test

install:
	mkdir -p $(DATA_DIR)/{profiles,groups,ignition,cloud,assets}
	touch $(ENV_FILE)
	cp bin/bootcfg $(BIN_DIR)
	cp bin/bootcmd $(BIN_DIR)
	@echo "**************"
	@echo "INSTALL SUCCESS"
	@echo "**************"
	@echo "bootcfg was installed to $(BIN_DIR)/bootcfg"
	@echo "bootcmd was installed to $(BIN_DIR)/bootcmd"

uninstall:
	rm $(BIN_DIR)/bootcfg
	@echo "bootcfg was removed from $(BIN_DIR)"
	rm $(BIN_DIR)/bootcmd
	@echo "bootcmd was removed from $(BIN_DIR)"

.PHONY: build test install
