package font

// STATUS: 1%
// used in the '90s

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// PFB parses the Adobe Printer Font Binary format
func PFB(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isPFB(c.Header) {
		return nil, nil
	}
	return parsePFB(c.File, c.ParsedLayout)
}

func isPFB(b []byte) bool {

	// XXX just guessing ...
	s := string(b[6:16])
	return s == "%!FontType"
}

func parsePFB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 16, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 6, Info: "unknown", Type: parse.Bytes},
			{Offset: pos + 6, Length: 10, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
