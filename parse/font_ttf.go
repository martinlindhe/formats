package parse

// truetype fonts
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func TTF(file *os.File) (*ParsedLayout, error) {

	if !isTTF(file) {
		return nil, nil
	}
	return parseTTF(file)
}

func isTTF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 0 || b[1] != 1 || b[2] != 0 || b[3] != 0 {
		return false
	}

	return true
}

func parseTTF(file *os.File) (*ParsedLayout, error) {

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
