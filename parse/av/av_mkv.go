package av

// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func MKV(file *os.File) (*parse.ParsedLayout, error) {

	if !isMKV(file) {
		return nil, nil
	}
	return parseMKV(file)
}

func isMKV(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	// XXX what is magic sequence? just guessing
	if b[0] != 0x1a || b[1] != 0x45 || b[2] != 0xdf || b[3] != 0xa3 {
		return false
	}
	return true
}

func parseMKV(file *os.File) (*parse.ParsedLayout, error) {

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
