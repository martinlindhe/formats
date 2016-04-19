package doc

// Rich Type File (RTF)
// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func RTF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isRTF(file) {
		return nil, nil
	}
	return parseRTF(file, pl)
}

func isRTF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [5]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

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
