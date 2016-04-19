package image

// STATUS: 10%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	pcxPaletteType = map[uint16]string{
		1: "color",
		2: "grayscale",
	}
	pcxVersions = map[uint8]string{
		0: "2.5",
		2: "2.8 w/ palette",
		3: "2.8 w/out palette",
		5: "3.0 or better",
	}
)

func PCX(c *parse.ParseChecker)(*parse.ParsedLayout, error) {

	if !isPCX(&c.Header) {
		return nil, nil
	}
	return parsePCX(c.File, c.ParsedLayout)
}

func isPCX(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 0xa {
		return false
	}
	if b[1] != 0 && b[1] != 2 && b[1] != 3 && b[1] != 5 {
		return false
	}
	return true
}

func parsePCX(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	version, _ := parse.ReadUint8(file, pos+1)
	versionName := "?"
	if val, ok := pcxVersions[version]; ok {
		versionName = val
	}

	paletteType, _ := parse.ReadUint16le(file, pos+68)
	paletteTypeName := "?"
	if val, ok := pcxPaletteType[paletteType]; ok {
		paletteTypeName = val
	}

	pl.FileKind = parse.Image
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 128, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 1, Info: "magic", Type: parse.Uint8},
			{Offset: pos + 1, Length: 1, Info: "version = " + versionName, Type: parse.Uint8},
			{Offset: pos + 2, Length: 1, Info: "encoding", Type: parse.Uint8},
			{Offset: pos + 3, Length: 1, Info: "bits per plane", Type: parse.Uint8},
			{Offset: pos + 4, Length: 2, Info: "x min", Type: parse.Uint16le},
			{Offset: pos + 6, Length: 2, Info: "y min", Type: parse.Uint16le},
			{Offset: pos + 8, Length: 2, Info: "x max", Type: parse.Uint16le},
			{Offset: pos + 10, Length: 2, Info: "y max", Type: parse.Uint16le},
			{Offset: pos + 12, Length: 2, Info: "vertical dpi", Type: parse.Uint16le},
			{Offset: pos + 14, Length: 2, Info: "horizontal dpi", Type: parse.Uint16le},
			{Offset: pos + 16, Length: 48, Info: "palette", Type: parse.Bytes},
			{Offset: pos + 64, Length: 1, Info: "reserved", Type: parse.Uint8},
			{Offset: pos + 65, Length: 1, Info: "color planes", Type: parse.Uint8},
			{Offset: pos + 66, Length: 2, Info: "bytes per plane line", Type: parse.Uint16le},
			{Offset: pos + 68, Length: 2, Info: "palette type = " + paletteTypeName, Type: parse.Uint16le},
			{Offset: pos + 70, Length: 2, Info: "hScrSize", Type: parse.Uint16le},
			{Offset: pos + 72, Length: 2, Info: "vScrSize", Type: parse.Uint16le},
			{Offset: pos + 74, Length: 54, Info: "padding", Type: parse.Bytes}, // XXX may be 56 byte if horiz dpi is absent
		}}, {
		Offset: pos + 128,
		Length: pl.FileSize - 128,
		Info:   "image data",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos + 128, Length: pl.FileSize - 128, Info: "image data", Type: parse.Bytes},
		}}}

	return &pl, nil
}
