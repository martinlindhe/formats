package parse

// STATUS: 1% TODO implement format

import (
	"encoding/binary"
	"os"
)

func TIFF(file *os.File) (*ParsedLayout, error) {

	if !isTIFF(file) {
		return nil, nil
	}
	return parseTIFF(file)
}

func isTIFF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// XXX dont know magic numbers just guessing
	if b[0] != 'I' || b[1] != 'I' || b[2] != '*' || b[3] != 0 {
		return false
	}
	return true
}

func parseTIFF(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{}

	fileHeader := Layout{
		Offset: 0,
		Length: 4,
		Info:   "file header",
		Type:   Group,
		Childs: []Layout{
			Layout{Offset: 0, Length: 4, Info: "magic", Type: Bytes},
		},
	}

	res.Layout = append(res.Layout, fileHeader)

	return &res, nil
}
