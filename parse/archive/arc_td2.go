package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func TD2(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isTD2(&c.Header) {
		return nil, nil
	}
	return parseTD2(c.File, c.ParsedLayout)
}

func isTD2(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 't' || b[1] != 'd' || b[2] != 0 {
		return false
	}
	return true
}

func parseTD2(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 3,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 3, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
