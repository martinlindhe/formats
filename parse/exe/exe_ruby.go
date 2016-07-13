package exe

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// Ruby parses the Ruby bytecode format
func Ruby(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isRuby(c.Header) {
		return nil, nil
	}
	return parseRuby(c.File, c.ParsedLayout)
}

func isRuby(b []byte) bool {

	if b[0] == 4 && b[1] == 8 && b[2] == 0x55 && b[3] == 0x3a {
		return true
	}
	return false
}

func parseRuby(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
