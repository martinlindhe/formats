package macos

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// CodeSignature parses the code signature format
func CodeSignature(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isCodeSignature(c.Header) {
		return nil, nil
	}
	return parseCodeSignature(c.File, c.ParsedLayout)
}

func isCodeSignature(b []byte) bool {

	// XXX just guessing
	if b[0] != 0x30 || b[1] != 0x80 || b[2] != 6 || b[3] != 9 {
		return false
	}
	return true
}

func parseCodeSignature(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.MacOSResource
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
