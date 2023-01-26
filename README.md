networksetup -listallhardwareports
    Hardware Port: Wi-Fi
    Device: en0
    Ethernet Address: 3c:22:fb:ef:dd:1d

networksetup -getmacaddress en0
	Ethernet Address: 3c:22:fb:ef:dd:1d (Device: en0)

networksetup -getairportnetwork en0
	Current Wi-Fi Network: WIFI_196

networksetup -getnetworkserviceenabled Wi-Fi
	Enabled

system_profiler SPAirPortDataType -json | jq . -- reports system hardware and software configuration
networksetup -setairportnetwork en0 WIFI_196 -- connect to ap
networksetup -setnetworkserviceenabled Wi-Fi off/on -- turn off/on wifi
sudo airport en0 -z -- disconnect from ap
networksetup -removepreferredwirelessnetwork en0 WIFI_196
	Removed WIFI_196 from the preferred networks list

arp -an
	? (192.168.1.1) at 4c:32:2d:10:80:70 on en0 ifscope [ethernet]
	? (192.168.1.255) at ff:ff:ff:ff:ff:ff on en0 ifscope [ethernet]
	? (224.0.0.251) at 1:0:5e:0:0:fb on en0 ifscope permanent [ethernet]