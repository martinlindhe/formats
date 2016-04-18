package windows

// STATUS: 1%
// https://github.com/libyal/liblnk/blob/master/documentation/Windows%20Shortcut%20File%20%28LNK%29%20format.asciidoc

import (
	"encoding/binary"
	"github.com/martinlindhe/formats/parse"
	"os"
)

func LNK(file *os.File) (*parse.ParsedLayout, error) {

	if !isLNK(file) {
		return nil, nil
	}
	return parseLNK(file)
}

func isLNK(file *os.File) bool {

	file.Seek(0, os.SEEK_SET)
	var b uint32
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b != 0x4c {
		return false
	}
	return true
}

func parseLNK(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.WindowsResource,
		Layout: []parse.Layout{{
			Offset: pos,
			Length: 76, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},

				// XXX type GUID / CLSID  GUID: {00021401-0000-0000-c000-000000000046}
				{Offset: pos + 4, Length: 16, Info: "guid", Type: parse.Bytes},

				{Offset: pos + 20, Length: 4, Info: "data flags", Type: parse.Uint32le},
				{Offset: pos + 24, Length: 4, Info: "file attribute flags", Type: parse.Uint32le},
				{Offset: pos + 28, Length: 8, Info: "created FILETIME", Type: parse.Bytes},  // XXX decode type
				{Offset: pos + 36, Length: 8, Info: "accessed FILETIME", Type: parse.Bytes}, // XXX decode type
				{Offset: pos + 44, Length: 8, Info: "modified FILETIME", Type: parse.Bytes}, // XXX decode type
				{Offset: pos + 52, Length: 4, Info: "file size", Type: parse.Uint32le},
				{Offset: pos + 56, Length: 4, Info: "icon index", Type: parse.Uint32le},
				{Offset: pos + 60, Length: 4, Info: "show window", Type: parse.Uint32le},
				{Offset: pos + 64, Length: 2, Info: "hot key", Type: parse.Uint16le},
				{Offset: pos + 66, Length: 2, Info: "reserved", Type: parse.Uint16le},
				{Offset: pos + 68, Length: 4, Info: "reserved", Type: parse.Uint32le},
				{Offset: pos + 72, Length: 4, Info: "reserved", Type: parse.Uint32le},
			}}}}

	return &res, nil
}
