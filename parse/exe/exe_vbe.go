package exe

// VBScript Encoded Script File

// http://lifeinhex.com/tag/vbe/
// https://en.wikipedia.org/wiki/VBScript
// http://fileformats.archiveteam.org/wiki/VBScript

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func VBE(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isVBE(c.Header) {
		return nil, nil
	}
	return parseVBE(c.File, c.ParsedLayout)
}

func isVBE(b []byte) bool {

	// XXX just guessing
	if b[0] != 0xff || b[1] != 0xfe || b[2] != 0x23 || b[3] != 0 {
		return false
	}
	return true
}

func parseVBE(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
