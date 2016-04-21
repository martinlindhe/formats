package bin

// Compiled terminfo entry
// /usr/share/file/magic/terminfo

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func Terminfo(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isTerminfo(c) {
		return nil, nil
	}
	return parseTerminfo(c.File, c.ParsedLayout)
}

func isTerminfo(c *parse.ParseChecker) bool {

	b := c.Header
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
