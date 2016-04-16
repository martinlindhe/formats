package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func MIDI(file *os.File) (*ParsedLayout, error) {

	if !isMIDI(file) {
		return nil, nil
	}
	return parseMIDI(file)
}

func isMIDI(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] != 'M' || b[1] != 'T' || b[2] != 'h' || b[3] != 'd' {
		return false
	}
	return true
}

func parseMIDI(file *os.File) (*ParsedLayout, error) {

	res := ParsedLayout{
		FileKind: AudioVideo,
		Layout: []Layout{{
			Offset: 0,
			Length: 4, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: 0, Length: 4, Info: "magic", Type: ASCII},
			}}}}

	return &res, nil
}
