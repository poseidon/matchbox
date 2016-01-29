
# Scripts

## Ignition Generation

Transform human friendly `*.yaml` Ignition files to machine-friendly Ignition configs using the [github.com/coreos/fuze](https://github.com/coreos/fuze) utility.

First, install Fuze.

    cd $GOPATH/src/github.com/coreos
    git clone https://github.com/coreos/fuze.git
    cd fuze
    ./build
    cp bin/fuze $GOPATH/bin

Use `gen-ignition` to generate a JSON Ignition config for each `*.yaml` file in a directory of Ignition file sources.

    ./scripts/gen-ignition examples/dev/ignition
