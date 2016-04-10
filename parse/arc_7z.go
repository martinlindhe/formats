package parse

// STATUS 1% , XXX

import (
	"encoding/binary"
	"os"
)

func SEVENZIP(file *os.File) (*ParsedLayout, error) {

	if !isSEVENZIP(file) {
		return nil, nil
	}
	return parseSEVENZIP(file)
}

func isSEVENZIP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [2]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != '7' || b[1] != 'z' {
		return false
	}

	return true
}

func parseSEVENZIP(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 2,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic", Type: ASCII},
		},
	})

	return &res, nil
}
