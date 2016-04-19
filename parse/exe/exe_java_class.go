package exe

// https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func JAVA(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isJAVA(file) {
		return nil, nil
	}
	return parseJAVA(file, pl)
}

func isJAVA(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b uint32
	if err := binary.Read(file, binary.BigEndian, &b); err != nil {
		return false
	}
	if b == 0xcafebabe {
		return true
	}

	return false
}

func parseJAVA(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32be},
		}}}

	return &pl, nil
}
