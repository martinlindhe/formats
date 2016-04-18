package archive

// STATUS 1%, see https://golang.org/src/compress/bzip2/bzip2.go

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func BZIP2(file *os.File) (*parse.ParsedLayout, error) {

	if !isBZIP2(file) {
		return nil, nil
	}
	return parseBZIP2(file)
}

func isBZIP2(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [3]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'B' || b[1] != 'Z' {
		return false
	}
	if b[2] != 'h' {
		// NOTE: onlu huffman encoding is used in the format (?)
		return false
	}

	return true
}

func parseBZIP2(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)

	res := parse.ParsedLayout{
		FileKind: parse.Archive,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 4,
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 2, Info: "magic", Type: parse.ASCII},
				{Offset: pos + 2, Length: 1, Info: "encoding", Type: parse.Uint8},          // XXX h = huffman
				{Offset: pos + 3, Length: 1, Info: "compression level", Type: parse.ASCII}, // 0=worst, 9=best<
			}}}}

	return &res, nil
}
