package parse

// https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func JAVA(file *os.File) (*ParsedLayout, error) {

	if !isJAVA(file) {
		return nil, nil
	}
	return parseJAVA(file)
}

func isJAVA(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b uint32
	if err := binary.Read(file, binary.BigEndian, &b); err != nil {
		return false
	}
	if b == 0xcafebabe {
		return true
	}

	return false
}

func parseJAVA(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)
	res := ParsedLayout{
		FileKind: Executable,
		Layout: []Layout{{
			Offset: pos,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: Uint32be},
			}}}}

	return &res, nil
}
