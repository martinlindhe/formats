package image

// X11 mouse cursor image

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func XCursor(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isXCursor(&c.Header) {
		return nil, nil
	}
	return parseXCursor(c.File, c.ParsedLayout)
}

func isXCursor(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] == 'X' && b[1] == 'c' && b[2] == 'u' && b[3] == 'r' {
		return true
	}
	return false
}

func parseXCursor(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Image
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
