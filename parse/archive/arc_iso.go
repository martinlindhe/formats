package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func ISO(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isISO(&hdr) {
		return nil, nil
	}
	return parseISO(file, pl)
}

func isISO(hdr *[0xffff]byte) bool {

	pos := 0x8000
	b := *hdr
	if b[pos] != 1 || b[pos+1] != 'C' || b[pos+2] != 'D' {
		return false
	}
	return true
}

func parseISO(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0x8000)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 3,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 3, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}