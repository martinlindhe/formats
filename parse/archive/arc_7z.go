package archive

// STATUS 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func SEVENZIP(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isSEVENZIP(file) {
		return nil, nil
	}
	return parseSEVENZIP(file, pl)
}

func isSEVENZIP(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [6]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != '7' || b[1] != 'z' || b[2] != 0xbc || b[3] != 0xaf ||
		b[4] != 0x27 || b[5] != 0x1c {
		return false
	}

	return true
}

func parseSEVENZIP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 6,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 6, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
