package parse

// STATUS 1%

import (
	"encoding/binary"
	"os"
)

func EOT(file *os.File) (*ParsedLayout, error) {

	if !isEOT(file) {
		return nil, nil
	}
	return parseEOT(file)
}

func isEOT(file *os.File) bool {

	file.Seek(34, os.SEEK_SET)
	var b uint16
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return false
	}
	if b != 0x504c {
		return false
	}
	return true
}

func parseEOT(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)
	res := ParsedLayout{
		FileKind: Font,
		Layout: []Layout{{
			Offset: pos,
			Length: 80, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				// XXX Version 0x00010000:
				{Offset: pos, Length: 4, Info: "eot size", Type: Uint32le},
				{Offset: pos + 4, Length: 4, Info: "font data size", Type: Uint32le},
				{Offset: pos + 8, Length: 4, Info: "version", Type: MajorMinor32le},
				{Offset: pos + 12, Length: 4, Info: "flags", Type: Uint32le},
				{Offset: pos + 16, Length: 10, Info: "font panose", Type: Bytes},
				{Offset: pos + 26, Length: 1, Info: "charset", Type: Uint8}, // XXX decode
				{Offset: pos + 27, Length: 1, Info: "italic", Type: Uint8},
				{Offset: pos + 28, Length: 4, Info: "weight", Type: Uint32le},
				{Offset: pos + 32, Length: 2, Info: "fs type", Type: Uint16le},
				{Offset: pos + 34, Length: 2, Info: "magic", Type: Uint16le},

				{Offset: pos + 36, Length: 4, Info: "unicode range 1", Type: Uint32le},
				{Offset: pos + 40, Length: 4, Info: "unicode range 2", Type: Uint32le},
				{Offset: pos + 44, Length: 4, Info: "unicode range 3", Type: Uint32le},
				{Offset: pos + 48, Length: 4, Info: "unicode range 4", Type: Uint32le},

				{Offset: pos + 52, Length: 4, Info: "code page range 1", Type: Uint32le},
				{Offset: pos + 56, Length: 4, Info: "code page range 2", Type: Uint32le},

				{Offset: pos + 60, Length: 4, Info: "reserved 1", Type: Uint32le},
				{Offset: pos + 64, Length: 4, Info: "reserved 2", Type: Uint32le},
				{Offset: pos + 68, Length: 4, Info: "reserved 3", Type: Uint32le},
				{Offset: pos + 72, Length: 4, Info: "reserved 4", Type: Uint32le},

				{Offset: pos + 76, Length: 2, Info: "padding 1", Type: Uint16le},

				{Offset: pos + 78, Length: 2, Info: "family name size", Type: Uint16le},
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
			}}}}

	return &res, nil
}
