
`pxe` provides a Vagrantfile and scripts for setting up a PXE server in libvirt or on physical hardware.

`pixiecore` provides a Vagrantfile and scripts for setting up a Pixiecore server in libvirt or on physical hardware.

## Setup

To develop with Vagrant, install the dependencies

	# Fedora 22/23
	dnf install vagrant vagrant-libvirt virt-manager

## Usage

The Vagrantfile will setup a `pxe_default` VM running a PXE server with a configured static IP address, DHCP range, CoreOS kernel image, and cloud-config. The VM will be connected to a network called `vagrant-pxe`.

### libvirt Provider

    vagrant up --provider libivrt
    vagrant ssh

Once the PXE server has started, you can start client VMs within the `vagrant-libvirt` network which should boot as PXE clients.

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
