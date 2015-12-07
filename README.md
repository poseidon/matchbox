
Boot and cloud config service for PXE, iPXE, and Pixiecore.

## Validation

The config service can be validated in scenarios which use PXE, iPXE, or Pixiecore on a libvirt virtual network or on a physical network of bare metal machines. 

### libvirt

A libvirt virtual network of containers and VMs can be used to validate PXE booting of VM clients in various scenarios (PXE, iPXE, Pixiecore).

To do this, start the appropriate set of containered services for the scenario, then boot a VM configured to use the PXE boot method.

Docker starts containers with virtual ethernet connections to the `docker0` bridge

    $ brctl show
    $ docker network inspect bridge  # docker client names the bridge "bridge"

which uses the default subnet 172.17.0.0/16. It is also possible to create your own network bridge and reconfigure Docker to start containers on that bridge, but that approach is not used here.

PXE boot client VMs can be started within the same subnet by attaching to the `docker0` bridge.

Create a VM using the virt-manager UI, select Network Boot with PXE, and for the network selection, choose "Specify Shared Device" with bridge name `docker0`.

The VM should PXE boot using the boot config determined by the MAC address of the virtual network card which can be inspected in virt-manager.

### iPXE

#### Pixecore

Run the config service container.

    make run-docker

Run the Pixiecore server container which uses `api` mode to call through to the config service.

    make run-pixiecore

Finally, run the `vethdhcp` to create a virtual ethernet connection on the `docker0` bridge, assign an IP address, and run dnsmasq to provide DHCP service to VMs we'll add to the bridge.

    make run-dhcp

Create a PXE boot client VM using virt-manager as described above. Let the VM PXE boot using the boot config determined the MAC address to config mapping set in the config service.

