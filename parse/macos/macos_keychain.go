package macos

//  Mac OS X .keychain

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func Keychain(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isKeychain(&c.Header) {
		return nil, nil
	}
	return parseKeychain(c.File, c.ParsedLayout)
}

func isKeychain(hdr *[0xffff]byte) bool {

	b := *hdr
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
