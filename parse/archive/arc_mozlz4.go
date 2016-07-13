package archive

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// MozillaLZ4 parses the Mozilla lz4 format
func MozillaLZ4(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isMozillaLZ4(c.Header) {
		return nil, nil
	}
	return parseMozillaLZ4(c.File, c.ParsedLayout)
}

func isMozillaLZ4(b []byte) bool {

	// XXX not proper magic , need other check
	if b[0] != 'm' || b[1] != 'o' || b[2] != 'z' ||
		b[3] != 'L' || b[4] != 'z' || b[5] != '4' {
		return false
	}
	return true
}

func parseMozillaLZ4(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 6, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 6, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
