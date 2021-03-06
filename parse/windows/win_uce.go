package windows

// STATUS: 1%
// found in win10-os/Windows/System32/bopomofo.uce

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// UCE parses the uce format
func UCE(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isUCE(c.Header) {
		return nil, nil
	}
	return parseUCE(c.File, c.ParsedLayout)
}

func isUCE(b []byte) bool {

	if b[0] != 'U' || b[1] != 'C' || b[2] != 'E' || b[3] != 'X' {
		return false
	}
	return true
}

func parseUCE(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
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
