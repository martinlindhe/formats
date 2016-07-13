package bin

// variant of RIFF
// http://fileformats.archiveteam.org/wiki/RIFX

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// RIFX parses the rifx format
func RIFX(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isRIFX(c.Header) {
		return nil, nil
	}
	return parseRIFX(c.File, c.ParsedLayout)
}

func isRIFX(b []byte) bool {

	s := string(b[0:4])
	return s == "RIFX" || s == "XFIR"
}

func parseRIFX(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	// XXX: "RIFX" = big endian, XFIR = little endian

	pos := int64(0)
	pl.FileKind = parse.Binary
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
