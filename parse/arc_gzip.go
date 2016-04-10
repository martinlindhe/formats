package parse

// STATUS 1% , XXX

import (
	"encoding/binary"
	"os"
)

func GZIP(file *os.File) (*ParsedLayout, error) {

	if !isGZIP(file) {
		return nil, nil
	}
	return parseGZIP(file)
}

func isGZIP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [2]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 0x1f || b[1] != 0x8b {
		return false
	}
	return true
}

func parseGZIP(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 2,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic", Type: Uint16le}, // XXX le/be ?
		},
	})

	return &res, nil
}
