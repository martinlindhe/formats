package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// TD2 parses the td2 format
func TD2(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isTD2(c.Header) {
		return nil, nil
	}
	return parseTD2(c.File, c.ParsedLayout)
}

func isTD2(b []byte) bool {

	if b[0] != 't' || b[1] != 'd' || b[2] != 0 {
		return false
	}
	return true
}

func parseTD2(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 3,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 3, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
