package windows

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// XBF parses the Visual Studio XAML Binary File format
func XBF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isXBF(c.Header) {
		return nil, nil
	}
	return parseXBF(c.File, c.ParsedLayout)
}

func isXBF(b []byte) bool {

	// XXX just guessing
	if b[0] != 'X' || b[1] != 'B' || b[2] != 'F' || b[3] != 0 {
		return false
	}
	return true
}

func parseXBF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
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
