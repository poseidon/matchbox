
# Lifecycle of a Physical Machine

A physical machine [network boots](network-booting.md) in an network boot environment created by [coreos/dnsmasq](../contrib/dnsmasq) or a custom DHCP/TFTP/DNS setup.

`bootcfg` serves iPXE, GRUB, or Pixiecore boot configs via HTTP to machines matching attribute selectors (UUID, MAC, region, etc.). The referenced kernel and initrd images are fetched and booted with an initial Ignition config for installing CoreOS. CoreOS is installed to disk and the Ignition config for the machine is fetched from `bootcfg` before rebooting.

The CoreOS machine boots (first boot from disk) and runs its Ignition config to provision its disk with systemd units, files, keys, etc. On subsequent reboots, systemd units may fetch dynamic metadata if needed. Ignition is not run again.

CoreOS hosts should have automatic updates enabled and use a system like fleet or Kubernetes to run containers to tolerate node updates or failures without operator intervention. Use IPMI or vendor utilities to re-provision machines to change their role, rather than mutation.

![Machine Lifecycle](img/machine-lifecycle.png)

