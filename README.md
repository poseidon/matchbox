
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

Next, prepare a [data config directory](#configs) or use the example in [data](data). Optionally create an [image assets](#assets) directory.

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
            └── 1cff2cd8-f00a-42c8-9426-f55e6a1847f6

To find boot configs and cloud configs, the `FileStore` searches the `uuid` directory for a file matching a client machine's UUID, then searches `mac` for file matching the client's MAC address, and finally falls back to using the `default` file if present.

A typical boot config can be written as

    {
        "kernel": "/images/stable/coreos_production_pxe.vmlinuz",
        "initrd": ["/images/stable/coreos_production_pxe_image.cpio.gz"],
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

    ./scripts/get-coreos

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

If you're running `bootcfg` on Docker an easy way to try PXE/iPXE booting is by running the included `ipxe` container on the same Docker bridge.

The `ipxe` Docker image runs DHCP and sends options to point PXE/iPXE clients to the boot config service (hardcoded to http://172.17.0.2:8080). For PXE clients, iPXE firmware is chainloaded.

    cd dockerfiles/ipxe
    ./docker-build
    ./docker-run

Create a PXE boot client VM or attach a bare metal PXE boot client to the `docker0` virtual bridge. See [clients](#Clients).

### PXE

To use `bootcfg` with PXE, you must [chainload iPXE](http://ipxe.org/howto/chainloading). This can be done by configuring your PXE/DHCP/TFTP server to send `undionly.kpxe` over TFTP to PXE clients.

    # dnsmasq.conf
    enable-tftp
    tftp-root=/var/lib/tftpboot
    # if regular PXE firmware, serve iPXE firmware (via TFTP)
    dhcp-boot=tag:!ipxe,undionly.kpxe

`bootcfg` does not respond to DHCP requests or serve files over TFTP.

#### Pixecore

Run the Pixiecore server container which uses `api` mode to call through to the config service.

    make run-pixiecore

Finally, run the `vethdhcp` script to create a virtual ethernet connection on the `docker0` bridge, assign an IP address, and run dnsmasq to provide DHCP service to VMs we'll add to the bridge.

    make run-dhcp

Create a PXE boot client VM using virt-manager as described above.

### Troubleshooting

* Check your firewall settings to ensure you've allowed DHCP to run.
* On some platforms, SELinux prevents file serving unless contexts are changed appropriately.
* If you get an "address is already in use" error, try stopping the `default` network created in virt-manager.
* If you find that the `docker0` bridge receives DHCP Offers from dnsmasq, but the VM does not you may need to change an iptables rule.

Change subnet MASQUERADE.

    $ iptables -L -t nat
    $ iptables -t nat -R POSTROUTING 1 -s 172.17.0.0/16 ! -d 172.17.0.0/16 -j MASQUERADE

    POSTROUTING
    MASQUERADE  all  --  172.17.0.0/16        0.0.0.0/0      (Original Rule 1)
    MASQUERADE  all  --  172.17.0.0/16       !172.17.0.0/16  (Updated Rule 1)

## Clients

Once boot services are running, create a PXE boot client VM or attach a bare metal machine to your host.

### VM

Create a VM using the virt-manager UI, select Network Boot with PXE, and for the network selection, choose "Specify Shared Device" with bridge name `docker0`. The VM should PXE boot using the boot configuration determined by its MAC address, which can be tweaked in virt-manager.

### Bare Metal

Link a bare metal machine, which has boot firmware (BIOS) support for PXE, to your host with a network adapter. Get the link and attach it to the bridge.

    ip link show                      # find new link e.g. enp0s20u2
    brctl addif docker0 enp0s20u2

Configure the boot firmware to prefer PXE booting or network booting and restart the machine. It should PXE boot using the boot configuration determined by its MAC address.
