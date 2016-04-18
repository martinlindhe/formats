package doc

// HLP help file (Windows)
// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func HLP(file *os.File) (*parse.ParsedLayout, error) {

	if !isHLP(file) {
		return nil, nil
	}
	return parseHLP(file)
}

func isHLP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// TODO what is right magic bytes? just guessing
	if b[0] != 0x3f || b[1] != 0x5f || b[2] != 3 || b[3] != 0 {
		return false
	}

	return true
}

func parseHLP(file *os.File) (*parse.ParsedLayout, error) {

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
