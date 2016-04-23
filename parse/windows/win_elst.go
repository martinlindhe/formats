package windows

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func ELST(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isELST(c.Header) {
		return nil, nil
	}
	return parseELST(c.File, c.ParsedLayout)
}

func isELST(b []byte) bool {

	if b[0] != 'E' || b[1] != 'L' || b[2] != 'S' || b[3] != 'T' {
		return false
	}
	return true
}

func parseELST(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
