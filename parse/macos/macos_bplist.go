package macos

// used by .nib (NeXT Interface Builder)
// cant find any docs

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// BPLIST parses the Binary plist format
func BPLIST(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isBPLIST(c.Header) {
		return nil, nil
	}
	return parseBPLIST(c.File, c.ParsedLayout)
}

func isBPLIST(b []byte) bool {

	if b[0] != 'b' || b[1] != 'p' || b[2] != 'l' || b[3] != 'i' ||
		b[4] != 's' || b[5] != 't' {
		return false
	}
	return true
}

func parseBPLIST(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.MacOSResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 8, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 8, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
