
BIN_DIR=/usr/local/bin
CONF_DIR=/etc/bootcfg
VAR_DIR=/var/bootcfg

CONF_FILE=/etc/bootcfg.conf
ENV_FILE=/etc/bootcfg.env

all: build

build:
	./build

test:
	./test

install:
	cp bin/bootcfg $(BIN_DIR)
	cp -n examples/default.yaml $(CONF_FILE)
	touch $(ENV_FILE)
	mkdir -p $(CONF_DIR)/{profiles,ignition,cloud}
	mkdir -p $(VAR_DIR)
	cp -n -R examples/profiles $(CONF_DIR)
	cp -n -R examples/ignition $(CONF_DIR)
	cp -n -R examples/cloud $(CONF_DIR)
	@echo "*****************"
	@echo "bootcfg INSTALLED"
	@echo "*****************"
	@echo "bootcfg was installed to /usr/local/bin/bootcfg"
	@echo "The config file is located at /etc/bootcfg.conf"
	@echo "The environment file is located at /etc/bootcfg.env"

uninstall:
	rm $(BIN_DIR)/bootcfg
	rm $(CONF_FILE)
	rm $(ENV_FILE)
	rm -rf $(CONF_DIR)
	rm -rf $(VAR_DIR)

.PHONY: build
