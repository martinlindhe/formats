package parse

// VBScript Encoded Script File
// https://en.wikipedia.org/wiki/VBScript
// http://fileformats.archiveteam.org/wiki/VBScript

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func VBE(file *os.File) (*ParsedLayout, error) {

	if !isVBE(file) {
		return nil, nil
	}
	return parseVBE(file)
}

func isVBE(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// XXX just guessing
	if b[0] != 0xff || b[1] != 0xfe || b[2] != 0x23 || b[3] != 0 {
		return false
	}
	return true
}

func parseVBE(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)

	res := ParsedLayout{
		FileKind: Executable,
		Layout: []Layout{{
			Offset: offset,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: offset, Length: 4, Info: "magic", Type: Uint32le},
			}}}}

	return &res, nil
}
