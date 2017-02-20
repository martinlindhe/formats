package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// LZMA parse the lzma format
func LZMA(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isLZMA(c.Header) {
		return nil, nil
	}
	return parseLZMA(c.File, c.ParsedLayout)
}

func isLZMA(b []byte) bool {

	// XXX not proper magic , need other check
	if b[0] != 0x5d || b[1] != 0x00 {
		return false
	}
	return true
}

func parseLZMA(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 13, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			// XXX unsure of this stuff
			{Offset: pos, Length: 1, Info: "properties", Type: parse.Uint8},
			{Offset: pos + 1, Length: 4, Info: "dict cap", Type: parse.Uint32le},
			{Offset: pos + 5, Length: 8, Info: "uncompressed size", Type: parse.Uint64le},
		}}}

	return &pl, nil
}
