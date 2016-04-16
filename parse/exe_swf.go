package parse

// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func SWF(file *os.File) (*ParsedLayout, error) {

	if !isSWF(file) {
		return nil, nil
	}
	return parseSWF(file)
}

func isSWF(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [3]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b[0] == 'F' || b[0] == 'C' || b[0] == 'Z' {
		if b[1] == 'W' && b[2] == 'S' {
			return true
		}
	}
	return false
}

func parseSWF(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)

	res := ParsedLayout{
		FileKind: Executable,
		Layout: []Layout{{
			Offset: offset,
			Length: 14, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: offset, Length: 3, Info: "magic", Type: ASCII}, // F = uncompressed, C = zlib compressed, Z = LZMA compressed
				{Offset: offset + 3, Length: 1, Info: "version", Type: Uint8},
				{Offset: offset + 4, Length: 4, Info: "file length", Type: Uint32le},

				// XXX "RECT" type
				// . This field is stored as a RECT structure, meaning that its size may vary according to the number of bits needed to encode the coordinates. The FrameSize RECT always has Xmin and Ymin value of 0; the Xmax and Ymax members define the width and height (see Using bit values).
				{Offset: offset + 8, Length: 2, Info: "frame size", Type: Uint16le},

				{Offset: offset + 10, Length: 2, Info: "frame rate", Type: Uint16le},
				{Offset: offset + 12, Length: 2, Info: "frame count", Type: Uint16le},
			}}}}

	return &res, nil
}
