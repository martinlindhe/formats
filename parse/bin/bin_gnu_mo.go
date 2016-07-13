package bin

// GNU nlsutils message catalog (.mo)
// /usr/share/file/magic/gnu

// TODO format can be big endian, depending on first 4 bytes, need sample!

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// GnuMO parses the GNU mo format
func GnuMO(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isGnuMO(c.Header) {
		return nil, nil
	}
	return parseGnuMO(c.File, c.ParsedLayout)
}

func isGnuMO(b []byte) bool {

	if b[0] == 0xde && b[1] == 0x12 && b[2] == 4 && b[3] == 0x95 {
		return true // le
	}
	if b[3] == 0xde && b[2] == 0x12 && b[1] == 4 && b[0] == 0x95 {
		return true // be
	}
	return false
}

func parseGnuMO(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 12, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
			{Offset: pos, Length: 4, Info: "revision id", Type: parse.Uint32le},
			{Offset: pos, Length: 4, Info: "messages", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
