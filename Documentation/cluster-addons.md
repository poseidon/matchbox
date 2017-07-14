## Cluster Addons

Kubernetes clusters run cluster addons atop Kubernetes itself. Addons may be considered essential for bootstrapping (non-optional), important (highly recommended), or optional.

## Essential

Several addons are considered essential. CoreOS cluster creation tools ensure these addons are included. Kubernetes clusters deployed via the Matchbox examples or using our Terraform Modules include these addons as well.

### kube-proxy

`kube-proxy` is deployed as a DaemonSet.

### kube-dns

`kube-dns` is deployed as a Deployment.

## Important

### Container Linux Update Operator

The [Container Linux Update Operator](https://github.com/coreos/container-linux-update-operator) (i.e. CLUO) coordinates reboots of auto-updating Container Linux nodes so that one node reboots at a time and nodes are drained before reboot. CLUO enables the auto-update behavior Container Linux clusters are known for, but does it in a Kubernetes native way. Deploying CLUO is strongly recommended.

Create the `update-operator` deployment and `update-agent` DaemonSet.

```
kubectl apply -f examples/addons/cluo/update-operator.yaml
kubectl apply -f examples/addons/cluo/update-agent.yaml
```

*Note, CLUO replaces `locksmithd` reboot coordination. The `update_engine` systemd unit on hosts still performs the Container Linux update check, download, and install to the inactive partition.*
