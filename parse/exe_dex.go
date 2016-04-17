package parse

// Dalvik Executable (android java code)
// STATUS: 1%

import (
	"encoding/binary"
	"os"
)

func DEX(file *os.File) (*ParsedLayout, error) {

	if !isDEX(file) {
		return nil, nil
	}
	return parseDEX(file)
}

func isDEX(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b [8]byte
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}

	if b[0] == 'd' && b[1] == 'e' && b[2] == 'x' && b[3] == '\n' &&
		b[4] == '0' && b[5] == '3' && b[6] == '5' && b[7] == 0 {
		return true
	}

	return false
}

func parseDEX(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)

	res := ParsedLayout{
		FileKind: Executable,
		Layout: []Layout{{
			Offset: offset,
			Length: 112, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: offset, Length: 4, Info: "magic", Type: ASCII},
				{Offset: offset + 4, Length: 4, Info: "version", Type: ASCII},
				{Offset: offset + 8, Length: 4, Info: "checksum", Type: Uint32le},
				{Offset: offset + 12, Length: 20, Info: "sha1", Type: Bytes},
				{Offset: offset + 32, Length: 4, Info: "file size", Type: Uint32le},
				{Offset: offset + 36, Length: 4, Info: "header size", Type: Uint32le},
				{Offset: offset + 40, Length: 4, Info: "endian tag", Type: Uint32le},
				{Offset: offset + 44, Length: 4, Info: "link size", Type: Uint32le},
				{Offset: offset + 48, Length: 4, Info: "link offset", Type: Uint32le},
				{Offset: offset + 52, Length: 4, Info: "map offset", Type: Uint32le},
				{Offset: offset + 56, Length: 4, Info: "string ids size", Type: Uint32le},
				{Offset: offset + 60, Length: 4, Info: "string ids offset", Type: Uint32le},
				{Offset: offset + 64, Length: 4, Info: "type ids size", Type: Uint32le},
				{Offset: offset + 68, Length: 4, Info: "type ids offset", Type: Uint32le},
				{Offset: offset + 72, Length: 4, Info: "proto ids size", Type: Uint32le},
				{Offset: offset + 76, Length: 4, Info: "proto ids offset", Type: Uint32le},
				{Offset: offset + 80, Length: 4, Info: "field ids size", Type: Uint32le},
				{Offset: offset + 84, Length: 4, Info: "field ids offset", Type: Uint32le},
				{Offset: offset + 88, Length: 4, Info: "method ids size", Type: Uint32le},
				{Offset: offset + 92, Length: 4, Info: "method ids offset", Type: Uint32le},
				{Offset: offset + 96, Length: 4, Info: "class definition size", Type: Uint32le},
				{Offset: offset + 100, Length: 4, Info: "class definition offset", Type: Uint32le},
				{Offset: offset + 104, Length: 4, Info: "data size", Type: Uint32le},
				{Offset: offset + 108, Length: 4, Info: "data offset", Type: Uint32le},
			}}}}

	return &res, nil
}
