package av

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func MIDI(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isMIDI(&c.Header) {
		return nil, nil
	}
	return parseMIDI(c.File, c.ParsedLayout)
}

func isMIDI(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'M' || b[1] != 'T' || b[2] != 'h' || b[3] != 'd' {
		return false
	}
	return true
}

func parseMIDI(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
