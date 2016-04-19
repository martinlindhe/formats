package font

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func WOFF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isWOFF(&hdr) {
		return nil, nil
	}
	return parseWOFF(file, pl)
}

func isWOFF(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'w' || b[1] != 'O' || b[2] != 'F' || b[3] != 'F' {
		return false
	}
	return true
}

func parseWOFF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
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