package macos

// STATUS: 1%
// Extensions: .keychain

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// Keychain parses the Mac OS X keychain format
func Keychain(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isKeychain(c.Header) {
		return nil, nil
	}
	return parseKeychain(c.File, c.ParsedLayout)
}

func isKeychain(b []byte) bool {

	if b[0] != 'k' || b[1] != 'y' || b[2] != 'c' || b[3] != 'h' {
		return false
	}
	return true
}

func parseKeychain(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.MacOSResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
