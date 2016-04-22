package doc

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func PDF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isPDF(c.Header) {
		return nil, nil
	}
	return parsePDF(c.File, c.ParsedLayout)
}

func isPDF(b []byte) bool {

	if b[0] != '%' || b[1] != 'P' || b[2] != 'D' || b[3] != 'F' {
		return false
	}
	return true
}

func parsePDF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Document
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
