package archive

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func WINIMG(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isWINIMG(file) {
		return nil, nil
	}
	return parseWINIMG(file, pl)
}

func isWINIMG(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [11]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 0xEB || b[1] != 'X' || b[2] != 0x90 {
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
