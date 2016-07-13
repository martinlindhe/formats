package windows

// STATUS: 1%
// Windows Compiled Resources File
// introduced by Microsoft for Windows 8 applications

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// PRI parses the pri format
func PRI(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isPRI(c.Header) {
		return nil, nil
	}
	return parsePRI(c.File, c.ParsedLayout)
}

func isPRI(b []byte) bool {

	// XXX just guessing
	if b[0] != 'm' || b[1] != 'r' || b[2] != 'm' || b[3] != '_' ||
		b[4] != 'p' || b[5] != 'r' || b[6] != 'i' || b[7] != '2' {
		return false
	}
	return true
}

func parsePRI(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 8, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 8, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
