package windows

// STATUS: 1%
// Extensions: .pf
// found on Windows 10 Windows/Prefetch/SEARCHINDEXER.EXE-4A6353B9.pf

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// MAM parses the mam format
func MAM(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isMAM(c.Header) {
		return nil, nil
	}
	return parseMAM(c.File, c.ParsedLayout)
}

func isMAM(b []byte) bool {

	// XXX just guessing
	if b[0] != 'M' || b[1] != 'A' || b[2] != 'M' || b[3] != 4 {
		return false
	}
	return true
}

func parseMAM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
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
