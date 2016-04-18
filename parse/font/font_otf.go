package font

// STATUS 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func OTF(file *os.File) (*parse.ParsedLayout, error) {

	if !isOTF(file) {
		return nil, nil
	}
	return parseOTF(file)
}

func isOTF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'O' || b[1] != 'T' || b[2] != 'T' || b[3] != 'O' {
		return false
	}

	return true
}

func parseOTF(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.Font,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 4, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
			}}}}

	return &res, nil
}
