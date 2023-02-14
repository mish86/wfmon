package mac

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const MAC48 = 48
const base10 = 10
const base16 = 16
const bitsInNibble = 4

// const bitsInOctet = 8

// Subnet Hardware Address presentation.
// Supports IEEE 802 MAC-48 only.
// TODO write tests.
type HardwareAddr struct {
	addr   *big.Int
	prefix uint8
}

// Accepts hardware address in IEEE 802 MAC-48 format with subnet prefix.
// E.g.
// 00:03:93
// 00:55:DA:00:00:00/28
// 40:D8:55:00:20:00/36
// Returns new hardware address.
func (a *HardwareAddr) WithAddr(mac string) *HardwareAddr {
	macAndMask := strings.Split(mac, "/")

	// mac in hex
	h := hex(macAndMask[0])
	// MAC-48 supported only
	if len(h) > MAC48/bitsInNibble {
		return a
	}

	// prefix
	var p uint8
	{
		if len(macAndMask) > 1 {
			if pI, err := strconv.Atoi(prefix(macAndMask[1])); err != nil {
				p = uint8(pI)
			}
		}
		if p == 0 {
			p = uint8(len(h) * bitsInNibble)
		}
	}

	if addrI, ok := new(big.Int).SetString(h, base16); ok {
		a = &HardwareAddr{
			addr:   addrI,
			prefix: p,
		}

		a = a.align()
	}

	return a
}

// Returns a hardware address with applied prefix.
func (a *HardwareAddr) WithPrefix(mask string) *HardwareAddr {
	if len(mask) == 0 {
		return a
	}

	// prefix
	p := a.prefix
	if pI, err := strconv.Atoi(prefix(mask)); err != nil {
		p = uint8(pI)
	}

	addr := a.addr
	if addr == nil {
		addr = new(big.Int)
	}

	a = &HardwareAddr{
		addr:   new(big.Int).Set(addr),
		prefix: p,
	}

	return a.align()
}

// Alings hardware address with prefix.
func (a *HardwareAddr) align() *HardwareAddr {
	if a.addr == nil || a.prefix == 0 {
		return a
	}

	got := a.addr.BitLen() / bitsInNibble
	expected := a.prefix / bitsInNibble
	shift := got - int(expected)

	addr := a.addr
	if shift > 0 {
		addr = new(big.Int).Rsh(a.addr, uint(shift))
	}

	return &HardwareAddr{
		addr:   addr,
		prefix: a.prefix,
	}
}

// Returns parent hardware address with subnet which includes current.
func (a *HardwareAddr) Parent() *HardwareAddr {
	if a.addr == nil || a.prefix == 0 {
		return nil
	}

	addr := new(big.Int).Rsh(a.addr, bitsInNibble)
	p := a.prefix - bitsInNibble

	// no align required
	return &HardwareAddr{
		addr:   addr,
		prefix: p,
	}
}

// Removes octets delimiters.
func hex(addr string) string {
	hex := addr
	hex = strings.ReplaceAll(hex, ":", "")
	hex = strings.ReplaceAll(hex, "-", "")
	hex = strings.ReplaceAll(hex, ".", "")
	return hex
}

// Removes prefix delimiter.
func prefix(mask string) string {
	return strings.Replace(mask, "/", "", 1)
}

// Converts hardware address into string presentation.
func (a *HardwareAddr) String() string {
	res := strings.Builder{}
	for i, b := range a.addr.Bytes() {
		res.WriteString(fmt.Sprintf("%x", b))
		if i%2 == 0 {
			res.WriteString(":")
		}
	}

	res.WriteString(fmt.Sprintf("/%d", a.prefix))

	return res.String()
}

// Retruns subnet prefix bits capacity.
func wildcard(prefix uint8) uint8 {
	return MAC48 - prefix
}

// Prefix dot hardware address presentation.
// Useful as key in map.
type WildcardDotBigInt HardwareAddr

func (a WildcardDotBigInt) String() string {
	return fmt.Sprintf("%d.%s", wildcard(a.prefix), a.addr.Text(base10))
}
