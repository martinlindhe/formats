package exe

// Dalvik Executable (android java code)
// STATUS: 1%

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func DEX(file *os.File) (*parse.ParsedLayout, error) {

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

func parseDEX(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.Executable,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 112, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
				{Offset: pos + 4, Length: 4, Info: "version", Type: parse.ASCII},
				{Offset: pos + 8, Length: 4, Info: "checksum", Type: parse.Uint32le},
				{Offset: pos + 12, Length: 20, Info: "sha1", Type: parse.Bytes},
				{Offset: pos + 32, Length: 4, Info: "file size", Type: parse.Uint32le},
				{Offset: pos + 36, Length: 4, Info: "header size", Type: parse.Uint32le},
				{Offset: pos + 40, Length: 4, Info: "endian tag", Type: parse.Uint32le},
				{Offset: pos + 44, Length: 4, Info: "link size", Type: parse.Uint32le},
				{Offset: pos + 48, Length: 4, Info: "link offset", Type: parse.Uint32le},
				{Offset: pos + 52, Length: 4, Info: "map offset", Type: parse.Uint32le},
				{Offset: pos + 56, Length: 4, Info: "string ids size", Type: parse.Uint32le},
				{Offset: pos + 60, Length: 4, Info: "string ids offset", Type: parse.Uint32le},
				{Offset: pos + 64, Length: 4, Info: "type ids size", Type: parse.Uint32le},
				{Offset: pos + 68, Length: 4, Info: "type ids offset", Type: parse.Uint32le},
				{Offset: pos + 72, Length: 4, Info: "proto ids size", Type: parse.Uint32le},
				{Offset: pos + 76, Length: 4, Info: "proto ids offset", Type: parse.Uint32le},
				{Offset: pos + 80, Length: 4, Info: "field ids size", Type: parse.Uint32le},
				{Offset: pos + 84, Length: 4, Info: "field ids offset", Type: parse.Uint32le},
				{Offset: pos + 88, Length: 4, Info: "method ids size", Type: parse.Uint32le},
				{Offset: pos + 92, Length: 4, Info: "method ids offset", Type: parse.Uint32le},
				{Offset: pos + 96, Length: 4, Info: "class definition size", Type: parse.Uint32le},
				{Offset: pos + 100, Length: 4, Info: "class definition offset", Type: parse.Uint32le},
				{Offset: pos + 104, Length: 4, Info: "data size", Type: parse.Uint32le},
				{Offset: pos + 108, Length: 4, Info: "data offset", Type: parse.Uint32le},
			}}}}

	return &res, nil
}
