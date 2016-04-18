package doc

// CHM help file (Windows)
// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func CHM(file *os.File) (*parse.ParsedLayout, error) {

	if !isCHM(file) {
		return nil, nil
	}
	return parseCHM(file)
}

func isCHM(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// TODO what is right magic bytes? just guessing
	if b[0] != 'I' || b[1] != 'T' || b[2] != 'S' || b[3] != 'F' {
		return false
	}

	return true
}

func parseCHM(file *os.File) (*parse.ParsedLayout, error) {

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
