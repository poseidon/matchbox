# Terraform Modules

Matchbox provides Terraform [modules](https://www.terraform.io/docs/modules/usage.html) you can re-use directly within your own Terraform configs. Modules are updated regularly so it is **recommended** that you pin the module version (e.g. `ref=sha`) to keep your configs deterministic.

```hcl
module "profiles" {
  source = "git::https://github.com/coreos/matchbox.git//examples/terraform/modules/profiles?ref=4451425db8f230012c36de6e6628c72aa34e1c10"
  matchbox_http_endpoint = "${var.matchbox_http_endpoint}"
  container_linux_version = "${var.container_linux_version}"
  container_linux_channel = "${var.container_linux_channel}"
}
```

Download referenced Terraform  modules.

```sh
$ terraform get            # does not check for updates
$ terraform get --update   # checks for updates
```

Available modules:

| Module   | Includes  | Description |
|----------|-----------|-------------|
| profiles | *         | Creates machine profiles you can reference in matcher groups |
|          | container-linux-install | Install Container Linux to disk from core-os.net |
|          | cached-container-linux-install | Install Container Linux to disk from matchbox assets cache |
|          | etcd3    | Provision an etcd3 peer node |
|          | etcd3-gateway | Provision an etcd3 gateway node |

## Customization

You are encouraged to look through the examples and modules. Implement your own profiles or package them as modules to meet your needs. We've just provided a starting point. Learn more about [matchbox](../../Documentation/matchbox.md) and [Container Linux configs](../../Documentation/container-linux-config.md).
