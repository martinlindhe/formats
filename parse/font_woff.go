package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func WOFF(file *os.File) (*ParsedLayout, error) {

	if !isWOFF(file) {
		return nil, nil
	}
	return parseWOFF(file)
}

func isWOFF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'w' || b[1] != 'O' || b[2] != 'F' || b[3] != 'F' {
		return false
	}

	return true
}

func parseWOFF(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{
		FileKind: Font,
		Layout: []Layout{{
			Offset: 0,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: 0, Length: 4, Info: "magic", Type: ASCII},
			}}}}

	return &res, nil
}
