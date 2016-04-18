package av

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func OGG(file *os.File) (*parse.ParsedLayout, error) {

	if !isOGG(file) {
		return nil, nil
	}
	return parseOGG(file)
}

func isOGG(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'O' || b[1] != 'g' || b[2] != 'g' {
		return false
	}

	return true
}

func parseOGG(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.AudioVideo,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 3, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 3, Info: "magic", Type: parse.ASCII},
			}}}}

	return &res, nil
}
