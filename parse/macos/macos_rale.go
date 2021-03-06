package macos

// ??? found in Mac OS X, /System/Library/PrivateFrameworks/GeoServices.framework/Versions/A/Resources/Ales/1.ale
// Extensions: .ale

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// RALE parses the rale format
func RALE(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isRALE(c.Header) {
		return nil, nil
	}
	return parseRALE(c.File, c.ParsedLayout)
}

func isRALE(b []byte) bool {

	if b[0] != 'R' || b[1] != 'A' || b[2] != 'L' || b[3] != 'E' {
		return false
	}
	return true
}

func parseRALE(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
