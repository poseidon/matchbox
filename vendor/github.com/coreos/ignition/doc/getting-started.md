# Getting Started with Ignition

*Ignition* is a low-level system configuration utility. The Ignition executable is part of the temporary initial root filesystem, the *initramfs*. When Ignition runs, it finds configuration data in a named location for a given environment, such as a file or URL, and applies it to the machine before `switch_root` is called to pivot to the machine's root filesystem.

Ignition uses a JSON configuration file to represent the set of changes to be made. The format of this config is detailed [in the specification][configspec]. One of the most important parts of this config is the version number. This **must** match the version number accepted by Ignition. If the config version isn't accepted by Ignition, Ignition will fail to run and prevent the machine from booting. This can be seen by inspecting the console output of the failed instance. For more information, check out the [troubleshooting section][troubleshooting].

## Providing a Config

Ignition will choose where to look for configuration based on the underlying platform. A list of [supported platforms][platforms] and metadata sources is provided for reference.

The configuration must be passed to Ignition through the designated data source. Please refer to Ignition [config examples][examples] to learn about writing config files. The provided configuration will be appended to the universal base configuration:

```json
{
  "storage": {
    "filesystems": [{
      "name": "root",
      "path": "/sysroot"
    }]
  }
}
```

## Troubleshooting

The single most useful piece of information needed when troubleshooting is the log from Ignition. Ignition runs in multiple stages so it's easiest to filter by the syslog identifier: `ignition`. When using systemd, this can be accomplished with the following command:

```
journalctl --identifier=ignition
```

In the event that this doesn't yield any results, running as root may help. There are circumstances where the journal isn't owned by the systemd-journal group or the current user is not a part of that group.

In the vast majority of cases, it will be immediately obvious why Ignition failed. If it's not, inspect the config that Ignition wrote into the log. This shows how Ignition interpreted the supplied configuration. The user-provided config may have a misspelled section or maybe an incorrect hierarchy.

[configspec]: configuration.md
[examples]: https://github.com/coreos/docs/blob/master/ignition/examples.md
[platforms]: supported-platforms.md
[troubleshooting]: #troubleshooting
