package font

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// OTF parses the otf format
func OTF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isOTF(c.Header) {
		return nil, nil
	}
	return parseOTF(c.File, c.ParsedLayout)
}

func isOTF(b []byte) bool {

	if b[0] != 'O' || b[1] != 'T' || b[2] != 'T' || b[3] != 'O' {
		return false
	}
	return true
}

func parseOTF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
	pl.MimeType = "application/vnd.ms-opentype"
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
