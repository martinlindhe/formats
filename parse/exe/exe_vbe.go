package exe

// VBScript Encoded Script File
// https://en.wikipedia.org/wiki/VBScript
// http://fileformats.archiveteam.org/wiki/VBScript

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func VBE(file *os.File) (*parse.ParsedLayout, error) {

	if !isVBE(file) {
		return nil, nil
	}
	return parseVBE(file)
}

func isVBE(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// XXX just guessing
	if b[0] != 0xff || b[1] != 0xfe || b[2] != 0x23 || b[3] != 0 {
		return false
	}
	return true
}

func parseVBE(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.Executable,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 4, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
			}}}}

	return &res, nil
}
