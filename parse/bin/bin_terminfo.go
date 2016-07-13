package bin

// /usr/share/file/magic/terminfo

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// Terminfo parses the Compiled terminfo entry format
func Terminfo(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isTerminfo(c.Header) {
		return nil, nil
	}
	return parseTerminfo(c.File, c.ParsedLayout)
}

func isTerminfo(b []byte) bool {

	if b[0] != 0x1a || b[1] != 1 {
		return false
	}
	return true
}

func parseTerminfo(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
