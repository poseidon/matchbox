# Vagrant Boot Servers

The `pxe`, `ipxe`, and `pixiecore` Vagrantfiles setup example PXE, iPXE, or Pixiecore boot/provisioner servers which can each be used to boot libvirt VM clients on a shared network into CoreOS and provision them with a simple cloud-config. This illustrates how the different network boot server setups work.

To get started, install the dependencies

	# Fedora 22/23
	dnf install vagrant vagrant-libvirt virt-manager

## Usage

Select one of the boot servers and create a boot server VM with `vagrant up`.

    vagrant up --provider libvirt
    vagrant ssh

The **PXE server** uses dnsmasq for DHCP and TFTP and an HTTP server. DHCP grants authoritative DHCP leases on 192.168.32.0/24 and the boot server has static IP 192.168.32.10. TFTP serves the `pxelinux.0` bootloader, default pxelinux cfg, kernel image, and init RAM filesystem image. The HTTP server hosts a cloud config with a configurable authorized SSH key.

The **iPXE server** uses dnsmasq for DHCP and TFTP and an HTTP server. DHCP grants authoritative DHCP leases on 192.168.34.0/24 and the boot server has static IP 192.168.34.10. TFTP serves the `undionly.kpxe` bootloader. The HTTP server hosts a boo.ipxe config script, the kernel image, the init RAM filesystem, and a cloud config with a configurable authorized SSH key.

The **Pixiecore server** itself is a proxy DHCP server, TFTP server, and HTTP server for `lpxelinux.0`, the kernel image, and init RAM filesystem image. The network is configured to grant DHCP leases in 192.168.33.0/24 and the boot server has static IP address 192.168.33.10. A standalone HTTP server is used to serve the cloud-config with a configurable authorized SSH key.

 and will grant DHCP leases, run a TFTP server with a CoreOS kernel image and init RAM fs, and host a cloud-config over HTTP. 

### Configuration

The Vagrantfile parses the `config.rb` file for several configurable variables including

* network_range
* server_ip
* dhcp_range
* ssh_authorized_keys

### Clients

Any of the boot servers allow PXE boot enabled client VMs in the same network to boot into CoreOS and configure themselves with cloud-config.

Launch `virt-manager` to create a new virtual machine. When prompted, select Network Boot (PXE), skip adding a disk, and choose the `vagrant-pxe`, `vagrant-ipxe`, or `vagrant-pixiecore` network.

If you see "Nothing to boot", try force resetting the client VM, there can be DHCP contention on Vagrant.

Use SSH to connect to a client VM after boot and cloud-config succeed. The CLIENT_IP will be visible in the virt-manager console.

    ssh core@CLIENT_IP  # requires ssh_authorized_keys entry in cloud-config

### Reload

If you change the Vagrantfile or a configuration variable, reload the VM with

    vagrant reload --provision

To try a new cloud-config, you can also scp the file onto the dev PXE server.

	scp new-config.yml core@NODE_IP:/var/www/html/cloud-config.yml
