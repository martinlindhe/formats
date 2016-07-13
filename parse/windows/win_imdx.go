package windows

// STATUS: 1%
// Extensions: .lex
// ??? found on Windows 10, Windows/InputMethod/CHS/ChsPinyinDM10.lex

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// IMDX parses the imdx format
func IMDX(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isIMDX(c.Header) {
		return nil, nil
	}
	return parseIMDX(c.File, c.ParsedLayout)
}

func isIMDX(b []byte) bool {

	if b[0] != 'I' || b[1] != 'M' || b[2] != 'D' || b[3] != 'X' {
		return false
	}
	return true
}

func parseIMDX(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
