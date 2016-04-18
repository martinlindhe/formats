package archive

// STATUS 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func TD2(file *os.File) (*parse.ParsedLayout, error) {

	if !isTD2(file) {
		return nil, nil
	}
	return parseTD2(file)
}

func isTD2(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [3]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 't' || b[1] != 'd' || b[2] != 0 {
		return false
	}

	return true
}

func parseTD2(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)

	res := parse.ParsedLayout{
		FileKind: parse.Archive,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 3,
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 3, Info: "magic", Type: parse.Bytes},
			}}}}

	return &res, nil
}
