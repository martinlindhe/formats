package macos

// ??? found in Mac OS X, /System/Library/PrivateFrameworks/GeoServices.framework/Versions/A/Resources/globe-default-454_2x.styl
// Extensions: .styl

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// STYL parses the styl format
func STYL(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isSTYL(c.Header) {
		return nil, nil
	}
	return parseSTYL(c.File, c.ParsedLayout)
}

func isSTYL(b []byte) bool {

	if b[0] != 'S' || b[1] != 'T' || b[2] != 'Y' || b[3] != 'L' {
		return false
	}
	return true
}

func parseSTYL(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
