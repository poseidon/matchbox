## gRPC TLS Generation

The Matchbox gRPC API allows clients (`terraform-provider-matchbox`) to create and update Matchbox resources. TLS credentials are used for client authentication and to establish a secure communication channel. When the gRPC API is [enabled](../../docs/deployment.md#customization), the server requires a TLS server certificate, key, and CA certificate ([locations](../../docs/config.md#files-and-directories)).

The `cert-gen` helper script generates a self-signed CA, server certificate, and client certificate. **Prefer your organization's PKI, if possible**

Navigate to the `scripts/tls` directory.

```sh
$ cd scripts/tls
```

Export `SAN` to set the Subject Alt Names which should be used in certificates. Provide the fully qualified domain name or IP (discouraged) where Matchbox will be installed.

```sh
# DNS or IP Subject Alt Names where matchbox runs
$ export SAN=DNS.1:matchbox.example.com,IP.1:172.18.0.2
```

Generate a `ca.crt`, `server.crt`, `server.key`, `client.crt`, and `client.key`.

```sh
$ ./cert-gen
Creating FAKE CA, server cert/key, and client cert/key...
...
...
...
******************************************************************
WARNING: Generated credentials are self-signed. Prefer your
organization's PKI for production deployments.
```

Move TLS credentials to the matchbox server's default location.

```sh
$ sudo mkdir -p /etc/matchbox
$ sudo cp ca.crt server.crt server.key /etc/matchbox
```

Save `client.crt`, `client.key`, and `ca.crt` for later use (e.g. `~/.matchbox`).

*If you are using the local Matchbox [development environment](../../docs/getting-started-docker.md), move server credentials to `examples/etc/matchbox`.*

## Inspect

Inspect the generated certificates if desired.

```sh
openssl x509 -noout -text -in ca.crt
openssl x509 -noout -text -in server.crt
openssl x509 -noout -text -in client.crt
```

## Verify

Verify that the server and client certificates were signed by the self-signed CA.

```sh
openssl verify -CAfile ca.crt server.crt
openssl verify -CAfile ca.crt client.crt
```
