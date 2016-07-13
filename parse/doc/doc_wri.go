package doc

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// WRI parses the Windows wri document format
func WRI(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isWRI(c.Header) {
		return nil, nil
	}
	return parseWRI(c.File, c.ParsedLayout)
}

func isWRI(b []byte) bool {

	// TODO what is right magic bytes? just guessing
	// FIXME IT IS     if data.find(b'\xBE\x00\x00\x00\xAB\x00\x00\x00\x00\x00\x00\x00\x00') == 1
	if b[0] != 0x31 || b[1] != 0xbe || b[2] != 0 || b[3] != 0 {
		return false
	}
	return true
}

func parseWRI(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Document
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
