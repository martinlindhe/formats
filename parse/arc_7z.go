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
	var b [6]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != '7' || b[1] != 'z' || b[2] != 0xbc || b[3] != 0xaf || b[4] != 0x27 || b[5] != 0x1c {
		return false
	}

	return true
}

func parseSEVENZIP(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 6,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 6, Info: "magic", Type: Bytes},
		},
	})

	return &res, nil
}
