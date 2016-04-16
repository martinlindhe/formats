package parse

// CHM help file (Windows)
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func CHM(file *os.File) (*ParsedLayout, error) {

	if !isCHM(file) {
		return nil, nil
	}
	return parseCHM(file)
}

func isCHM(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// TODO what is right magic bytes? just guessing
	if b[0] != 'I' || b[1] != 'T' || b[2] != 'S' || b[3] != 'F' {
		return false
	}

	return true
}

func parseCHM(file *os.File) (*ParsedLayout, error) {

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
