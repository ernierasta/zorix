# Zorix

Remote server/service monitoring and alerting system deployed in 5 minutes. No agents, no dependencies.

## Features:

 - monitor web services,
 - run any command to run custom checks,
 - configurable notifications, templates, sending interval, recovery message ...
 - one file to deploy, scalable, minimal requirements.

## Status:

 It is under development, but already usable. Look at TODO list below to see what is planned.

## TODO:

- documentation,
- installation instructions,
- implement flags for test notification, config and log,
- init scripts for: sysVinit (for legacy installations), runit, OpenRC and yes ... systemd. ;-)
- Windows specific info/scripts (will it currently run?),
- implement better API checking (inspired by [statusOK](https://github.com/sanathp/statusok)),
- implement jabber (xmpp) notifications,
- check code for any panics, allow them only on process start, but not when it is running (should be ok already),
- add check types: ping and other shorthands,
- port testing,
- implement [rtop](https://github.com/rapidloop/rtop) funkcionality.
  Configure ssh access, set thresholds and you have remote system resources monitored (cpu, ram, hdd, ...).

