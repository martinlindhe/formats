package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func MP4(file *os.File) (*ParsedLayout, error) {

	if !isMP4(file) {
		return nil, nil
	}
	return parseMP4(file)
}

func isMP4(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// TODO what is right magic bytes? just guessing
	if b[0] != 0 || b[1] != 0 || b[2] != 0 || b[3] != 0x18 {
		return false
	}

	return true
}

func parseMP4(file *os.File) (*ParsedLayout, error) {

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
