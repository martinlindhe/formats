package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// LUKS parse the luks format
func LUKS(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isLUKS(c.Header) {
		return nil, nil
	}
	return parseLUKS(c.File, c.ParsedLayout)
}

func isLUKS(b []byte) bool {

	if b[0] != 'L' || b[1] != 'U' || b[2] != 'K' || b[3] != 'S' {
		return false
	}
	return true
}

func parseLUKS(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	pl.FileKind = parse.Archive
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
