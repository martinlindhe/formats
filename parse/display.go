package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
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

func (pl *ParsedLayout) PrettyASCIIView(file *os.File, hexView HexViewState) string {

	ascii := ""
	base := hexView.StartingRow * int64(hexView.RowWidth)
	ceil := base + int64(hexView.VisibleRows*hexView.RowWidth)

	for i := base; i < ceil; i += int64(hexView.RowWidth) {

		ofs, err := file.Seek(i, os.SEEK_SET)
		if i != ofs {
			log.Fatalf("err: unexpected offset %04x, expected %04x\n", ofs, i)
		}
		line, err := pl.GetASCII(file, hexView)

		ascii += line + "\n"
		if err != nil {
			fmt.Println("got err", err)
			break
		}
	}
	return ascii
}

// PrettyOffsetView ...
func (pl *ParsedLayout) PrettyOffsetView(file *os.File, hexView HexViewState) string {

	ofsFmt := "%08x"
	padding := 0
	if pl.FileSize <= 0xffff {
		ofsFmt = "%04x"
		padding = 3
	} else if pl.FileSize <= 0xffffff {
		ofsFmt = "%06x"
		padding = 1
	}

	base := hexView.StartingRow * int64(hexView.RowWidth)
	ceil := base + int64(hexView.VisibleRows*hexView.RowWidth)

	pad := strings.Repeat(" ", padding)

	res := ""
	for i := base; i < ceil; i += int64(hexView.RowWidth) {
		res += pad + fmt.Sprintf(ofsFmt, i) + "\n"
	}
	return res
}

// PrettyHexView ...
func (pl *ParsedLayout) PrettyHexView(file *os.File, hexView HexViewState) string {

	hex := ""
	base := hexView.StartingRow * int64(hexView.RowWidth)
	ceil := base + int64(hexView.VisibleRows*hexView.RowWidth)

	for i := base; i < ceil; i += int64(hexView.RowWidth) {

		ofs, err := file.Seek(i, os.SEEK_SET)
		if i != ofs {
			log.Fatalf("err: unexpected offset %04x, expected %04x\n", ofs, i)
		}
		line, err := pl.GetHex(file, hexView)
		hex += line + "\n"
		if err != nil {
			break
		}
	}
	return hex
}

func (pl *ParsedLayout) GetASCII(file *os.File, hexView HexViewState) (string, error) {

	if len(pl.Layout) == 0 {
		return "", fmt.Errorf("pl.Layout is empty")
	}

	layout := pl.Layout[hexView.CurrentField]

	reader := io.Reader(file)

	symbols := []string{}

	base, err := file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return "", err
	}

	formatting := HexFormatting{
		BetweenSymbols: "",
		GroupSize:      1,
	}

	for w := int64(0); w < 16; w++ {
		var b byte
		if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
			if err == io.EOF {
				return combineHexRow(symbols, formatting), nil
			}
			return "", err
		}

		ceil := base + w

		colorName := "fg-white"
		if !pl.isOffsetKnown(base + w) {
			colorName = "fg-red"
		}
		if ceil >= layout.Offset && ceil < layout.Offset+int64(layout.Length) {
			colorName = "fg-cyan"
		}
		// XXX byte as string
		tok := "."
		if b > 31 && b < 128 && b != '[' && b != ']' {
			// [] is used by termui, so we cant display them + colors :P
			tok = fmt.Sprintf("%c", b)
		}
		group := fmt.Sprintf("[%s](%s)", tok, colorName)
		symbols = append(symbols, group)
	}

	return combineHexRow(symbols, formatting), nil
}

// GetHex dumps a row of hex from io.Reader
func (pl *ParsedLayout) GetHex(file *os.File, hexView HexViewState) (string, error) {

	formatting := HexFormatting{
		BetweenSymbols: " ",
		GroupSize:      1,
	}

	if len(pl.Layout) == 0 {
		return "", fmt.Errorf("pl.Layout is empty")
	}

	layout := pl.Layout[hexView.CurrentField]

	reader := io.Reader(file)

	symbols := []string{}

	base, err := file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return "", err
	}

	for w := int64(0); w < 16; w++ {
		var b byte
		if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
			if err == io.EOF {
				return combineHexRow(symbols, formatting), nil
			}
			return "", err
		}

		ceil := base + w

		colorName := "fg-white"
		if !pl.isOffsetKnown(base + w) {
			colorName = "fg-red"
		}
		if ceil >= layout.Offset && ceil < layout.Offset+int64(layout.Length) {
			colorName = "fg-cyan"
		}

		group := fmt.Sprintf("[%02x](%s)", b, colorName)
		symbols = append(symbols, group)
	}

	return combineHexRow(symbols, formatting), nil
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
