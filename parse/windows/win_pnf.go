package windows

// ???
// found on Windows 10 Windows/INF/acpipmi.PNF
// extensions: .pnf

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func PNF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isPNF(c.Header) {
		return nil, nil
	}
	return parsePNF(c.File, c.ParsedLayout)
}

func isPNF(b []byte) bool {

	// XXX just guessing
	if b[0] != 1 || b[1] != 3 || b[2] != 2 || b[3] != 0 ||
		b[4] != 0x83 || b[5] != 0 {
		return false
	}
	return true
}

func parsePNF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 6, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 6, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
