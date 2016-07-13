package macos

// ??? found in Mac OS X, /System/Library/Frameworks/SceneKit.framework/Versions/A/Resources/default.metallib
// Extensions: .metallib

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// MTLB parses the metallib format
func MTLB(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isMTLB(c.Header) {
		return nil, nil
	}
	return parseMTLB(c.File, c.ParsedLayout)
}

func isMTLB(b []byte) bool {

	if b[0] != 'M' || b[1] != 'T' || b[2] != 'L' || b[3] != 'B' {
		return false
	}
	return true
}

func parseMTLB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
