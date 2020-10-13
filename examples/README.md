# Examples

Matchbox automates network booting and provisioning of clusters. These examples show how to use matchbox on-premise or locally with [QEMU/KVM](scripts/README.md#libvirt).

## Terraform Examples

These examples use [Terraform](https://www.terraform.io/intro/) as a client to Matchbox.

| Name                          | Description                   |
|-------------------------------|-------------------------------|
| [simple-install](terraform/simple-install/) | Install Container Linux with an SSH key |
| [etcd3-install](terraform/etcd3-install/) | Install a 3-node etcd3 cluster |

### Customization

Look through the examples and Terraform modules and use them as a starting point. Learn more about [matchbox](../docs/matchbox.md) and [Container Linux configs](../docs/container-linux-config.md).

## Manual Examples

These examples mount raw Matchbox objects into a Matchbox server's `/var/lib/matchbox/` directory.

| Name          | Description                  | FS  | Docs  |
|---------------|------------------------------|-----|-------|
| fedora-coreos | Fedora CoreOS live PXE       | RAM | [docs](https://docs.fedoraproject.org/en-US/fedora-coreos/live-booting-ipxe/) |
| fedora-coreos-install | Fedora CoreOS install | Disk | [docs](https://docs.fedoraproject.org/en-US/fedora-coreos/bare-metal/) |
| flatcar       | Flatcar Linux live PXE       | RAM | [docs](https://docs.flatcar-linux.org/os/booting-with-ipxe/) |
| flatcar-install | Flatcar Linux install      | Disk | [docs](https://docs.flatcar-linux.org/os/booting-with-ipxe/) |

### Customization

For Fedora CoreOS, add an SSH authorized key to Fedora CoreOS Config (`ignition/fedora-coreos.yaml`) and regenerate the Ignition Config.

```
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-rsa pub-key-goes-here
```

```
podman run -i --rm quay.io/coreos/fcct:release --pretty --strict < fedora-coreos.yaml > fedora-coreos.ign
```

For Flatcar Linux, add a Matchbox variable to a Group to set the SSH authorized key (or directly update the Container Linux Config).

```
# groups/flatcar-install/flatcar.json
{
  "id": "stage-1",
  "name": "Flatcar Linux",
  "profile": "flatcar",
  "selector": {
    "os": "installed"
  },
  "metadata": {
    "ssh_authorized_keys": ["ssh-rsa pub-key-goes-here"]
  }
}
```

