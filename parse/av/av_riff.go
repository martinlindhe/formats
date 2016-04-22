package av

// RIFF format (WAV, AVI)
// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func RIFF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isRIFF(c.Header) {
		return nil, nil
	}
	return parseRIFF(c.File, c.ParsedLayout)
}

func isRIFF(b []byte) bool {

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
