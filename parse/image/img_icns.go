package image

// XXX some osx/ios icon collection ?

// STATUS 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func ICNS(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isICNS(c.Header) {
		return nil, nil
	}
	return parseICNS(c.File, c.ParsedLayout)
}

func isICNS(b []byte) bool {

	if b[0] != 'i' || b[1] != 'c' || b[2] != 'n' || b[3] != 's' {
		return false
	}
	return true
}

func parseICNS(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Image
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
