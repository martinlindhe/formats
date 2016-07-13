package windows

// STATUS: 1%
// Extensions: .lnk
// http://lifeinhex.com/analyzing-malicious-lnk-file/
// https://github.com/libyal/liblnk/blob/master/documentation/Windows%20Shortcut%20File%20%28LNK%29%20format.asciidoc
// https://en.wikipedia.org/wiki/File_shortcut#Microsoft_Windows

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// LNK parses the Windows shortcut format
func LNK(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isLNK(c.Header) {
		return nil, nil
	}
	return parseLNK(c.File, c.ParsedLayout)
}

func isLNK(b []byte) bool {

	if b[0] != 0x4c || b[1] != 0 || b[2] != 0 || b[3] != 0 {
		return false
	}
	return true
}

func parseLNK(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.MimeType = "application/x-ms-shortcut"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 76, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},

			// XXX type GUID / CLSID  GUID: {00021401-0000-0000-c000-000000000046}
			{Offset: pos + 4, Length: 16, Info: "guid", Type: parse.Bytes},

			{Offset: pos + 20, Length: 4, Info: "data flags", Type: parse.Uint32le, Masks: []parse.Mask{
				// XXX flags
				{Low: 0, Length: 18, Info: "unused"},
				{Low: 18, Length: 1, Info: "encrypted"},
				{Low: 19, Length: 1, Info: "not content indexed"},
				{Low: 20, Length: 1, Info: "offline"},
				{Low: 21, Length: 1, Info: "compressed"},
				{Low: 22, Length: 1, Info: "reparse point"},
				{Low: 23, Length: 1, Info: "sparse file"},
				{Low: 24, Length: 1, Info: "temporary"},
				{Low: 25, Length: 1, Info: "normal"},
				{Low: 26, Length: 1, Info: "archive"},
				{Low: 27, Length: 1, Info: "directory"},
				{Low: 28, Length: 1, Info: "reserved"},
				{Low: 29, Length: 1, Info: "system"},
				{Low: 30, Length: 1, Info: "hidden"},
				{Low: 31, Length: 1, Info: "readonly"},
			}},
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
