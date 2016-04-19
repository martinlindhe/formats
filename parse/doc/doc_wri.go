package doc

// WRI document (Win16)
// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func WRI(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isWRI(&hdr) {
		return nil, nil
	}
	return parseWRI(file, pl)
}

func isWRI(hdr *[0xffff]byte) bool {

	b := *hdr
	// TODO what is right magic bytes? just guessing
	// FIXME IT IS     if data.find(b'\xBE\x00\x00\x00\xAB\x00\x00\x00\x00\x00\x00\x00\x00') == 1
	if b[0] != 0x31 || b[1] != 0xbe || b[2] != 0 || b[3] != 0 {
		return false
	}

	return true
}

func parseWRI(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Document
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
