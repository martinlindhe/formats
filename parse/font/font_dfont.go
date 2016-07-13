package font

// Used on MacOS
// https://en.wikipedia.org/wiki/Datafork_TrueType

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// DFont parses the Datafork TrueType format
func DFont(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isDFont(c.Header) {
		return nil, nil
	}
	return parseDFont(c.File, c.ParsedLayout)
}

func isDFont(b []byte) bool {

	// XXX just guessing
	if b[0] != 0 || b[1] != 0 || b[2] != 1 || b[3] != 0 || b[4] != 0 {
		return false
	}
	return true
}

func parseDFont(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
	pl.MimeType = "application/x-dfont"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 5, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 5, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
