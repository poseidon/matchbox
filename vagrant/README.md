# Vagrant Development

`pxe` and `pixiecore` provide Vagrantfiles and scripts for setting up a PXE or Pixiecore provisioning server in libvirt for development.

To get started, install the dependencies

	# Fedora 22/23
	dnf install vagrant vagrant-libvirt virt-manager

## Usage

Create a PXE or Pixiecore server VM with `vagrant up`.

    vagrant up --provider libivrt
    vagrant ssh

The PXE server will allocate DHCP leases, run a TFTP server with a CoreOS kernel image and init RAM fs, and host a cloud-config over HTTP. The Pixiecore server itself is a proxy DHCP, TFTP, and HTTP server for images.

By default, the PXE server runs at 192.168.32.10 on the `vagrant-pxe` virtual network. The Pixiecore server runs at 192.168.33.10 on the `vagrant-pixiecore` virtual network.

### Clients

Once the provisioning server has started, PXE boot enabled client VMs in the same network should boot with CoreOS.

Launch `virt-manager` to create a new virtual machine. When prompted, select Network Boot (PXE), skip adding a disk, and choose the `vagrant-libvirt` network.

If you see "Nothing" to boot, try force resetting the client VM.

Use SSH to connect to a client VM after boot and cloud-config succeed. The CLIENT_IP will be visible in the virt-manager console.

    ssh core@CLIENT_IP  # requires ssh_authorized_keys entry in cloud-config

### Configuration

The Vagrantfile parses the `config.rb` file for several variables you can use to configure network settings.

### Reload

If you change the Vagrantfile or a configuration variable, reload the VM with

    vagrant reload --provision

To try a new cloud-config, you can also scp the file onto the dev PXE server.

	scp new-config.yml core@NODE_IP:/var/www/html/cloud-config.yml
