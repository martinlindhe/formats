package windows

// ???
// found on Windows 10 ProgramData/Microsoft/Diagnosis/events00.rbs
// extensions: .rbs

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func RBS(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isRBS(&c.Header) {
		return nil, nil
	}
	return parseRBS(c.File, c.ParsedLayout)
}

func isRBS(hdr *[0xffff]byte) bool {

	b := *hdr
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
