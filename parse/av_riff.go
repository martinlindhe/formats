package parse

// RIFF format (WAV, AVI)
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func RIFF(file *os.File) (*ParsedLayout, error) {

	if !isRIFF(file) {
		return nil, nil
	}
	return parseRIFF(file)
}

func isRIFF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [4]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] != 'R' || b[1] != 'I' || b[2] != 'F' || b[3] != 'F' {
		return false
	}

	return true
}

func parseRIFF(file *os.File) (*ParsedLayout, error) {

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
