package image

// Better Portable Graphics
// https://en.wikipedia.org/wiki/Better_Portable_Graphics
// Extensions: .bpg

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func BPG(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isBPG(c.Header) {
		return nil, nil
	}
	return parseBPG(c.File, c.ParsedLayout)
}

func isBPG(b []byte) bool {

	if b[0] == 'B' && b[1] == 'P' && b[2] == 'G' && b[3] == 0xfb {
		return true
	}
	return false
}

func parseBPG(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Image
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
