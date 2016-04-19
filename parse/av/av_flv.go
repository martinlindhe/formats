package av

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func FLV(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isFLV(file) {
		return nil, nil
	}
	return parseFLV(file, pl)
}

func isFLV(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [3]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'F' || b[1] != 'L' || b[2] != 'V' {
		return false
	}

	return true
}

func parseFLV(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 3, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 3, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
