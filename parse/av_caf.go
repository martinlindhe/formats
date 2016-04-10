package parse

// STATUS: 1%
// Core Audio Format (CAF)
// Modern audio format container by Apple, commonly used in OSX

import (
	"encoding/binary"
	"os"
)

func CAF(file *os.File) (*ParsedLayout, error) {

	if !isCAF(file) {
		return nil, nil
	}
	return parseCAF(file)
}

func isCAF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'c' || b[1] != 'a' || b[2] != 'f' || b[3] != 'f' {
		return false
	}

	return true
}

func parseCAF(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 4, // XXX
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 4, Info: "magic", Type: Bytes},
		}})
	return &res, nil
}
