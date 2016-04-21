package bin

// timezone data
// /usr/share/file/magic/timezone

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func Timezone(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isTimezone(&c.Header) {
		return nil, nil
	}
	return parseTimezone(c.File, c.ParsedLayout)
}

func isTimezone(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] == 'T' && b[1] == 'Z' && b[2] == 'i' && b[3] == 'f' {
		return true
	}
	return false
}

func parseTimezone(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			// XXX this is le-format only
			{Offset: pos, Length: 1, Info: "version", Type: parse.Bytes},
			{Offset: pos + 1, Length: 3, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
