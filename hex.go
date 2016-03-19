package formats

import (
	"strings"
)

// HexFormatting ...
type HexFormatting struct {
	BetweenSymbols string
	GroupSize      byte
}

// HexViewState ...
type HexViewState struct {
	StartingRow  int64
	VisibleRows  int
	RowWidth     int
	CurrentField int
}

var (
	formatting = HexFormatting{
		BetweenSymbols: " ",
		GroupSize:      1,
	}
	hexViewState = HexViewState{
		StartingRow:  0,
		VisibleRows:  10,
		RowWidth:     16,
		CurrentField: 0,
	}
)

// Formatting ...
func Formatting(fmt HexFormatting) { formatting = fmt }

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
