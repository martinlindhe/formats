package font

// STATUS: borked

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func TTC(c *parse.ParseChecker)(*parse.ParsedLayout, error) {

	if !isTTC(&c.Header) {
		return nil, nil
	}
	return parseTTC(c.File, c.ParsedLayout)
}

func isTTC(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 't' || b[1] != 't' || b[2] != 'c' || b[3] != 'f' {
		return false
	}
	return true
}

func parseTTC(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
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
