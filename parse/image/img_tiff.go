package image

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func TIFF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isTIFF(file) {
		return nil, nil
	}
	return parseTIFF(file, pl)
}

func isTIFF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// XXX dont know magic numbers just guessing
	if b[0] != 'I' || b[1] != 'I' || b[2] != '*' || b[3] != 0 {
		return false
	}
	return true
}

func parseTIFF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Image
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4,
		Info:   "file header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
