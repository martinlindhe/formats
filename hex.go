package formats

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// HexFormatting ...
type HexFormatting struct {
	BetweenSymbols string
	GroupSize      byte
}

var (
	formatting = HexFormatting{
		BetweenSymbols: " ",
		GroupSize:      1, // XXX TODO other group sizes
	}
)

// Formatting ...
func Formatting(fmt HexFormatting) {
	formatting = fmt
}

// GetHex ...
func GetHex(r *io.Reader, height int) (res []string, err error) {
	// XXX dump one screen of hex (300 byte or so, ) from io.Reader

	for h := 0; h < height; h++ {
		symbols := []string{}

		for w := 0; w < 16; w++ {
			var b byte
			if err = binary.Read(*r, binary.LittleEndian, &b); err != nil {
				res = append(res, combineHexRow(symbols))
				return
			}
			group := fmt.Sprintf("%02x", b)
			symbols = append(symbols, group)
		}
		res = append(res, combineHexRow(symbols))
	}
	return
}

func combineHexRow(symbols []string) string {

	group := []string{}
	row := []string{}
	cur := byte(0)

	for _, sym := range symbols {
		cur++
		group = append(group, sym)
		if cur == formatting.GroupSize {
			row = append(row, strings.Join(group, ""))
			group = nil
			cur = 0
		}
	}
	return strings.Join(row, formatting.BetweenSymbols)
}
