package windows

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// GRP parses the grp format
func GRP(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isGRP(c.Header) {
		return nil, nil
	}
	return parseGRP(c.File, c.ParsedLayout)
}

func isGRP(b []byte) bool {

	if b[0] != 'P' || b[1] != 'M' || b[2] != 'C' || b[3] != 'C' {
		return false
	}
	return true
}

func parseGRP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
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
