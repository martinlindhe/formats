package formats

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
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

// ...
var (
	HexView = HexViewState{
		StartingRow:  0,
		VisibleRows:  11,
		RowWidth:     16,
		CurrentField: 0,
	}
)

// CurrentFieldInfo renders info of current field
func (f *HexViewState) CurrentFieldInfo(file *os.File, pl ParsedLayout) string {

	field := pl.Layout[f.CurrentField]

	res := "field: " + field.Info

	if field.Length > 4 {
		res += fmt.Sprintf(" (%d bytes)", field.Length)
	}

	res += "\nvalue: "

	file.Seek(field.Offset, os.SEEK_SET)

	// decode data based on type and show
	r := io.Reader(file)

	switch field.Type {
	case Int8:
		var i int8
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint8:
		var i uint8
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Int16le:
		var i int16
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint16le:
		var i uint16
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Int32le:
		var i int32
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case Uint32le:
		var i uint32
		if err := binary.Read(r, binary.LittleEndian, &i); err != nil {
			panic(err)
		}
		res += fmt.Sprintf("%d", i)

	case ASCII, ASCIIZ:
		buf := make([]byte, field.Length)
		_, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		res += string(buf)

	default:
		res += "XXX unhandled " + field.Type.String()
	}

	res += " (" + field.Type.String() + ")"

	return res
}

// Next moves focus to the next field
func (f *HexViewState) Next(max int) {
	f.CurrentField++
	if f.CurrentField >= max {
		f.CurrentField = max - 1
	}
}

// Prev moves focus to the previous field
func (f *HexViewState) Prev() {
	f.CurrentField--
	if f.CurrentField < 0 {
		f.CurrentField = 0
	}
}

func combineHexRow(symbols []string, formatting HexFormatting) string {

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
