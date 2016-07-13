package av

// Developed by Apple, popular on Mac OS in the 90's

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// AIFF parses the Audio Interchange File Format
func AIFF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isAIFF(c.Header) {
		return nil, nil
	}
	return parseAIFF(c.File, c.ParsedLayout)
}

func isAIFF(b []byte) bool {

	// TODO also detect "AIFF" string
	if b[0] != 'F' || b[1] != 'O' || b[2] != 'R' || b[3] != 'M' {
		return false
	}
	return true
}

func parseAIFF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
	pl.MimeType = "audio/x-aiff"
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
