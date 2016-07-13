package windows

// STATUS: 1%
// Extensions: .rbs
// found on Windows 10 ProgramData/Microsoft/Diagnosis/events00.rbs

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// RBS parses the rbs format
func RBS(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isRBS(c.Header) {
		return nil, nil
	}
	return parseRBS(c.File, c.ParsedLayout)
}

func isRBS(b []byte) bool {

	if b[0] != 'U' || b[1] != 'T' || b[2] != 'C' || b[3] != 'R' ||
		b[4] != 'B' || b[5] != 'E' || b[6] != 'S' || b[7] != '3' {
		return false
	}
	return true
}

func parseRBS(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 8, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 8, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
