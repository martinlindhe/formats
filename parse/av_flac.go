package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func FLAC(file *os.File) (*ParsedLayout, error) {

	if !isFLAC(file) {
		return nil, nil
	}
	return parseFLAC(file)
}

func isFLAC(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'f' || b[1] != 'L' || b[2] != 'a' || b[3] != 'C' {
		return false
	}

	return true
}

func parseFLAC(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)
	res := ParsedLayout{
		FileKind: AudioVideo,
		Layout: []Layout{{
			Offset: pos,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: Bytes},
			}}}}

	return &res, nil
}
