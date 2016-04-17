package parse

// STATUS: borked

import (
	"encoding/binary"
	"os"
)

func TTC(file *os.File) (*ParsedLayout, error) {

	if !isTTC(file) {
		return nil, nil
	}
	return parseTTC(file)
}

func isTTC(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 't' || b[1] != 't' || b[2] != 'c' || b[3] != 'f' {
		return false
	}

	return true
}

func parseTTC(file *os.File) (*ParsedLayout, error) {

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
