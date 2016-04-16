package parse

// MS Word document (.doc)
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func WORD(file *os.File) (*ParsedLayout, error) {

	if !isWORD(file) {
		return nil, nil
	}
	return parseWORD(file)
}

func isWORD(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [5]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// TODO what is right magic bytes? just guessing
	if b[0] != 0xd0 || b[1] != 0xcf || b[2] != 0x11 || b[3] != 0xe0 {
		return false
	}

	return true
}

func parseWORD(file *os.File) (*ParsedLayout, error) {

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
