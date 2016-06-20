
# Torus Storage

The `torus` example provisions a 3 node CoreOS cluster, with `etcd3` and Torus, to demonstrate a stand-alone storage cluster. Each of the 3 nodes runs a Torus instance which makes 1GiB of space available (configured per node by "torus_storage_size" in machine group metadata).

## Requirements

Ensure that you've gone through the [bootcfg with rkt](getting-started-rkt.md) guide and understand the basics. In particular, you should be able to:

* Use rkt to start `bootcfg`
* Create a network boot environment with `coreos/dnsmasq`
* Create the example libvirt client VMs
* Install the Torus [binaries](https://github.com/coreos/torus/releases)

## Examples

The [examples](..examples) statically assign IP addresses (172.15.0.21, 172.15.0.22, 172.15.0.23) to libvirt client VMs created by `scripts/libvirt`. You can use the same examples for real hardware, but you'll need to update the MAC/IP addresses.

* [torus](../examples/groups/torus) - iPXE boot a Torus cluster (use rkt)

## Assets

Download the CoreOS image assets referenced in the target [profile](../examples/profiles).

    ./scripts/get-coreos alpha 1053.2.0 ./examples/assets

## Containers

Run the latest `bootcfg` ACI with rkt and the `torus` example.

    sudo rkt run --net=metal0:IP=172.15.0.2 --mount volume=data,target=/var/lib/bootcfg --volume data,kind=host,source=$PWD/examples --mount volume=groups,target=/var/lib/bootcfg/groups --volume groups,kind=host,source=$PWD/examples/groups/torus quay.io/coreos/bootcfg:latest -- -address=0.0.0.0:8080 -log-level=debug

Create a network boot environment with `coreos/dnsmasq` and create VMs with `scripts/libvirt` as covered in [bootcfg with rkt](getting-started-rkt.md). Client machines should network boot and provision themselves.

## Verify

Install the Torus [binaries](https://github.com/coreos/torus/releases) on your laptop. Torus uses etcd3 for coordination and metadata storage, so any etcd node in the cluster can be queried with `torusctl`.

    ./torusctl --etcd 172.15.0.21:2379 list-peers

Run `list-peers` to report the status of data nodes in the Torus cluster.

```
+--------------------------+--------------------------------------+---------+------+--------+---------------+--------------+
|         ADDRESS          |                 UUID                 |  SIZE   | USED | MEMBER |    UPDATED    | REB/REP DATA |
+--------------------------+--------------------------------------+---------+------+--------+---------------+--------------+
| http://172.15.0.21:40000 | 016fad6a-2e23-11e6-8ced-525400a19cae | 1.0 GiB | 0 B  | OK     | 1 second ago  | 0 B/sec      |
| http://172.15.0.23:40000 | 0408cbba-2e23-11e6-9871-525400c36177 | 1.0 GiB | 0 B  | OK     | 2 seconds ago | 0 B/sec      |
| http://172.15.0.22:40000 | 0c67d31c-2e23-11e6-91f5-525400b22f86 | 1.0 GiB | 0 B  | OK     | 3 seconds ago | 0 B/sec      |
+--------------------------+--------------------------------------+---------+------+--------+---------------+--------------+
```

Torus has already initialized its metadata within etcd3 to format the cluster and added all peers to the pool. Each node provides 1 GiB of storage and has `MEMBER` status `OK`.

### Volume Creation

Create a new replicated, virtual block device or `volume` on Torus.

    ./torusblk --etcd=172.15.0.21:2379 volume create hello 500MiB

List the current volumes,

    ./torusctl --etcd=172.15.0.21:2379 volume list

and verify that `hello` was created.

```
+-------------+---------+
| VOLUME NAME |  SIZE   |
+-------------+---------+
| hello       | 500 MiB |
+-------------+---------+
```

### Filesystems and Mounting

Let's attach the Torus volume, create a filesystem, and add some files. Add the `nbd` kernel module.

    sudo modprobe nbd
    sudo ./torusblk --etcd=172.15.0.21:2379 nbd hello

In a new shell, create a new filesystem on the volume and mount it on your system.

    sudo mkfs.ext4 /dev/nbd0
    sudo mkdir -p /mnt/hello
    sudo mount /dev/nbd0 -o discard,noatime /mnt/hello

Check that the mounted filesystem is present.

    $ mount | grep nbd
    /dev/nbd0 on /mnt/hello type ext4 (rw,noatime,seclabel,discard,data=ordered)

By default, Torus uses a replication factor of 2. You may write some data and poweroff one of the three nodes if you wish.

    sudo sh -c "echo 'hello world' > /mnt/hello/world"
    sudo virsh destroy node3            # actually equivalent to poweroff

Check the Torus data nodes.

    $ ./torusctl --etcd 172.15.0.21:2379 list-peers

```
+--------------------------+--------------------------------------+---------+--------+--------+---------------+--------------+
|         ADDRESS          |                 UUID                 |  SIZE   |  USED  | MEMBER |    UPDATED    | REB/REP DATA |
+--------------------------+--------------------------------------+---------+--------+--------+---------------+--------------+
| http://172.15.0.21:40000 | 016fad6a-2e23-11e6-8ced-525400a19cae | 1.0 GiB | 22 MiB | OK     | 3 seconds ago | 0 B/sec      |
| http://172.15.0.22:40000 | 0c67d31c-2e23-11e6-91f5-525400b22f86 | 1.0 GiB | 22 MiB | OK     | 3 seconds ago | 0 B/sec      |
|                          | 0408cbba-2e23-11e6-9871-525400c36177 | ???     | ???    | DOWN   | Missing       |              |
+--------------------------+--------------------------------------+---------+--------+--------+---------------+--------------+
Balanced: true Usage:  2.15%
```

## Going Further

See the [Torus](https://github.com/coreos/torus) project to learn more about Torus and contribute.
