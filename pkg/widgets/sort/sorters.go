package sort

import (
	netdata "wfmon/pkg/data/net"
)

// Sort by BSSID asc.
func ByBSSIDSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) string { return n[i].BSSID })
}

// Sort by station manufacture short name asc.
func ByManufSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) string { return n[i].Manuf })
}

// Sort by station manufacture full name asc.
func ByManufLongSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) string { return n[i].ManufLong })
}

// Sort by station manufacture full name asc.
func BySSIDSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) string { return n[i].NetworkName })
}

// Sort by Channel asc.
func ByChannelSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) int { return int(n[i].Channel) })
}

// Sort by Channel Width asc.
func ByChannelWidthSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) int { return int(n[i].ChannelWidth) })
}

// Sort by Bandwidth asc.
func ByBandwidthSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) int { return int(n[i].Band) })
}

// Sort by RSSI asc.
func ByRSSISorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) int { return int(n[i].RSSI) })
}

// Sort by Quality asc.
func ByQualitySorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) int { return int(n[i].Quality) })
}

// Sort by Quality/Bars asc.
func ByBarsSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) int { return int(n[i].Quality) })
}

// Sort by Noise asc.
func ByNoiseSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) int { return int(n[i].Noise) })
}

// Sort by SNR asc.
func BySNRSorter() FncSorter {
	return Sorter(func(n netdata.Slice, i int) int { return int(n[i].SNR) })
}
