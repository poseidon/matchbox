
# CoreOS on Hardware

`bootcfg` serves configs to PXE/iPXE/Pixiecore network boot clients machines based on their hardware attributes, to declaratively setup clusters on virtual or physical hardware.

Boot configs (kernel, options) and cloud configs can be defined for machines with a specific UUID, a specific MAC address, or as the default. `bootcfg` serves configs as iPXE scripts and as JSON to implement the Pixiecore [API spec](https://github.com/danderson/pixiecore/blob/master/README.api.md).

`bootcfg` can run as a container or on a provisioner machine directly and operates alongside your existing DHCP and PXE boot setups.

## Usage

Pull the container image

    docker pull quay.io/dghubble/bootcfg:latest

or build the binary from source.

    ./build
    ./docker-build

Next, prepare a [data config directory](#configs) or use the example in [data](data). Optionally create an [image assets](#image assets) directory.

Run the application on your host or as a Docker container with flag or environment variable [arguments](docs/config.md).

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/data:/data -v $PWD/images:/images dghubble/bootcfg:latest -address=0.0.0.0:8080

You can quickly check that `/ipxe?uuid=val` and `/pixiecore/v1/boot/:mac` serve the expected boot config and `/cloud?uuid=val` serves the expected cloud config. Proceed to [Networking](#networking) for discussion about connecting clients to your new boot config service.

### Configs

A `Store` maintains associations between machine attributes and different types of bootstrapping configurations. Currently, `bootcfg` includes a `FileStore` which can search a filesystem directory for `boot` and `cloud` configs.

Prepare a config data directory. If you keep a versioned repository of declarative configs, consider keeping this directory there.

    data
    ├── boot
    │   └── default
    └── cloud
        ├── default
        └── uuid
            └── 1cff2cd8-f00a-42c8-9426-f55e6a1847f6
        └── mac
            └── 52:54:00:c7:b6:64

To find boot configs and cloud configs, the `FileStore` searches the `uuid` directory for a file matching a client machine's UUID, then searches `mac` for file matching the client's MAC address, and finally falls back to using the `default` file if present.

A typical boot config can be written as

    {
        "kernel": "/images/coreos/835.9.0/coreos_production_pxe.vmlinuz",
        "initrd": ["/images/coreos/835.9.0/coreos_production_pxe_image.cpio.gz"],
        "cmdline": {
            "cloud-config-url": "http://172.17.0.2:8080/cloud?uuid=${uuid}",
            "coreos.autologin": ""
        }
    }

Point kernel and initrd to either the URIs of images or to local [assets](#assets) served by `bootcfg`. If the OS in the boot config supports it, point `cloud-config-url` to the `/cloud` endpoint at the name or IP where you plan to run `bootcfg`.

For this example, we use the internal IP Docker will assign to the first container on its bridge, because we'll attach PXE clients to the same bridge.

A typical cloud config script:

    #cloud-config
    coreos:
      units:
        - name: etcd2.service
          command: start
        - name: fleet.service
          command: start
    write_files:
      - path: "/home/core/welcome"
        owner: "core"
        permissions: "0644"
        content: |
          File added by the default cloud-config.

See the [data](/data) directory for examples. Alternative `Store` backends and support for additional machine attributes is forthcoming.

### Image Assets

Optionally, `bootcfg` can serve free-form static assets (e.g. kernel and initrd images) if an `-images-path` argument to a mounted volume directory is provided.

    images/
    └── coreos
        └── 835.9.0
            ├── coreos_production_pxe.vmlinuz
            └── coreos_production_pxe_image.cpio.gz

Run the `get-coreos` script to quickly download kernel and initrd images from a recent CoreOS release into an `/images` directory.

    ./scripts/get-coreos                 # stable, 835.9.0
    ./scripts/get-coreos beta 877.1.0

To use the hosted images, tweak `kernel` and `initrd` in a boot config file. For example, change `http://stable.release.core-os.net/amd64-usr/current/coreos_production_pxe.vmlinuz` to `/images/coreos/835.9.0/coreos_production_pxe.vmlinuz` so your client machines don't all download a remote image.

## Networking

The boot config service can be validated in PXE, iPXE, and Pixiecore scenarios using a libvirt virtual network with VM or bare metal PXE clients.

An ethernet bridge can be used to connect container services and VM or bare metal boot clients on the same network. Docker starts containers with virtual ethernet connections to the `docker0` bridge

    $ brctl show
    $ docker network inspect bridge  # docker client names the bridge "bridge"

which uses the default subnet 172.17.0.0/16. Rather than reconfiguring Docker to start containers on a user defined bridge, it is easier to attach VMs and bare metal machines to `docker0`.

Start the container services specific to the scenario you wish to test, following the sections on PXE, iPXE, and Pixiecore below. Then create a client.

### iPXE

Configure your iPXE server to chainload the `bootcfg` iPXE boot script hosted at `$BOOTCFG_BASE_URL/boot.ipxe`. The base URL should be of the for `protocol://hostname:port` where where hostname is the IP address of the config server or a DNS name for it.

    # dnsmasq.conf
    # if request from iPXE, serve bootcfg's iPXE boot script (via HTTP)
    dhcp-boot=tag:ipxe,http://172.17.0.2:8080/boot.ipxe

#### docker0

Try running a PXE/iPXE server alongside `bootcfg` on the docker0 bridge. Use the included `ipxe` container to run an example PXE/iPXE server on that bridge.

The `ipxe` Docker image uses dnsmasq DHCP to point PXE/iPXE clients to the boot config service (hardcoded to http://172.17.0.2:8080). It also runs a TFTP server to serve the iPXE firmware to older PXE clients via chainloading.

    cd dockerfiles/ipxe
    ./docker-build
    ./docker-run

Now create local PXE boot [clients](#clients) as libvirt VMs or by attaching bare metal machines to the docker0 bridge.

### PXE

To use `bootcfg` with PXE, you must [chainload iPXE](http://ipxe.org/howto/chainloading). This can be done by configuring your PXE/DHCP/TFTP server to send `undionly.kpxe` over TFTP to PXE clients.

    # dnsmasq.conf
    enable-tftp
    tftp-root=/var/lib/tftpboot
    # if regular PXE firmware, serve iPXE firmware (via TFTP)
    dhcp-boot=tag:!ipxe,undionly.kpxe

`bootcfg` does not respond to DHCP requests or serve files over TFTP.

### Pixecore

Pixiecore is a ProxyDHCP, TFTP, and HTTP server and calls through to the `bootcfg` API to get a boot config for `pxelinux` to boot. No modification of your existing DHCP server is required in production.

#### docker0

Try running a DHCP server, Pixiecore, and `bootcfg` on the docker0 bridge. Use the included `dhcp` container to run an example DHCP server and the official Pixiecore container image.

    # DHCP
    cd dockerfiles/dhcp
    ./docker-build
    ./docker-run

Start Pixiecore using the script which attempts to detect the IP and port of `bootcfg` on the Docker host or do it manually.

    # Pixiecore
    ./scripts/pixiecore
    # manual
    docker run -v $PWD/images:/images:Z danderson/pixiecore -api http://$BOOTCFG_HOST:$BOOTCFG_PORT/pixiecore

Now create local PXE boot [clients](#clients) as libvirt VMs or by attaching bare metal machines to the docker0 bridge.

## Clients

Once boot services are running, create a PXE boot client VM or attach a bare metal machine to your host.

### VM

Create a VM using the virt-manager UI, select Network Boot with PXE, and for the network selection, choose "Specify Shared Device" with bridge name `docker0`. The VM should PXE boot using the boot configuration determined by its MAC address, which can be tweaked in virt-manager.

### Bare Metal

Link a bare metal machine, which has boot firmware (BIOS) support for PXE, to your host with a network adapter. Get the link and attach it to the bridge.

    ip link show                      # find new link e.g. enp0s20u2
    brctl addif docker0 enp0s20u2

Configure the boot firmware to prefer PXE booting or network booting and restart the machine. It should PXE boot using the boot configuration determined by its MAC address.
