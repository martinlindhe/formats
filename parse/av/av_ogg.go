package av

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func OGG(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isOGG(&hdr) {
		return nil, nil
	}
	return parseOGG(file, pl)
}

func isOGG(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'O' || b[1] != 'g' || b[2] != 'g' {
		return false
	}
	return true
}

func parseOGG(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 3, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 3, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}