package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func ELF(file *os.File) (*ParsedLayout, error) {

	if !isELF(file) {
		return nil, nil
	}
	return parseELF(file)
}

func isELF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] == 0x7f && b[1] == 'E' && b[2] == 'L' && b[3] == 'F' {
		return true
	}

	return false
}

func parseELF(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)

	res := ParsedLayout{
		FileKind: Executable,
		Layout: []Layout{{
			Offset: offset,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: offset, Length: 4, Info: "magic", Type: ASCII},
			}}}}

	return &res, nil
}
