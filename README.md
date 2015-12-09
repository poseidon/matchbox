
# Boot Config Service

Boot config service for PXE, iPXE, and Pixiecore.

## Validation

The boot config service can be validated in PXE, iPXE, and Pixiecore scenarios using a libvirt virtual network with VM or bare metal PXE clients.

An ethernet bridge can be used to connect container services and VM or bare metal boot clients on the same network. Docker starts containers with virtual ethernet connections to the `docker0` bridge

    $ brctl show
    $ docker network inspect bridge  # docker client names the bridge "bridge"

which uses the default subnet 172.17.0.0/16. Rather than reconfiguring Docker to start containers on a user defined bridge, it is easier to attach VMs and bare metal machines to `docker0`.

Start the container services specific to the scenario you wish to test, following the sections on PXE, iPXE, and Pixiecore below. Then create a client.

### iPXE

#### Pixecore

Run the boot config service container.

    make run-docker

Run the Pixiecore server container which uses `api` mode to call through to the config service.

    make run-pixiecore

Finally, run the `vethdhcp` to create a virtual ethernet connection on the `docker0` bridge, assign an IP address, and run dnsmasq to provide DHCP service to VMs we'll add to the bridge.

    make run-dhcp

If you get an "address is already in use" error, try stopping the `default` network created in virt-manager. If you find that the `docker0` bridge receives DHCP Offers from dnsmasq, but the VM does not you may need to change an iptables rule.

    $ ip tables -L -t nat
    $ iptables -t nat -R POSTROUTING 1 -s 172.17.0.0/16 ! -d 172.17.0.0/16 -j MASQUERADE

    POSTROUTING
    MASQUERADE  all  --  172.17.0.0/16        0.0.0.0/0      (Original Rule 1)
    MASQUERADE  all  --  172.17.0.0/16       !172.17.0.0/16  (Updated Rule 1)

Create a PXE boot client VM using virt-manager as described above. Let the VM PXE boot using the boot config determined the MAC address to config mapping set in the config service.

## Clients

Once boot services are running, create a PXE boot client VM or attach a bare metal machine to your host.

### VM

Create a VM using the virt-manager UI, select Network Boot with PXE, and for the network selection, choose "Specify Shared Device" with bridge name `docker0`. The VM should PXE boot using the boot configuration determined by its MAC address, which can be tweaked in virt-manager.

### Bare Metal

Link a bare metal machine, which has boot firmware (BIOS) support for PXE, to your host with a network adapter. Get the link and attach it to the bridge.

    ip link show                      # find new link e.g. enp0s20u2
    brctl addif docker0 enp0s20u2

Configure the boot firmware to prefer PXE booting or network booting and restart the machine. It should PXE boot using the boot configuration determined by its MAC address.
