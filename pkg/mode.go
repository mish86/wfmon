package app

import (
	"strings"
)

type Mode int

const (
	Dev Mode = iota
	Prod
)

func (m Mode) String() string {
	var modeStrings = []string{
		Dev:  "DEV",
		Prod: "PROD",
	}

	if m < 0 || int(m) >= len(modeStrings) {
		return Prod.String()
	}
	return modeStrings[m]
}

func FromString(s string) Mode {
	var modeMap = map[string]Mode{
		"DEV":  Dev,
		"PROD": Prod,
	}

	if m, ok := modeMap[strings.ToUpper(s)]; ok {
		return m
	}

	return Prod
}
