package parse

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func BMP(file *os.File) *ParsedLayout {

	if !isBMP(file) {
		return nil
	}

	fmt.Println("XXX TODO parse BMP")

	return parseBMP(file)
}

func parseBMP(file *os.File) *ParsedLayout {

	res := ParsedLayout{}

	bitmapFileHeader := Layout{
		Length: 14,
		Info:   "bitmap file header",
		Childs: []Layout{
			Layout{Offset: 0, Length: 2, Info: "magic (BMP image)", Type: ASCII},
			Layout{Offset: 2, Length: 4, Info: "file size", Type: Uint32le},
			Layout{Offset: 6, Length: 2, Info: "reserved 1", Type: Uint16le},
			Layout{Offset: 8, Length: 2, Info: "reserved 2", Type: Uint16le},
			Layout{Offset: 10, Length: 4, Info: "offset to pixel data", Type: Uint32le},
		},
	}
	res.Layout = append(res.Layout, bitmapFileHeader)

	return &res
}

func isBMP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	r := io.Reader(file)
	var b [2]byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return false
	}
	return b[0] == 'B' && b[1] == 'M'
}
