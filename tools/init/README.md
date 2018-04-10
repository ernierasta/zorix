Zorix init scripts
==================

In this directory you can find init scripts for different inits.

Installation notes (for install scripts):

```bash
sudo cp zorix /usr/local/bin/zorix
sudo mkdir /etc/zorix
sudo cp config.sample.toml /etc/zorix/config.toml
sudo useradd zorix -s /sbin/nologin -M
```

Paths:  
binary(zorix): `/usr/local/bin/zorix`  
config: `/etc/zorix/config.toml`

runit
-----

Tested on Void Linux.

```bash
sudo cp tools/init/runit/zorix /etc/sv/
sudo chmod +x /etc/sv/zorix/run /etc/sv/zorix/log/run
sudo mkdir /var/log/zorix
```

sysvinit
--------

Tested on strange debian 9.4, but with sysvinit. ;-)
Logging directly to /var/log/zorix.log

```bash
cp tools/init/sysvinit/debian/zorix /etc/init.d/
chmod +x /etc/init.d/zorix
update-rc.d zorix defaults
service zorix start
```


