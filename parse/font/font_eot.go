package font

// Embedded OpenType
// https://en.wikipedia.org/wiki/Embedded_OpenType

// STATUS 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func EOT(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isEOT(&hdr) {
		return nil, nil
	}
	return parseEOT(file, pl)
}

func isEOT(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[34] != 0x4c || b[35] != 0x50 {
		return false
	}
	return true
}

func parseEOT(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Font
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 80, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			// XXX Version 0x00010000:
			{Offset: pos, Length: 4, Info: "eot size", Type: parse.Uint32le},
			{Offset: pos + 4, Length: 4, Info: "font data size", Type: parse.Uint32le},
			{Offset: pos + 8, Length: 4, Info: "version", Type: parse.MajorMinor32le},
			{Offset: pos + 12, Length: 4, Info: "flags", Type: parse.Uint32le},
			{Offset: pos + 16, Length: 10, Info: "font panose", Type: parse.Bytes},
			{Offset: pos + 26, Length: 1, Info: "charset", Type: parse.Uint8}, // XXX decode
			{Offset: pos + 27, Length: 1, Info: "italic", Type: parse.Uint8},
			{Offset: pos + 28, Length: 4, Info: "weight", Type: parse.Uint32le},
			{Offset: pos + 32, Length: 2, Info: "fs type", Type: parse.Uint16le},
			{Offset: pos + 34, Length: 2, Info: "magic", Type: parse.Uint16le},
			{Offset: pos + 36, Length: 4, Info: "unicode range 1", Type: parse.Uint32le},
			{Offset: pos + 40, Length: 4, Info: "unicode range 2", Type: parse.Uint32le},
			{Offset: pos + 44, Length: 4, Info: "unicode range 3", Type: parse.Uint32le},
			{Offset: pos + 48, Length: 4, Info: "unicode range 4", Type: parse.Uint32le},
			{Offset: pos + 52, Length: 4, Info: "code page range 1", Type: parse.Uint32le},
			{Offset: pos + 56, Length: 4, Info: "code page range 2", Type: parse.Uint32le},
			{Offset: pos + 60, Length: 4, Info: "reserved 1", Type: parse.Uint32le},
			{Offset: pos + 64, Length: 4, Info: "reserved 2", Type: parse.Uint32le},
			{Offset: pos + 68, Length: 4, Info: "reserved 3", Type: parse.Uint32le},
			{Offset: pos + 72, Length: 4, Info: "reserved 4", Type: parse.Uint32le},
			{Offset: pos + 76, Length: 2, Info: "padding 1", Type: parse.Uint16le},
			{Offset: pos + 78, Length: 2, Info: "family name size", Type: parse.Uint16le},
			// XXX byte 	FamilyName[FamilyNameSize]
			/*
			   unsigned short 	Padding2
			   unsigned short 	StyleNameSize
			   byte 	StyleName[StyleNameSize]
			   unsigned short 	Padding3
			   unsigned short 	VersionNameSize
			   bytes 	VersionName[VersionNameSize]
			   unsigned short 	Padding4
			   unsigned short 	FullNameSize
			   byte 	FullName[FullNameSize]
			   byte 	FontData[FontDataSize]
			*/
		}}}

	return &pl, nil
}
