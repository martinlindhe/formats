package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func PDF(file *os.File) (*ParsedLayout, error) {

	if !isPDF(file) {
		return nil, nil
	}
	return parsePDF(file)
}

func isPDF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != '%' || b[1] != 'P' || b[2] != 'D' || b[3] != 'F' {
		return false
	}

	return true
}

func parsePDF(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 4, // XXX
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 4, Info: "magic", Type: ASCII},
		}})
	return &res, nil
}
