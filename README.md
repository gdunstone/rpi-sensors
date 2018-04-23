## rpi-sensors

[![Build Status](https://travis-ci.org/gdunstone/rpi-sensors.svg?branch=master)](https://travis-ci.org/gdunstone/rpi-sensors)

This project is an attempt to modularise deployment of sensor monitoring software onto raspberry pis.

It is primarily an ansible playbook to facilitate the distribution onto raspberry pis with [Arch Linux ARM](https://archlinuxarm.org/) installed on them.

There is also golang programs for supported sensors that output data in the influx line protocol to be run by telegraf.


To perform an install (on an rpi 2) you need to follow these directions:
https://archlinuxarm.org/platforms/armv7/broadcom/raspberry-pi-2

then ssh into the device to install the required dependencies:
 $ ssh alarm@<rpi ip>
 alarm@<rpi ip>'s password: alarm
 [alarm@alarmpi]# su
 Password: root
 [root@alarmpi alarm] pacman -Sy python2 i2c-tools
 [root@alarmpi alarm] reboot

You will need to add your raspberry pi to the roles you want it to have in `hosts`

So if I had connected a DHT22 sensor on pin 14 to my raspberry pi, I would add it under the dht22 role and add pin=14
