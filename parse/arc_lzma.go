package parse

import (
	"encoding/binary"
	"os"
)

func LZMA(file *os.File) (*ParsedLayout, error) {

	if !isLZMA(file) {
		return nil, nil
	}
	return parseLZMA(file)
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

func parseLZMA(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)

	res := ParsedLayout{
		FileKind: Archive,
		Layout: []Layout{{
			Offset: offset,
			Length: 13, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				// XXX unsure of this stuff
				{Offset: offset, Length: 1, Info: "properties", Type: Uint8},
				{Offset: offset + 1, Length: 4, Info: "dict cap", Type: Uint32le},
				{Offset: offset + 5, Length: 8, Info: "uncompressed size", Type: Uint64le},
			}}}}

	return &res, nil
}
