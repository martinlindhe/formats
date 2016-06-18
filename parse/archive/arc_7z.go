package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func SevenZIP(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isSevenZIP(c.Header) {
		return nil, nil
	}
	return parseSevenZIP(c.File, c.ParsedLayout)
}

func isSevenZIP(b []byte) bool {

	if b[0] != '7' || b[1] != 'z' || b[2] != 0xbc || b[3] != 0xaf ||
		b[4] != 0x27 || b[5] != 0x1c {
		return false
	}
	return true
}

func parseSevenZIP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 6,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 6, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
