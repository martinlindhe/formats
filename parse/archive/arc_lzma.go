package archive

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func LZMA(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isLZMA(file) {
		return nil, nil
	}
	return parseLZMA(file, pl)
}

func isLZMA(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [6]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// XXX not proper magic , need other check
	if b[0] != 0x5d || b[1] != 0x00 {
		return false
	}

	return true
}

func parseLZMA(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 13, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			// XXX unsure of this stuff
			{Offset: pos, Length: 1, Info: "properties", Type: parse.Uint8},
			{Offset: pos + 1, Length: 4, Info: "dict cap", Type: parse.Uint32le},
			{Offset: pos + 5, Length: 8, Info: "uncompressed size", Type: parse.Uint64le},
		}}}

	return &pl, nil
}
