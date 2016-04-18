package parse

// STATUS 1% , XXX

import (
	"encoding/binary"
	"os"
)

func TD2(file *os.File) (*ParsedLayout, error) {

	if !isTD2(file) {
		return nil, nil
	}
	return parseTD2(file)
}

func isTD2(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [3]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 't' || b[1] != 'd' || b[2] != 0 {
		return false
	}

	return true
}

func parseTD2(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)

	res := ParsedLayout{
		FileKind: Archive,
		Layout: []Layout{{
			Offset: pos,
			Length: 3,
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 3, Info: "magic", Type: Bytes},
			}}}}

	return &res, nil
}
