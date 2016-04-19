package font

// Adobe Printer Font Binary (used in the '90s)
// STATUS: 1%

import (
	"github.com/martinlindhe/formats/parse"
	"os"
)

func PFB(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isPFB(file) {
		return nil, nil
	}
	return parsePFB(file, pl)
}

func isPFB(file *os.File) bool {

	// XXX need ways to work on hdr []byte

	// XXX just guessing ...
	s, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, 6, 10)
	if s != "%!FontType" {
		return false
	}
	return true
}

func parsePFB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 16, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 6, Info: "unknown", Type: parse.Bytes},
			{Offset: pos + 6, Length: 10, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
