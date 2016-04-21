package font

// X11 font files in Server Natural Format (SNF)

// STATUS: 1%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func X11FontSNF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isX11FontSNF(c.File) {
		return nil, nil
	}
	return parseX11FontSNF(c.File, c.ParsedLayout)
}

func isX11FontSNF(file *os.File) bool {

	val, _ := parse.ReadUint32le(file, 0)
	if val == 4 {
		return true // le
	}
	if val == 0x04000000 {
		fmt.Println("sample please! x11 snf font big-endian")
		return true // be
	}
	return false
}

func parseX11FontSNF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
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