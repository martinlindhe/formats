package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func OGG(file *os.File) (*ParsedLayout, error) {

	if !isOGG(file) {
		return nil, nil
	}
	return parseOGG(file)
}

func isOGG(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'O' || b[1] != 'g' || b[2] != 'g' {
		return false
	}

	return true
}

func parseOGG(file *os.File) (*ParsedLayout, error) {

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
