package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// GZIP parses the gz format
func GZIP(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isGZIP(c.Header) {
		return nil, nil
	}
	return parseGZIP(c.File, c.ParsedLayout)
}

func isGZIP(b []byte) bool {

	if b[0] != 0x1f || b[1] != 0x8b {
		return false
	}
	return true
}

func parseGZIP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.MimeType = "application/x-gzip"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 2,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.Uint16le}, // XXX le/be ?
		}}}

	return &pl, nil
}
