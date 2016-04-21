package bin

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func GPG(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isGPG(&c.Header) {
		return nil, nil
	}
	return parseGPG(c.File, c.ParsedLayout)
}

func isGPG(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 0x99 || b[1] != 1 {
		return false
	}
	return true
}

func parseGPG(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 2, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
