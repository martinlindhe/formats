package archive

// Xar Archive
// Extensions: .xar .pkg

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func XAR(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isXAR(c.Header) {
		return nil, nil
	}
	return parseXAR(c.File, c.ParsedLayout)
}

func isXAR(b []byte) bool {

	if b[0] != 'x' || b[1] != 'a' || b[2] != 'r' || b[3] != '!' {
		return false
	}
	return true
}

func parseXAR(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
