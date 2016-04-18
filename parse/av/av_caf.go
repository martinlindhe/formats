package av

// STATUS: 1%
// Core Audio Format (CAF)
// Modern audio format container by Apple, commonly used in OSX

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func CAF(file *os.File) (*parse.ParsedLayout, error) {

	if !isCAF(file) {
		return nil, nil
	}
	return parseCAF(file)
}

func isCAF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'c' || b[1] != 'a' || b[2] != 'f' || b[3] != 'f' {
		return false
	}

	return true
}

func parseCAF(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.AudioVideo,
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
