package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func WINIMG(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isWINIMG(&hdr) {
		return nil, nil
	}
	return parseWINIMG(file, pl)
}

func isWINIMG(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 0xeb || b[1] != 'X' || b[2] != 0x90 {
		return false
	}
	if b[3] != 'W' || b[4] != 'I' || b[5] != 'N' || b[6] != 'I' ||
		b[7] != 'M' || b[8] != 'A' || b[9] != 'G' || b[10] != 'E' {
		return false
	}
	return true
}

func parseWINIMG(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 2,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.Uint16le}, // XXX le/be ?
		}}}

	return &pl, nil
}
