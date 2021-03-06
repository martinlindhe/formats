package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// ISO parses the iso format
func ISO(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isISO(c.Header) {
		return nil, nil
	}
	return parseISO(c.File, c.ParsedLayout)
}

func isISO(b []byte) bool {

	pos := 0x8000
	if b[pos] != 1 || b[pos+1] != 'C' || b[pos+2] != 'D' {
		return false
	}
	return true
}

func parseISO(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0x8000)
	pl.FileKind = parse.Archive
	pl.MimeType = "application/x-iso9660-image"
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
