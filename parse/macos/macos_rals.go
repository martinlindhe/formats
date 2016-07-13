package macos

// ??? found in Mac OS X, /System/Library/PrivateFrameworks/GeoServices.framework/Versions/A/Resources/Ales/1.alestrings
// Extensions: .alestrings

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// RALS parses the rals format
func RALS(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isRALS(c.Header) {
		return nil, nil
	}
	return parseRALS(c.File, c.ParsedLayout)
}

func isRALS(b []byte) bool {

	if b[0] != 'R' || b[1] != 'A' || b[2] != 'L' || b[3] != 'S' {
		return false
	}
	return true
}

func parseRALS(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
