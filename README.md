# Zorix

Remote server/service monitoring and alerting system deployed in 5 minutes. 
No agents, no dependencies, one binary.

## Features:

 - monitor web services,
 - launch command as custom checks (cmd type check),
 - configurable notifications, templates, sending interval, recovery message ...,
 - separate timers for checks and for notifications,
 - one file to deploy, scalable, minimal requirements,
 - configured in few minutes.

## Status:

 Zorix is under development, but usable. Look at TODO list below to see planned features.

### No tests?

Yes, not a line. I find testing very important. Some of my projects are TDD (or almost TDD). I needed this ASAP.
But do not worry, tests will come, adding tests will be good opportunity to refactor where needed. 

## TODO:

- [ ] documentation,
- [ ] tests (!),
- [ ] installation instructions,
- [ ] implement flags for test notification, config and log,
- [ ] init scripts for: sysVinit (for legacy installations), runit, OpenRC and yes ... systemd. ;-)
-  [ ] Windows specific info/functions (will it currently run?):
  - [ ] run as service,
  - [ ] hide cmd,
  - [ ] ping on Windows.
- [ ] implement better API checking (inspired by [statusOK](https://github.com/sanathp/statusok)),
- [ ] implement jabber (xmpp) notifications,
- [ ] implement cmd notifications (you can use any command as alerter),
- [ ] add [[global]] section to config, allow defining notification templates there,
- [ ] check code for any panics, allow them only on process start, but not when it is running (should be ok already),
- [ ] add check types: ping and other helper types to make config clearer,
- [ ] port testing,
- [ ] database storage: influxdb, maybe more if needed,
- [ ] document db usage, grafana integration,
- [ ] implement [rtop](https://github.com/rapidloop/rtop) funkcionality.
  Configure ssh access, set thresholds and you have remote system resources monitored (cpu, ram, hdd, ...).

## Contributions

Contributions are very welcome! 

Bug reports are great way how can you improve this project. So please if something is not working as expected, create new issue.

But please wait for first binary release here on github Release page. When packages will be there it means
I want you to deploy it and documentation is in place. :-)

## Additional info

### Why another monitoring solution?

I have servers, I needed to monitor them and I have 128MB ram on cheap monitoring VM. I found two solutions: old-school monitoring systems like `nagios`, `zabbix` or stats/metrics collectors like `prometheus`. While they are all great in many situations (I love zabbix), they are heavy, bulky and need deploying agents. I needed something simple, but flexible, ideally written in Go (self-containing, easily hackable). 
I found some interesting projects (some of them mentioned already), but at the end decided to write `zorix`, because I believe it is universal solution to many needs.

### Name ... it resembles zabbix ...

Yes, that is what I wanted. "Zori" is something like my trademark and at the past adding "x" at the end of program name makes application instantly cool (like nowadays "d"?), so `zorix` it be!
Reference to `zabbix` is intentional despite the fact, that `zorix` is incomparably simpler, less powerful and has completely different design.

