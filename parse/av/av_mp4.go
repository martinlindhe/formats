package av

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func MP4(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isMP4(&c.Header) {
		return nil, nil
	}
	return parseMP4(c.File, c.ParsedLayout)
}

func isMP4(hdr *[0xffff]byte) bool {

	b := *hdr
	// TODO what is right magic bytes? just guessing
	if b[0] != 0 || b[1] != 0 || b[2] != 0 || b[3] != 0x18 {
		return false
	}
	return true
}

func parseMP4(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
