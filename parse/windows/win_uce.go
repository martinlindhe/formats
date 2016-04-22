package windows

// ???, found in win10-os/Windows/System32/bopomofo.uce

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func UCE(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isUCE(&c.Header) {
		return nil, nil
	}
	return parseUCE(c.File, c.ParsedLayout)
}

func isUCE(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'U' || b[1] != 'C' || b[2] != 'E' || b[3] != 'X' {
		return false
	}
	return true
}

func parseUCE(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
