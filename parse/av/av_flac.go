package av

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func FLAC(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isFLAC(c.Header) {
		return nil, nil
	}
	return parseFLAC(c.File, c.ParsedLayout)
}

func isFLAC(b []byte) bool {

	if b[0] != 'f' || b[1] != 'L' || b[2] != 'a' || b[3] != 'C' {
		return false
	}
	return true
}

func parseFLAC(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
	pl.MimeType = "audio/x-flac"
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
