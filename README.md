# Zorix

Remote server/service monitoring and alerting system deployed in 5 minutes. 
No agents, no dependencies, one binary.

## Features:

 - monitor web services,
 - run any command to run custom checks (cmd type check),
 - configurable notifications, templates, sending interval, recovery message ...
 - one file to deploy, scalable, minimal requirements.

## Status:

 It is under development, but already usable. Look at TODO list below to see what is planned.

### No tests?

Yes, not a line. I find testing very important. Some of my projects are even TDD (or almost TDD). But I needed this ASAP.
But do not worry, tests will come, adding tests will be good opportunity to refactor where needed. 

## TODO:

- documentation,
- tests (!),
- installation instructions,
- implement flags for test notification, config and log,
- init scripts for: sysVinit (for legacy installations), runit, OpenRC and yes ... systemd. ;-)
- Windows specific info/scripts (will it currently run?),
- implement better API checking (inspired by [statusOK](https://github.com/sanathp/statusok)),
- implement jabber (xmpp) notifications,
- implement cmd notifications (you can use any command as alerter),
- check code for any panics, allow them only on process start, but not when it is running (should be ok already),
- add check types: ping and other shorthands,
- port testing,
- implement [rtop](https://github.com/rapidloop/rtop) funkcionality.
  Configure ssh access, set thresholds and you have remote system resources monitored (cpu, ram, hdd, ...).

## Contributions

Contributions are very welcome! As well any bug reports are great way how can you improve this project.
But please wait for first binary release here on github Release page. When packages will be there it means
I want you to deploy it and documentation is in place. :-)


