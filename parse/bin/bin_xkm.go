package bin

// XXX need a big-endian sample file

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// XKM parses the Compiled XKB Keymap format
func XKM(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isXKM(c.Header) {
		return nil, nil
	}
	return parseXKM(c.File, c.ParsedLayout)
}

func isXKM(b []byte) bool {

	if b[1] == 'm' && b[2] == 'k' && b[3] == 'x' {
		return true // le
	}
	if b[0] == 'x' && b[1] == 'k' && b[2] == 'k' {
		return true // be
	}
	return false
}

func parseXKM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			// XXX this is le-format only
			{Offset: pos, Length: 1, Info: "version", Type: parse.Bytes},
			{Offset: pos + 1, Length: 3, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
