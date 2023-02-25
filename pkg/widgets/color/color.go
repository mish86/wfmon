package color

import (
	"fmt"
	"math/rand"
	"time"

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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var next HexIter
	next = func() (HexColor, HexIter) {
		return HexColor{
			byte(r.Intn(255)),
			byte(r.Intn(255)),
			byte(r.Intn(255)),
		}, next
	}

	return next
}
