package font

// TrueType font

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func TTF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isTTF(&hdr) {
		return nil, nil
	}
	return parseTTF(file, pl)
}

func isTTF(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 0 || b[1] != 1 || b[2] != 0 || b[3] != 0 {
		return false
	}
	return true
}

func parseTTF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
