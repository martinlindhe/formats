package bin

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// GnuPG parses the gpg key public ring format
func GnuPG(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isGPG(c.Header) {
		return nil, nil
	}
	return parseGPG(c.File, c.ParsedLayout)
}

func isGPG(b []byte) bool {

	if b[0] != 0x99 || b[1] != 1 {
		return false
	}
	return true
}

func parseGPG(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.MimeType = "application/x-gnupg-keyring"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 2, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
