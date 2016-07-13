package font

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// TTF parses the TrueType font format
func TTF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isTTF(c.Header) {
		return nil, nil
	}
	return parseTTF(c.File, c.ParsedLayout)
}

func isTTF(b []byte) bool {

	if b[0] != 0 || b[1] != 1 || b[2] != 0 || b[3] != 0 {
		return false
	}
	return true
}

func parseTTF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
	pl.MimeType = "application/x-font-ttf"
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
