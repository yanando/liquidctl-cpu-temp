# liquidctl-cpu-temp
Script that monitors cpu temps and sets cpu cooler temps according to entered fan/pump curves.
Only tested on NZXT kraken z63

~~requires lm-sensors and liquidctl to be installed~~

vendorID and ProductID can be found by using the lsusb command

I.E: Bus 001 Device 004: ID 1e71:3008 NZXT NZXT KrakenZ Device

with 0x1e71 being the vendorID and 0x3008 being the productID

get the temperatureDevice from the sensors command, I.E k10temp-pci-00c3
temperaturePath is optional and should point to a file containing the temperature of the cpu, usually located somewhere in /sys/devices
I.E /sys/devices/pci0000:00/0000:00:18.3/hwmon/hwmon2/temp1_input
