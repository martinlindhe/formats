package parse

// WRI document (Win16)
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func WRI(file *os.File) (*ParsedLayout, error) {

	if !isWRI(file) {
		return nil, nil
	}
	return parseWRI(file)
}

func isWRI(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [5]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// TODO what is right magic bytes? just guessing
	// FIXME IT IS     if data.find(b'\xBE\x00\x00\x00\xAB\x00\x00\x00\x00\x00\x00\x00\x00') == 1
	if b[0] != 0x31 || b[1] != 0xbe || b[2] != 0 || b[3] != 0 {
		return false
	}

	return true
}

func parseWRI(file *os.File) (*ParsedLayout, error) {

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
