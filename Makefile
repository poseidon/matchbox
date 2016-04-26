
BIN_DIR=/usr/local/bin
DATA_DIR=/var/lib/bootcfg
ENV_FILE=/etc/bootcfg.env

all: build

build:
	./build

test:
	./test

install:
	cp bin/bootcfg $(BIN_DIR)
	cp bin/bootcmd $(BIN_DIR)
	mkdir -p $(DATA_DIR)/{profiles,groups,ignition,cloud,assets}
	cp -n -R examples/profiles $(DATA_DIR)
	cp -n -R examples/groups $(DATA_DIR)
	cp -n -R examples/ignition $(DATA_DIR)
	cp -n -R examples/cloud $(DATA_DIR)
	touch ${ENV_FILE}
	@echo "**************"
	@echo "INSTALL SUCESS"
	@echo "**************"
	@echo "bootcfg was installed to /usr/local/bin/bootcfg"
	@echo "bootcmd was installed to /usr/local/bin/bootcmd"
	@echo "The default data directory is located at /var/lib/bootcfg"

uninstall:
	rm $(BIN_DIR)/bootcfg
	rm $(BIN_DIR)/bootcmd

.PHONY: build
