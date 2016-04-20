package macos

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func CodeDirectory(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isCodeDirectory(&c.Header) {
		return nil, nil
	}
	return parseCodeDirectory(c.File, c.ParsedLayout)
}

func isCodeDirectory(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 0xfa || b[1] != 0xde || b[2] != 0x0c {
		if b[3] != 0x01 && b[3] != 0x02 {
			return false
		}
	}
	return true
}

func parseCodeDirectory(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	// byte 3 = 1: CodeRequirements file
	// byte 3 = 2: CodeDirectory file

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
