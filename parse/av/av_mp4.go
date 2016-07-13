package av

// A/V container format
// .mp4, .mov (quicktime)

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// MP4 parses the mp4 format
func MP4(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isMP4(c.Header) {
		return nil, nil
	}
	return parseMP4(c.File, c.ParsedLayout)
}

func isMP4(b []byte) bool {

	// XXX
	return false
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
