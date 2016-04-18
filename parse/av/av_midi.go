package av

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func MIDI(file *os.File) (*parse.ParsedLayout, error) {

	if !isMIDI(file) {
		return nil, nil
	}
	return parseMIDI(file)
}

func isMIDI(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'M' || b[1] != 'T' || b[2] != 'h' || b[3] != 'd' {
		return false
	}
	return true
}

func parseMIDI(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.AudioVideo,
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
