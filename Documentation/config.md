
# Flags and Variables

Configuration arguments can be provided as flags or as environment variables.

| flag | variable | example |
|------|----------|---------|
| -address | BOOTCFG_ADDRESS | 0.0.0.0:8080 |
| -config | BOOTCFG_CONFIG | ./data/config.yaml |
| -data-path | BOOTCFG_DATA_PATH | ./data |
| -assets-path | BOOTCFG_ASSETS_PATH | ./assets |
| -key-ring-path | BOOTCFG_KEY_RING_PATH | ~/.secrets/vault/bootcfg/secring.gpg |
| Disallowed | BOOTCFG_PASSPHRASE | secret passphrase |
| -log-level | BOOTCFG_LOG_LEVEL | critical, error, warning, notice, info, debug |

## Examples

Build the static binary.

    ./build

Run

    ./bin/bootcfg -address=0.0.0.0:8080 -log-level=debug -data-path examples/ -config examples/etcd-rkt.yaml

Run with a fake signing key.

    export BOOTCFG_PASSPHRASE=test
    ./bin/bootcfg -address=0.0.0.0:8080 -key-ring-path sign/fixtures/secring.gpg -data-path examples/ -config examples/etcd-rkt.yaml


