package windows

// ???, Windows OS resource file
// extension .cat

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func CAT(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isCAT(c.Header) {
		return nil, nil
	}
	return parseCAT(c.File, c.ParsedLayout)
}

func isCAT(b []byte) bool {

	if b[0] != 0x30 || b[1] != 0x82 {
		return false
	}
	return true
}

func parseCAT(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 2, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
