package parse

// STATUS 1%, see https://golang.org/src/compress/bzip2/bzip2.go

import (
	"encoding/binary"
	"os"
)

func BZIP2(file *os.File) (*ParsedLayout, error) {

	if !isBZIP2(file) {
		return nil, nil
	}
	return parseBZIP2(file)
}

func isBZIP2(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [2]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'B' || b[1] != 'Z' {
		return false
	}

	return true
}

func parseBZIP2(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	res.Layout = append(res.Layout, Layout{
		Offset: 0,
		Length: 4,
		Info:   "header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic", Type: ASCII},
			Layout{Offset: 2, Length: 1, Info: "encoding", Type: Uint8},          // XXX h = huffman
			Layout{Offset: 3, Length: 1, Info: "compression level", Type: ASCII}, // 0=worst, 9=best<
		},
	})

	return &res, nil
}
