package av

// Core Audio Format (CAF)
// Modern audio format container by Apple, commonly used in OSX

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func CAF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isCAF(&c.Header) {
		return nil, nil
	}
	return parseCAF(c.File, c.ParsedLayout)
}

func isCAF(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'c' || b[1] != 'a' || b[2] != 'f' || b[3] != 'f' {
		return false
	}
	return true
}

func parseCAF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
