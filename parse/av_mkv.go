package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func MKV(file *os.File) (*ParsedLayout, error) {

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

func parseMKV(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)
	res := ParsedLayout{
		FileKind: AudioVideo,
		Layout: []Layout{{
			Offset: pos,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: ASCII},
			}}}}

	return &res, nil
}
