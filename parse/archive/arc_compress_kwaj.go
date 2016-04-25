package archive

// MS-DOS COMPRESS.EXE
// Extensions: .??_
// http://www.cabextract.org.uk/libmspack/doc/szdd_kwaj_format.html

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	kwajCompressions = map[uint16]string{
		0: "none",
		1: "xor 255",
		2: "regular SZDD",
		3: "LZ + Huffman 'Jeff Johnson'",
		4: "MS-ZIP",
	}
)

func CompressKWAJ(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isCompressKWAJ(c.Header) {
		return nil, nil
	}
	return parseCompressKWAJ(c.File, c.ParsedLayout)
}

func isCompressKWAJ(b []byte) bool {

	if b[0] != 'K' || b[1] != 'W' || b[2] != 'A' || b[3] != 'J' ||
		b[4] != 0x88 || b[5] != 0xf0 || b[6] != 0x27 || b[7] != 0xd1 {
		return false
	}
	return true
}

func parseCompressKWAJ(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	compressionName, _ := parse.ReadToMap(file, parse.Uint16le, pos+8, kwajCompressions)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 14, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 8, Info: "magic", Type: parse.Bytes},
			{Offset: pos + 8, Length: 2, Info: "compression = " + compressionName, Type: parse.Uint16le},
			{Offset: pos + 10, Length: 2, Info: "data offset", Type: parse.Uint16le},
			{Offset: pos + 12, Length: 2, Info: "flags", Type: parse.Uint16le}, // XXX bit mask
		}}}

	return &pl, nil
}
