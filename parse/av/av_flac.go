package av

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func FLAC(file *os.File) (*parse.ParsedLayout, error) {

	if !isFLAC(file) {
		return nil, nil
	}
	return parseFLAC(file)
}

func isFLAC(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'f' || b[1] != 'L' || b[2] != 'a' || b[3] != 'C' {
		return false
	}

	return true
}

func parseFLAC(file *os.File) (*parse.ParsedLayout, error) {

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
