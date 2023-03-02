package wifitable

import (
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/widgets/wifitable/sort"
)

// Sort by BSSID asc.
func ByBSSIDSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) string { return n[i].BSSID })
}

// Sort by station manufacture short name asc.
func ByManufSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) string { return n[i].Manuf })
}

// Sort by station manufacture full name asc.
func ByManufLongSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) string { return n[i].ManufLong })
}

// Sort by station manufacture full name asc.
func BySSIDSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) string { return n[i].NetworkName })
}

// Sort by Channel asc.
func ByChannelSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) int { return int(n[i].Channel) })
}

// Sort by Channel Width asc.
func ByChannelWidthSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) int { return int(n[i].ChannelWidth) })
}

// Sort by Bandwidth asc.
func ByBandwidthSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) int { return int(n[i].Band) })
}

// Sort by RSSI asc.
func ByRSSISorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) int { return int(n[i].RSSI) })
}

// Sort by Quality asc.
func ByQualitySorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) int { return int(n[i].Quality) })
}

// Sort by Quality/Bars asc.
func ByBarsSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) int { return int(n[i].Quality) })
}

// Sort by Noise asc.
func ByNoiseSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) int { return int(n[i].Noise) })
}

// Sort by SNR asc.
func BySNRSorter() sort.FncSorter {
	return sort.Sorter(func(n netdata.Slice, i int) int { return int(n[i].SNR) })
}
