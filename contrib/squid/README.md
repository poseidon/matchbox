# Squid Proxy (DRAFT)

This guide shows how to setup a [Squid](http://www.squid-cache.org/) cache proxy for providing kernel/initrd files to PXE, iPXE, or GRUB2 client machines. This setup runs Squid as a Docker container using the [sameersbn/squid](https://quay.io/repository/sameersbn/squid)
image.

The Squid container requires a squid.conf file to run. Download the example squid.conf file from the [sameersbn/docker-squid](https://github.com/sameersbn/docker-squid) repo:
```
curl -O https://raw.githubusercontent.com/sameersbn/docker-squid/master/squid.conf
```

Squid [interception caching](http://wiki.squid-cache.org/SquidFaq/InterceptionProxy#Concepts_of_Interception_Caching) is required for proxying PXE, iPXE, or GRUB2 client machines. Set the intercept mode in squid.conf:
```
sed -ie 's/http_port 3128/http_port 3128 intercept/g' squid.conf
```

By default, Squid caches objects that are 4MB or less. Increase the maximum object size to cache large files such as kernel and initrd images. The following example increases the maximum object size to 300MB:
```
sed -ie 's/# maximum_object_size 4 MB/maximum_object_size 300 MB/g' squid.conf
```

Squid supports a wide range of cache configurations. Review the Squid [documentation](http://www.squid-cache.org/Doc/) to learn more about configuring Squid.

This example uses systemd to manage squid. Create the squid service systemd unit file:
```
cat /etc/systemd/system/squid.service
#/etc/systemd/system/squid.service
[Unit]
Description=squid proxy service
After=docker.service
Requires=docker.service

[Service]
Restart=always
TimeoutStartSec=0
ExecStart=/usr/bin/docker run --net=host --rm \
          -v /path/to/squid.conf:/etc/squid3/squid.conf:Z \
          -v /srv/docker/squid/cache:/var/spool/squid3:Z \
          quay.io/sameersbn/squid

[Install]
WantedBy=multi-user.target
```

Start Squid:
```
systemctl start squid
```

If your Squid host is running iptables or firewalld, modify rules to allow the interception and redirection of traffic. In the following example, 192.168.10.1 is the IP address of the interface facing PXE, iPXE, or GRUB2 client machines. The default port number used by squid is 3128.

For firewalld:
```
firewall-cmd --permanent --zone=internal --add-forward-port=port=80:proto=tcp:toport=3128:toaddr=192.168.10.1
firewall-cmd --permanent --zone=internal --add-port=3128/tcp
firewall-cmd --reload
firewall-cmd --zone=internal --list-all
```

For iptables:
```
iptables -t nat -A POSTROUTING -o enp15s0 -j MASQUERADE
iptables -t nat -A PREROUTING -i enp14s0 -p tcp --dport 80 -j REDIRECT --to-port 3128
```
**Note**: enp14s0 faces PXE, iPXE, or GRUB2 clients and enp15s0 faces Internet access.

Your DHCP server should be configured so the Squid host is the default gateway for PXE, iPXE, or GRUB2 clients. For deployments that run Squid on the same host as dnsmasq, remove any DHCP option 3 settings. For example ```--dhcp-option=3,192.168.10.1"```

Update Matchbox policies to use the url of the Container Linux kernel/initrd download site:
```
cat policy/etcd3.json
{
  "id": "etcd3",
  "name": "etcd3",
  "boot": {
    "kernel": "http://stable.release.core-os.net/amd64-usr/1235.9.0/coreos_production_pxe.vmlinuz",
    "initrd": ["http://stable.release.core-os.net/amd64-usr/1235.9.0/coreos_production_pxe_image.cpio.gz"],
    "args": [
      "coreos.config.url=http://matchbox.foo:8080/ignition?uuid=${uuid}&mac=${mac:hexhyp}",
      "coreos.first_boot=yes",
      "console=tty0",
      "console=ttyS0",
      "coreos.autologin"
    ]
  },
  "ignition_id": "etcd3.yaml"
}
```

(Optional) Configure Matchbox to not serve static assets by providing an empty assets-path value.
```
# /etc/systemd/system/matchbox.service.d/override.conf
[Service]
Environment="MATCHBOX_ASSETS_PATHS="
```

Boot your PXE, iPXE, or GRUB2 clients.
