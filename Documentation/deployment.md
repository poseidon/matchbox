
# Deployment

## Binary

Clone the coreos-baremetal project into your $GOPATH.

    go get github.com/coreos/coreos-baremetal
    cd $GOPATH/src/github.com/coreos/coreos-baremetal

Build `bootcfg` from source.

    make

Install the `bootcfg` static binary to `/usr/local/bin`.

    $ sudo make install

Run `bootcfg`

    $ sudo bootcfg -version
    $ sudo bootcfg -address 0.0.0.0:8080
    main: starting bootcfg HTTP server on 0.0.0.0:8080

See [flags and variables](config.md).

### systemd

Add and start bootcfg's example systemd unit.

    sudo cp contrib/systemd/bootcfg.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl start bootcfg.service

Check the logs with `journalctl`.

    journalctl -u bootcfg.service

Enable the `bootcfg` service if you'd like it to start at boot time.

    sudo systemctl enable bootcfg.service

### Uninstall

    sudo systemctl stop bootcfg.service
    sudo make uninstall