
# OpenPGP signing

The `matchbox` OpenPGP signature endpoints serve detached binary and ASCII armored signatures of rendered configs, if enabled. Each config endpoint has corresponding signature endpoints, typically suffixed with `.sig` or `.asc`.

To enable OpenPGP signing, provide the path to a secret keyring containing a single signing key with `-key-ring-path` or by setting `MATCHBOX_KEY_RING_PATH`. If a passphrase is required, set it via the `MATCHBOX_PASSPHRASE` environment variable.

Here are example signature endpoints without their query parameters.

| Endpoint   | Signature Endpoint | ASCII Signature Endpoint |
|------------|--------------------|-------------------------|
| iPXE       | `http://matchbox.foo/ipxe.sig` | `http://matchbox.foo/ipxe.asc` |
| GRUB2      | `http://bootcf.foo/grub.sig` | `http://matchbox.foo/grub.asc` |
| Ignition   | `http://matchbox.foo/ignition.sig` | `http://matchbox.foo/ignition.asc` |
| Cloud-Config | `http://matchbox.foo/cloud.sig` | `http://matchbox.foo/cloud.asc` |
| Metadata   | `http://matchbox.foo/metadata.sig` | `http://matchbox.foo/metadata.asc` |

In production, mount your signing keyring and source the passphrase from a [Kubernetes secret](https://kubernetes.io/docs/user-guide/secrets/). Use a signing subkey exported to a keyring by itself, which can be revoked by a primary key, if needed.

To try it locally, you may use the test fixture keyring. **Warning: The test fixture keyring is for examples only.**

## Verify

Verify a signature response and config response from the command line using the public key. Notice that most configs have a trailing newline.

**Warning: The test fixture keyring is for examples only.**

```sh
$ gpg --homedir sign/fixtures --verify sig_file response_file
gpg: Signature made Mon 08 Feb 2016 11:37:03 PM PST using RSA key ID 9896356A
gpg: sign/fixtures/trustdb.gpg: trustdb created
gpg: Good signature from "Fake Bare Metal Key (Do not use) <do-not-use@example.com>"
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
Primary key fingerprint: BE2F 12BC 3642 2594 570A  CCBB 8DC4 2020 9896 356A
```

## Signing key generation

Create a signing key or subkey according to your requirements and security policies. Here are some basic [guides](https://coreos.com/rkt/docs/latest/signing-and-verification-guide.html).

### gpg

```sh
$ mkdir -m 700 path/in/vault
$ gpg --homedir path/in/vault --expert --gen-key
...
```

### gpg2

```sh
$ mkdir -m 700 path/in/vault
$ gpg2 --homedir path/in/vault --expert --gen-key
...
$ gpg2 --homedir path/in/vault --export-secret-key KEYID > path/in/vault/secring.gpg
```
