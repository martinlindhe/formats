package doc

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// HLP parses the Windows HLP help file
func HLP(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isHLP(c.Header) {
		return nil, nil
	}
	return parseHLP(c.File, c.ParsedLayout)
}

func isHLP(b []byte) bool {

	// TODO what is right magic bytes? just guessing
	if b[0] != 0x3f || b[1] != 0x5f || b[2] != 3 || b[3] != 0 {
		return false
	}
	return true
}

func parseHLP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
