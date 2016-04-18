package doc

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func PDF(file *os.File) (*parse.ParsedLayout, error) {

	if !isPDF(file) {
		return nil, nil
	}
	return parsePDF(file)
}

func isPDF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != '%' || b[1] != 'P' || b[2] != 'D' || b[3] != 'F' {
		return false
	}

	return true
}

func parsePDF(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.Document,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 4, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
			}}}}

	return &res, nil
}
