### Demo
[![asciicast](https://asciinema.org/a/7WTEroHbU9YMrMYDXKW37ZxtX.svg)](https://asciinema.org/a/7WTEroHbU9YMrMYDXKW37ZxtX)

### TODO
- [x] Calcualte Width in Mghz (20/40/80/160).
- [x] Switch column view for RSSI/Quality/Bars with sorting support.
- [x] Coloring of rows according to RSSI/Quality/Bars value (custom rows render).
- [x] Mark a row with color or emoji, if row relates to a network associated before wfmon start.
- [x] Generate code from manuf dictionary. Add Manufacture data.
- [x] Switch column view for BSSID/Manufacture with sorting support.
- [x] Add Vendor data.
- [x] Add RSSI/Quality spectrum chart.
- [x] Add RSSI/Quality sparkline chart.
- [ ] Add flags support. -i --interface. if not provided then use default wifi interface.
- [ ] Add an option to start program with analyze of given pcap file and end execution. -f --file.
- [x] ?Determine default wifi interface using CoreWLAN api.
- [x] ?Deassociate interface from network before set on monitoring using CoreWLAN api.
- [x] ?Change radio channels during scan using CoreWLAN api.
- [ ] Support average sampling for RSSI and Noise values.
- [ ] Search network by SSID, hotkey /
- [ ] ?Verbose flag to print logs below the table and charts. -v
- [ ] ?Windows support
- [ ] ?Linux support
- [ ] ?Add packets received stats as a line above the table.
- [ ] ?Add Info (with more data) widget of highlighted network.
- [ ] ?Add Seen data/column.
- [ ] ?Add b/g/n/ac data.
- [ ] ?Add Rate data.
