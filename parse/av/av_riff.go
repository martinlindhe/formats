package av

// RIFF format (WAV, AVI)
// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func RIFF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isRIFF(&hdr) {
		return nil, nil
	}
	return parseRIFF(file, pl)
}

func isRIFF(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'R' || b[1] != 'I' || b[2] != 'F' || b[3] != 'F' {
		return false
	}
	return true
}

func parseRIFF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
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
