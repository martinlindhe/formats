package windows

// holds app signatures for APPX apps

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func P7X(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isP7X(&c.Header) {
		return nil, nil
	}
	return parseP7X(c.File, c.ParsedLayout)
}

func isP7X(hdr *[0xffff]byte) bool {

	b := *hdr
	// XXX just guessing
	if b[0] != 'P' || b[1] != 'K' || b[2] != 'C' || b[3] != 'X' {
		return false
	}
	return true
}

func parseP7X(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
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