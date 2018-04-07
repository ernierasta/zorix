Zorix init scripts
==================

In this directory you can find init scripts for different inits.

Installation notes (for install scripts):

sudo cp zorix /usr/local/bin/zorix
sudo mkdir /etc/zorix
sudo cp config.sample.toml /etc/zorix/config.toml
sudo useradd zorix -s /sbin/nologin -M

Paths:
binary(zorix): /usr/local/bin/zorix
config: /etc/zorix/config.toml

runit
-----

Tested on Void Linux.

`sudo cp tools/init/runit/zorix /etc/sv/`
`sudo chmod +x /etc/sv/zorix/run /etc/sv/zorix/log/run`
`sudo mkdir /var/log/zorix`
