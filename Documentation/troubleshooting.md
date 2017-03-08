# Troubleshooting

## Firewall

Running DHCP or proxyDHCP with `coreos/dnsmasq` on a host requires that the Firewall allow DHCP and TFTP (for chainloading) services to run.

## Port collision

Running DHCP or proxyDHCP can cause port already in use collisions depending on what's running. Fedora runs bootp listening on udp/67 for example. Find the service using the port.

```sh
$ sudo lsof -i :67
```

Evaluate whether you can configure the existing service or whether you'd like to stop it and test with `coreos/dnsmasq`.

## No boot filename received

PXE client firmware did not receive a DHCP Offer with PXE-Options after several attempts. If you're using the `coreos/dnsmasq` image with `-d`, each request should log to stdout. Using the wrong `-i` interface is the most common reason DHCP requests are not received. Otherwise, wireshark can be useful for investigating.
