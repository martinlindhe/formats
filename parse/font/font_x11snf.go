package font

// STATUS: 1%

import (
	"encoding/binary"
	"log"
	"os"

	"github.com/martinlindhe/formats/parse"
)

// X11FontSNF parses the X11 font files in Server Natural Format
func X11FontSNF(c *parse.Checker) (*parse.ParsedLayout, error) {
	if !isX11FontSNF(c.Header) {
		return nil, nil
	}
	return parseX11FontSNF(c.File, c.ParsedLayout)
}

func isX11FontSNF(b []byte) bool {
	val := binary.LittleEndian.Uint32(b)
	if val == 4 {
		return true // le
	}
	if val == 0x04000000 {
		log.Println("sample please! x11 snf font big-endian")
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
