package bin

// ??? created by libre office

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func RDB(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isRDB(c) {
		return nil, nil
	}
	return parseRDB(c.File, c.ParsedLayout)
}

func isRDB(c *parse.ParseChecker) bool {

	b := c.Header
	if b[0] != 'U' || b[1] != 'N' || b[2] != 'O' || b[3] != 'I' {
		return false
	}
	return true
}

func parseRDB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
