# Examples

Matchbox automates network booting and provisioning of clusters. These examples show how to use Matchbox on-premise or locally with QEMU/KVM.

## Terraform Examples

These examples use [Terraform](https://www.terraform.io/intro/) as a client to Matchbox.

| Name                          | Description                   |
|-------------------------------|-------------------------------|
| [fedora-coreos-install](terraform/fedora-coreos-install) | Fedora CoreOS disk install |
| [flatcar-install](terraform/flatcar-install) | Flatcar Linux disk install |

### Customization

Look through the examples and Terraform modules and use them as a starting point. Learn more about [matchbox](../docs/matchbox.md).

## Manual Examples

These examples mount raw Matchbox objects into a Matchbox server's `/var/lib/matchbox/` directory.

| Name          | Description                  | FS  | Docs  |
|---------------|------------------------------|-----|-------|
| fedora-coreos | Fedora CoreOS live PXE       | RAM | [docs](https://docs.fedoraproject.org/en-US/fedora-coreos/live-booting/) |
| fedora-coreos-install | Fedora CoreOS install | Disk | [docs](https://docs.fedoraproject.org/en-US/fedora-coreos/bare-metal/) |
| flatcar       | Flatcar Linux live PXE       | RAM | [docs](https://docs.flatcar-linux.org/os/booting-with-ipxe/) |
| flatcar-install | Flatcar Linux install      | Disk | [docs](https://docs.flatcar-linux.org/os/booting-with-ipxe/) |

### Customization

For Fedora CoreOS, add an SSH authorized key to Fedora CoreOS Config (`ignition/fedora-coreos.yaml`) and regenerate the Ignition Config.

```yaml
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

For Flatcar Linux, add a Matchbox variable to a Group (`groups/flatcar-install/flatcar.json`) to set the SSH authorized key (or directly update the Container Linux Config).

```json
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
