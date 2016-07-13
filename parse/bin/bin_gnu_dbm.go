package bin

// GDBM magic numbers
// /usr/share/file/magic/database

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// GnuDBM parses the gdbm format
func GnuDBM(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isGnuDBM(c.Header) {
		return nil, nil
	}
	return parseGnuDBM(c.File, c.ParsedLayout)
}

func isGnuDBM(b []byte) bool {

	if b[0] == 0xce && b[1] == 0x9a && b[2] == 0x57 && b[3] == 0x13 {
		return true // v1.x le
	}
	if b[3] == 0xce && b[2] == 0x9a && b[1] == 0x57 && b[0] == 0x13 {
		return true // v1.x be
	}
	if b[3] == 'G' && b[2] == 'D' && b[1] == 'B' && b[0] == 'M' {
		return true // v v2.x
	}
	return false
}

func parseGnuDBM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.MimeType = "application/x-gdbm"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
