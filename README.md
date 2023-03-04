### Must have features
- [x] Calcualte Width in Mghz (20/40/80/160).
- [x] Switch column view for RSSI/Quality/Bars with sorting support.
- [x] Coloring of rows according to RSSI/Quality/Bars value (custom rows render).
- [x] Mark a row with color or emoji, if row relates to a network associated before wfmon start.
- [x] Generate code from manuf dictionary. Add Manufacture data.
- [x] Switch column view for BSSID/Manufacture with sorting support.
- [ ] Add flags support. -i --interface. if not provided then use default wifi interface.
- [ ] Search network by SSID, hotkey /

### Next Minor features
- [ ] Add Seen data/column.
- [ ] Add b/g/n/ac data.
- [ ] Add Rate data.
- [ ] Add Vendor data.
- [ ] ?Select default wifi interface using CoreWLAN api.
- [ ] ?Deassociate interface from network using CoreWLAN api.
- [ ] ?Change Radio channels hopping using CoreWLAN api.

### Next Majors features
- [ ] Support average sampling (RSSI, Noise values).
- [ ] Add an option to start program with analyze of given pcap file and end execution. -f --file.
- [ ] Sniffering with timeout, print simple output and end program execution. -p --print columns; -t --timeout.
- [ ] ?Add packets received line above the table.
- [ ] ?Verbose flag to print logs below the table. -v
- [ ] Add Info (with more data) widget of selected network on right table side.
- [ ] Add sectrum chart right below table.
- [ ] Add RSSI/Quality/Bars line chart right below table.
- [ ] Windows support
- [ ] Linux support
