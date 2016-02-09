
# OpenPGP Signing

The `bootcfg` OpenPGP signature endpoints serve ASCII armored detached signatures of rendered configs, if enabled. Each config endpoint has a corresponding signature endpoint, usually suffixed by `.asc`.

To enable OpenPGP signing, provide the path to a secret keyring containing a single signing key with `-key-ring-path` or by setting `BOOTCFG_KEY_RING_PATH`. If a passphrase is required, set it via the `BOOTCFG_PASSPHRASE` environment variable.

Here are example signature endpoints without their query parameters.

| Endpoint   | ASCII Signature Endpoint |
|------------|-----------------|
| Ignition   | `http://bootcfg.foo/ignition.asc` |
| Cloud-init | `http://bootcfg.foo/cloud.asc` |
| iPXE       | `http://bootcfg.foo/boot.ipxe.asc` |
| iPXE       | `http://bootcfg.foo/ipxe.asc` |
| Pixiecore  | `http://bootcfg.foo/pixiecore/v1/boot.asc/:MAC` |

In production, mount your signing keyring and source the passphrase from a [Kubernetes secret](http://kubernetes.io/v1.1/docs/user-guide/secrets.html). Use a signing subkey exported to a keyring used only for config signing, which can be revoked by a master if needed.

To try it locally, you may use the test fixture keyring. **Warning: The test fixture keyring is for examples only.**

**Binary**

    export BOOTCFG_PASSPHRASE=test
    ./bin/bootcfg -address=0.0.0.0:8080 -key-ring-path sign/fixtures/secring.gpg -config examples/etcd-rkt.yaml -data-path examples/

**rkt**

    sudo rkt run --set-env=BOOTCFG_PASSPHRASE=test --mount volume=secrets,target=/secrets --volume secrets,kind=host,source=$PWD/sign/fixtures --mount volume=assets,target=/assets --volume assets,kind=host,source=$PWD/assets --mount volume=data,target=/data --volume data,kind=host,source=$PWD/examples quay.io/coreos/bootcfg -- -address=0.0.0.0:8080 -config /data/etcd-rkt.yaml -key-ring-path secrets/secring.gpg

**docker**

    sudo docker run -p 8080:8080 --rm --env BOOTCFG_PASSPHRASE=test -v $PWD/examples:/data:Z -v $PWD/assets:/assets:Z -v $PWD/sign/fixtures:/secrets:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -config=/data/etcd-docker.yaml -key-ring-path secrets/secring.gpg

## Verify

Verify a signature response and config response from the command line using the public key. Notice that most configs have a trailing newline.

**Warning: The test fixture keyring is for examples only.**

    $ gpg --homedir sign/fixtures --verify sig_file response_file
    gpg: Signature made Mon 08 Feb 2016 11:37:03 PM PST using RSA key ID 9896356A
    gpg: sign/fixtures/trustdb.gpg: trustdb created
    gpg: Good signature from "Fake Bare Metal Key (Do not use) <do-not-use@example.com>"
    gpg: WARNING: This key is not certified with a trusted signature!
    gpg:          There is no indication that the signature belongs to the owner.
    Primary key fingerprint: BE2F 12BC 3642 2594 570A  CCBB 8DC4 2020 9896 356A

## Signing Key Generation

Create a signing key or subkey according to your requirements and security policies. Here are some basic [guides](https://coreos.com/rkt/docs/latest/signing-and-verification-guide.html).

### gpg

    mkdir -m 700 path/in/vault
    gpg --homedir path/in/vault --expert --gen-key
    ...

### gpg2

    mkdir -m 700 path/in/vault
    gpg2 --homedir path/in/vault --expert --gen-key
    ...
    gpg2 --homedir path/in/vault --export-secret-key KEYID > path/in/vault/secring.gpg

