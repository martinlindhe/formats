package doc

// Rich Type File (RTF)

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func RTF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isRTF(&hdr) {
		return nil, nil
	}
	return parseRTF(file, pl)
}

func isRTF(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != '{' || b[1] != '\\' || b[2] != 'r' || b[3] != 't' || b[4] != 'f' {
		return false
	}
	return true
}

func parseRTF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Document
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 5, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 5, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
