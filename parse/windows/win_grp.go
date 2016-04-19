package windows

// GRP

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func GRP(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isGRP(&hdr) {
		return nil, nil
	}
	return parseGRP(file, pl)
}

func isGRP(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'P' || b[1] != 'M' || b[2] != 'C' || b[3] != 'C' {
		return false
	}
	return true
}

func parseGRP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}