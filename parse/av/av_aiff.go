package av

// Audio Interchange File Format (AIFF)
// Developed by Apple, popular on Mac OS in the 90's
// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func AIFF(file *os.File) (*parse.ParsedLayout, error) {

	if !isAIFF(file) {
		return nil, nil
	}
	return parseAIFF(file)
}

func isAIFF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'F' || b[1] != 'O' || b[2] != 'R' || b[3] != 'M' {
		return false
	}

	// TODO also detect "AIFF" string

	return true
}

func parseAIFF(file *os.File) (*parse.ParsedLayout, error) {

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
