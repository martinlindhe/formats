package av

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// FLV parses the flv format
func FLV(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isFLV(c.Header) {
		return nil, nil
	}
	return parseFLV(c.File, c.ParsedLayout)
}

func isFLV(b []byte) bool {

	if b[0] != 'F' || b[1] != 'L' || b[2] != 'V' {
		return false
	}
	return true
}

func parseFLV(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
	pl.MimeType = "video/x-flv"
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
