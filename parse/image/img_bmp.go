package image

// TODO parse & display CIEXYZTRIPLE endpoint data: FXPT2DOT30  X, Y, Z
// TODO samples using png / jpeg compression , and properly decode/extract to file that part as a sub-resource or sth....

// STATUS: 90%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	bmpCompressions = map[uint32]string{
		0: "rgb",
		1: "rle8",
		2: "rle4",
		3: "bitfields",
		4: "jpeg",
		5: "png",
	}
	v2len = 64
	v3len = 40
	v4len = 68
)

// BMP parses the bmp format
func BMP(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isBMP(c.Header) {
		return nil, nil
	}
	return parseBMP(c.File, c.ParsedLayout)
}

func isBMP(b []byte) bool {

	return b[0] == 'B' && b[1] == 'M'
}

func parseBMP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	fileHeaderLen := int64(14)
	pl.FileKind = parse.Image
	pl.MimeType = "image/x-ms-bmp"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: fileHeaderLen,
		Info:   "file header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 2, Length: 4, Info: "file size", Type: parse.Uint32le},
			{Offset: pos + 6, Length: 4, Info: "reserved", Type: parse.Uint32le},
			{Offset: pos + 10, Length: 4, Info: "image data offset", Type: parse.Uint32le},
		}}}

	infoHeader, err := parseBMPInfoHeader(file)
	if err != nil {
		return nil, err
	}

	pl.Layout = append(pl.Layout, infoHeader)

	// body
	headerLen := int64(fileHeaderLen + infoHeader.Length)

	dataLayout := parse.Layout{
		Offset: headerLen,
		Info:   "image data",
		Length: pl.FileSize - headerLen,
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: headerLen, Length: pl.FileSize - headerLen, Info: "image data", Type: parse.Bytes},
		}}

	pl.Layout = append(pl.Layout, dataLayout)

	return &pl, nil
}

func parseBMPInfoHeader(file *os.File) (parse.Layout, error) {

	pos := int64(14)
	layout := parse.Layout{
		Offset: pos,
		Type:   parse.Group}

	infoHdrSize, _ := parse.ReadUint32le(file, pos)
	layout.Length = int64(infoHdrSize)

	switch infoHdrSize {
	case 12: // OS/2 V1 - BITMAPCOREHEADER
		layout.Info = "info header V1"
		layout.Childs = parseBMPVersion1Header(file, pos)

	case 64: // OS/2 V2
		layout.Info = "info header V2"
		v3 := parseBMPVersion3Header(file, pos)
		v2 := parseBMPVersion2Header(file, pos+int64(v3len))
		layout.Childs = append(v3, v2...)

	case 40: // Windows V3 - BITMAPINFOHEADER
		layout.Info = "info header V3"
		layout.Childs = parseBMPVersion3Header(file, pos)

	case 108: // Windows V4 - BITMAPV4HEADER
		layout.Info = "info header V4"
		v3 := parseBMPVersion3Header(file, pos)
		v4 := parseBMPVersion4Header(file, pos+int64(v3len))
		layout.Childs = append(v3, v4...)

	case 124: // Windows V5 - BITMAPV5HEADER
		layout.Info = "info header V5"
		v3 := parseBMPVersion3Header(file, pos)
		v4 := parseBMPVersion4Header(file, pos+int64(v3len))
		v5 := parseBMPVersion5Header(file, pos+int64(v3len)+int64(v4len))
		layout.Childs = append(v3, v4...)
		layout.Childs = append(layout.Childs, v5...)

	default:
		return layout, fmt.Errorf("unrecognized header size %d", infoHdrSize)
	}

	return layout, nil
}

func parseBMPVersion1Header(file *os.File, pos int64) []parse.Layout {

	return []parse.Layout{
		{Offset: pos, Length: 4, Type: parse.Uint32le, Info: "info header size"},
		{Offset: pos + 4, Length: 2, Type: parse.Uint16le, Info: "width"},
		{Offset: pos + 6, Length: 2, Type: parse.Uint16le, Info: "height"},
		{Offset: pos + 8, Length: 2, Type: parse.Uint16le, Info: "planes"},
		{Offset: pos + 10, Length: 2, Type: parse.Uint16le, Info: "bpp"},
	}
}

func parseBMPVersion2Header(file *os.File, pos int64) []parse.Layout {

	return []parse.Layout{
		{Offset: pos, Length: 2, Type: parse.Uint16le, Info: "units"},
		{Offset: pos + 2, Length: 2, Type: parse.Uint16le, Info: "reserved"},
		{Offset: pos + 4, Length: 2, Type: parse.Uint16le, Info: "recording"},
		{Offset: pos + 6, Length: 2, Type: parse.Uint16le, Info: "rendering"},
		{Offset: pos + 8, Length: 4, Type: parse.Uint32le, Info: "size1"},
		{Offset: pos + 12, Length: 4, Type: parse.Uint32le, Info: "size2"},
		{Offset: pos + 16, Length: 4, Type: parse.Uint32le, Info: "color encoding"},
		{Offset: pos + 20, Length: 4, Type: parse.Uint32le, Info: "identifier"},
	}
}

func parseBMPVersion3Header(file *os.File, pos int64) []parse.Layout {

	compressionName, _ := parse.ReadToMap(file, parse.Uint32le, pos+16, bmpCompressions)
	return []parse.Layout{
		{Offset: pos, Length: 4, Type: parse.Uint32le, Info: "info header size"},
		{Offset: pos + 4, Length: 4, Type: parse.Uint32le, Info: "width"},
		{Offset: pos + 8, Length: 4, Type: parse.Uint32le, Info: "height"},
		{Offset: pos + 12, Length: 2, Type: parse.Uint16le, Info: "planes"},
		{Offset: pos + 14, Length: 2, Type: parse.Uint16le, Info: "bpp"},
		{Offset: pos + 16, Length: 4, Type: parse.Uint32le, Info: "compression = " + compressionName},
		{Offset: pos + 20, Length: 4, Type: parse.Uint32le, Info: "size of picture"},
		{Offset: pos + 24, Length: 4, Type: parse.Uint32le, Info: "horizontal resolution"},
		{Offset: pos + 28, Length: 4, Type: parse.Uint32le, Info: "vertical resolution"},
		{Offset: pos + 32, Length: 4, Type: parse.Uint32le, Info: "number of used colors"},
		{Offset: pos + 36, Length: 4, Type: parse.Uint32le, Info: "number of important colors"},
	}
}

func parseBMPVersion4Header(file *os.File, pos int64) []parse.Layout {

	return []parse.Layout{
		{Offset: pos, Length: 4, Type: parse.Uint32le, Info: "red mask"},
		{Offset: pos + 4, Length: 4, Type: parse.Uint32le, Info: "green mask"},
		{Offset: pos + 8, Length: 4, Type: parse.Uint32le, Info: "blue mask"},
		{Offset: pos + 12, Length: 4, Type: parse.Uint32le, Info: "alpha mask"},

		// NOTE: v5 file use "BGRs", while v4 use 0x1
		{Offset: pos + 16, Length: 4, Type: parse.Uint32le, Info: "cs type"},
		{Offset: pos + 20, Length: 12, Type: parse.Uint8, Info: "ciexyz red"},
		{Offset: pos + 32, Length: 12, Type: parse.Uint8, Info: "ciexyz green"},
		{Offset: pos + 44, Length: 12, Type: parse.Uint8, Info: "ciexyz blue"},
		{Offset: pos + 56, Length: 4, Type: parse.Uint32le, Info: "gamma red"},
		{Offset: pos + 60, Length: 4, Type: parse.Uint32le, Info: "gamma green"},
		{Offset: pos + 64, Length: 4, Type: parse.Uint32le, Info: "gamma blue"},
	}
}

func parseBMPVersion5Header(file *os.File, pos int64) []parse.Layout {

	return []parse.Layout{
		{Offset: pos, Length: 4, Type: parse.Uint32le, Info: "intent"},
		{Offset: pos + 4, Length: 4, Type: parse.Uint32le, Info: "profile data"},
		{Offset: pos + 8, Length: 4, Type: parse.Uint32le, Info: "profile size"},
		{Offset: pos + 12, Length: 4, Type: parse.Uint32le, Info: "reserved"},
	}
}
