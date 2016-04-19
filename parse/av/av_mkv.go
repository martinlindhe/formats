package av

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func MKV(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isMKV(&c.Header) {
		return nil, nil
	}
	return parseMKV(c.File, c.ParsedLayout)
}

func isMKV(hdr *[0xffff]byte) bool {

	b := *hdr
	// XXX what is magic sequence? just guessing
	if b[0] != 0x1a || b[1] != 0x45 || b[2] != 0xdf || b[3] != 0xa3 {
		return false
	}
	return true
}

func parseMKV(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
