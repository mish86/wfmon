package color

import (
	"fmt"

	"math/rand"
	"time"
	"wfmon/pkg/utils/cmp"

	"github.com/charmbracelet/lipgloss"
)

type HexColor [3]byte

func Black() HexColor {
	return HexColor{}
}

func (hex HexColor) String() string {
	return fmt.Sprintf("#%02x%02x%02x", hex[0], hex[1], hex[2])
}

func (hex HexColor) Lipgloss() lipgloss.Color {
	return lipgloss.Color(hex.String())
}

type HexIter func() (HexColor, HexIter)

func Random() HexIter {
	const (
		byteSize = 255
		darker   = 220
	)

	var dark = func(v byte) byte {
		return cmp.Nvl(v <= darker, v, byteSize-v)
	}
	var light = func(v byte) byte {
		return cmp.Nvl(v > darker, v, byteSize-v)
	}

	adaptive := dark
	if lipgloss.HasDarkBackground() {
		adaptive = light
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var gen = func() byte {
		return adaptive(byte(r.Intn(byteSize)))
	}

	var next HexIter
	next = func() (HexColor, HexIter) {
		return HexColor{gen(), gen(), gen()}, next
	}

	return next
}
