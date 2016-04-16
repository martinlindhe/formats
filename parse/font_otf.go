package parse

// STATUS 1%

import (
	"encoding/binary"
	"os"
)

func OTF(file *os.File) (*ParsedLayout, error) {

	if !isOTF(file) {
		return nil, nil
	}
	return parseOTF(file)
}

func isOTF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'O' || b[1] != 'T' || b[2] != 'T' || b[3] != 'O' {
		return false
	}

	return true
}

func parseOTF(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{
		FileKind: Font,
		Layout: []Layout{{
			Offset: 0,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: 0, Length: 4, Info: "magic", Type: Bytes},
			}}}}

	return &res, nil
}
