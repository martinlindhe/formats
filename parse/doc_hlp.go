package parse

// HLP help file (Windows)
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func HLP(file *os.File) (*ParsedLayout, error) {

	if !isHLP(file) {
		return nil, nil
	}
	return parseHLP(file)
}

func isHLP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// TODO what is right magic bytes? just guessing
	if b[0] != 0x3f || b[1] != 0x5f || b[2] != 3 || b[3] != 0 {
		return false
	}

	return true
}

func parseHLP(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{
		FileKind: Document,
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
