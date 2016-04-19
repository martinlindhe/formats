package windows

// https://github.com/libyal/liblnk/blob/master/documentation/Windows%20Shortcut%20File%20%28LNK%29%20format.asciidoc

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func LNK(c *parse.ParseChecker)(*parse.ParsedLayout, error) {

	if !isLNK(&c.Header) {
		return nil, nil
	}
	return parseLNK(c.File, c.ParsedLayout)
}

func isLNK(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 0x4c || b[1] != 0 || b[2] != 0 || b[3] != 0 {
		return false
	}
	return true
}

func parseLNK(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
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
		}}}

	return &pl, nil
}
