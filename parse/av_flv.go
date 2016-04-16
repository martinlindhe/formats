package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func FLV(file *os.File) (*ParsedLayout, error) {

	if !isFLV(file) {
		return nil, nil
	}
	return parseFLV(file)
}

func isFLV(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [3]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'F' || b[1] != 'L' || b[2] != 'V' {
		return false
	}

	return true
}

func parseFLV(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{
		FileKind: AudioVideo,
		Layout: []Layout{{
			Offset: 0,
			Length: 3, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: 0, Length: 3, Info: "magic", Type: ASCII},
			}}}}

	return &res, nil
}
