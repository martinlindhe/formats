package macos

// ??? found in Mac OS X, /System/Library/PrivateFrameworks/GeoServices.framework/Versions/A/Resources/v_water-1.vmap4
// Extensions: .vmap4

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func VMAP4(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isVMAP4(c.Header) {
		return nil, nil
	}
	return parseVMAP4(c.File, c.ParsedLayout)
}

func isVMAP4(b []byte) bool {

	if b[0] != 'V' || b[1] != 'M' || b[2] != 'P' || b[3] != '4' {
		return false
	}
	return true
}

func parseVMAP4(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
