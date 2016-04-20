package macos

//  Mac OS X bill of materials (BOM) file

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func BOMStore(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isBOMStore(&c.Header) {
		return nil, nil
	}
	return parseBOMStore(c.File, c.ParsedLayout)
}

func isBOMStore(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'B' || b[1] != 'O' || b[2] != 'M' || b[3] != 'S' ||
		b[4] != 't' || b[5] != 'o' || b[6] != 'r' || b[7] != 'e' {
		return false
	}
	return true
}

func parseBOMStore(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
