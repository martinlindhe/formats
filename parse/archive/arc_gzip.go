package archive

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func GZIP(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isGZIP(file) {
		return nil, nil
	}
	return parseGZIP(file, pl)
}

func isGZIP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [2]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 0x1f || b[1] != 0x8b {
		return false
	}
	return true
}

func parseGZIP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 2,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.Uint16le}, // XXX le/be ?
		}}}

	return &pl, nil
}
