package parse

// STATUS: 1%
// https://github.com/libyal/liblnk/blob/master/documentation/Windows%20Shortcut%20File%20%28LNK%29%20format.asciidoc

import (
	"encoding/binary"
	"os"
)

func LNK(file *os.File) (*ParsedLayout, error) {

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

func parseLNK(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)

	res := ParsedLayout{
		FileKind: WindowsResource,
		Layout: []Layout{{
			Offset: offset,
			Length: 76, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: offset, Length: 4, Info: "magic", Type: Uint32le},

				// XXX type GUID / CLSID  GUID: {00021401-0000-0000-c000-000000000046}
				{Offset: offset + 4, Length: 16, Info: "guid", Type: Bytes},

				{Offset: offset + 20, Length: 4, Info: "data flags", Type: Uint32le},
				{Offset: offset + 24, Length: 4, Info: "file attribute flags", Type: Uint32le},
				{Offset: offset + 28, Length: 8, Info: "created FILETIME", Type: Bytes},  // XXX decode type
				{Offset: offset + 36, Length: 8, Info: "accessed FILETIME", Type: Bytes}, // XXX decode type
				{Offset: offset + 44, Length: 8, Info: "modified FILETIME", Type: Bytes}, // XXX decode type
				{Offset: offset + 52, Length: 4, Info: "file size", Type: Uint32le},
				{Offset: offset + 56, Length: 4, Info: "icon index", Type: Uint32le},
				{Offset: offset + 60, Length: 4, Info: "show window", Type: Uint32le},
				{Offset: offset + 64, Length: 2, Info: "hot key", Type: Uint16le},
				{Offset: offset + 66, Length: 2, Info: "reserved", Type: Uint16le},
				{Offset: offset + 68, Length: 4, Info: "reserved", Type: Uint32le},
				{Offset: offset + 72, Length: 4, Info: "reserved", Type: Uint32le},
			}}}}

	return &res, nil
}
