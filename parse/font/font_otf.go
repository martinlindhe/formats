package font

// STATUS 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func OTF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isOTF(&hdr) {
		return nil, nil
	}
	return parseOTF(file, pl)
}

func isOTF(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'O' || b[1] != 'T' || b[2] != 'T' || b[3] != 'O' {
		return false
	}
	return true
}

func parseOTF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
