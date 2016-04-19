package av

// Audio Interchange File Format (AIFF)
// Developed by Apple, popular on Mac OS in the 90's

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func AIFF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isAIFF(&hdr) {
		return nil, nil
	}
	return parseAIFF(file, pl)
}

func isAIFF(hdr *[0xffff]byte) bool {

	// TODO also detect "AIFF" string
	b := *hdr
	if b[0] != 'F' || b[1] != 'O' || b[2] != 'R' || b[3] != 'M' {
		return false
	}
	return true
}

func parseAIFF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
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
