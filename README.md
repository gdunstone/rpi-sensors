## rpi-sensors

This project is an attempt to modularise deployment of sensor monitoring software onto raspberry pis.

It is primarily an ansible playbook to facilitate the distribution onto raspberry pis with [Arch Linux ARM](https://archlinuxarm.org/) installed on them.

There is also golang programs for supported sensors that output data in the influx line protocol to be run by telegraf.
