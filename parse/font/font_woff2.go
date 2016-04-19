package font

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func WOFF2(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isWOFF2(file) {
		return nil, nil
	}
	return parseWOFF2(file, pl)
}

func isWOFF2(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'w' || b[1] != 'O' || b[2] != 'F' || b[3] != '2' {
		return false
	}

	return true
}

func parseWOFF2(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
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
