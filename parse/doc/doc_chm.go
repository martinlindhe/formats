package doc

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// CHM parses the Windows CHM help file format
func CHM(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isCHM(c.Header) {
		return nil, nil
	}
	return parseCHM(c.File, c.ParsedLayout)
}

func isCHM(b []byte) bool {

	// TODO what is right magic bytes? just guessing
	if b[0] != 'I' || b[1] != 'T' || b[2] != 'S' || b[3] != 'F' {
		return false
	}
	return true
}

func parseCHM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
