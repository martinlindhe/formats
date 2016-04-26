package exe

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

var (
	lzexeVersions = map[uint16]string{
		0x3130: "0.9",  // "09"
		0x3139: "0.91", // "91"
	}
	lzexe090 = []byte{
		0x06, 0x0e, 0x1f, 0x8b, 0x0e, 0x0c, 0x00, 0x8b, 0xf1, 0x4e, 0x89, 0xf7,
		0x8c, 0xdb, 0x03, 0x1e, 0x0a, 0x00, 0x8e, 0xc3, 0xb4, 0x00, 0x31, 0xed,
		0xfd, 0xac, 0x01, 0xc5, 0xaa, 0xe2, 0xfa, 0x8b, 0x16, 0x0e, 0x00, 0x8a,
		0xc2, 0x29, 0xc5, 0x8a, 0xc6, 0x29, 0xc5, 0x39, 0xd5, 0x74, 0x0c, 0xba,
		0x91, 0x01, 0xb4, 0x09, 0xcd, 0x21, 0xb8, 0xff, 0x4c, 0xcd, 0x21, 0x53,
		0xb8, 0x53, 0x00, 0x50, 0xcb, 0x2e, 0x8b, 0x2e, 0x08, 0x00, 0x8c, 0xda,
		0x89, 0xe8, 0x3d, 0x00, 0x10, 0x76, 0x03, 0xb8, 0x00, 0x10, 0x29, 0xc5,
		0x29, 0xc2, 0x29, 0xc3, 0x8e, 0xda, 0x8e, 0xc3, 0xb1, 0x03, 0xd3, 0xe0,
		0x89, 0xc1, 0xd1, 0xe0, 0x48, 0x48, 0x8b, 0xf0, 0x8b, 0xf8, 0xf3, 0xa5,
		0x09, 0xed, 0x75, 0xd8, 0xfc, 0x8e, 0xc2, 0x8e, 0xdb, 0x31, 0xf6, 0x31,
		0xff, 0xba, 0x10, 0x00, 0xad, 0x89, 0xc5, 0xd1, 0xed, 0x4a, 0x75, 0x05,
		0xad, 0x89, 0xc5, 0xb2, 0x10, 0x73, 0x03, 0xa4, 0xeb, 0xf1, 0x31, 0xc9,
		0xd1, 0xed, 0x4a, 0x75, 0x05, 0xad, 0x89, 0xc5, 0xb2, 0x10, 0x72, 0x22,
		0xd1, 0xed, 0x4a, 0x75, 0x05, 0xad, 0x89, 0xc5, 0xb2, 0x10, 0xd1, 0xd1,
		0xd1, 0xed, 0x4a, 0x75, 0x05, 0xad, 0x89, 0xc5, 0xb2, 0x10, 0xd1, 0xd1,
		0x41, 0x41, 0xac, 0xb7, 0xff, 0x8a, 0xd8, 0xe9, 0x13, 0x00, 0xad, 0x8b,
		0xd8, 0xb1, 0x03, 0xd2, 0xef, 0x80, 0xcf, 0xe0, 0x80, 0xe4, 0x07, 0x74,
		0x0c, 0x88, 0xe1, 0x41, 0x41, 0x26, 0x8a, 0x01, 0xaa, 0xe2, 0xfa, 0xeb,
		0xa6, 0xac, 0x08, 0xc0, 0x74, 0x40, 0x3c, 0x01, 0x74, 0x05, 0x88, 0xc1,
		0x41, 0xeb, 0xea, 0x89}
	lzexe091 = []byte{
		0x06, 0x0e, 0x1f, 0x8b, 0x0e, 0x0c, 0x00, 0x8b, 0xf1, 0x4e, 0x89, 0xf7,
		0x8c, 0xdb, 0x03, 0x1e, 0x0a, 0x00, 0x8e, 0xc3, 0xfd, 0xf3, 0xa4, 0x53,
		0xb8, 0x2b, 0x00, 0x50, 0xcb, 0x2e, 0x8b, 0x2e, 0x08, 0x00, 0x8c, 0xda,
		0x89, 0xe8, 0x3d, 0x00, 0x10, 0x76, 0x03, 0xb8, 0x00, 0x10, 0x29, 0xc5,
		0x29, 0xc2, 0x29, 0xc3, 0x8e, 0xda, 0x8e, 0xc3, 0xb1, 0x03, 0xd3, 0xe0,
		0x89, 0xc1, 0xd1, 0xe0, 0x48, 0x48, 0x8b, 0xf0, 0x8b, 0xf8, 0xf3, 0xa5,
		0x09, 0xed, 0x75, 0xd8, 0xfc, 0x8e, 0xc2, 0x8e, 0xdb, 0x31, 0xf6, 0x31,
		0xff, 0xba, 0x10, 0x00, 0xad, 0x89, 0xc5, 0xd1, 0xed, 0x4a, 0x75, 0x05,
		0xad, 0x89, 0xc5, 0xb2, 0x10, 0x73, 0x03, 0xa4, 0xeb, 0xf1, 0x31, 0xc9,
		0xd1, 0xed, 0x4a, 0x75, 0x05, 0xad, 0x89, 0xc5, 0xb2, 0x10, 0x72, 0x22,
		0xd1, 0xed, 0x4a, 0x75, 0x05, 0xad, 0x89, 0xc5, 0xb2, 0x10, 0xd1, 0xd1,
		0xd1, 0xed, 0x4a, 0x75, 0x05, 0xad, 0x89, 0xc5, 0xb2, 0x10, 0xd1, 0xd1,
		0x41, 0x41, 0xac, 0xb7, 0xff, 0x8a, 0xd8, 0xe9, 0x13, 0x00, 0xad, 0x8b,
		0xd8, 0xb1, 0x03, 0xd2, 0xef, 0x80, 0xcf, 0xe0, 0x80, 0xe4, 0x07, 0x74,
		0x0c, 0x88, 0xe1, 0x41, 0x41, 0x26, 0x8a, 0x01, 0xaa, 0xe2, 0xfa, 0xeb,
		0xa6, 0xac, 0x08, 0xc0, 0x74, 0x34, 0x3c, 0x01, 0x74, 0x05, 0x88, 0xc1,
		0x41, 0xeb, 0xea, 0x89, 0xfb, 0x83, 0xe7, 0x0f, 0x81, 0xc7, 0x00, 0x20,
		0xb1, 0x04, 0xd3, 0xeb, 0x8c, 0xc0, 0x01, 0xd8, 0x2d, 0x00, 0x02, 0x8e,
		0xc0, 0x89, 0xf3, 0x83, 0xe6, 0x0f, 0xd3, 0xeb, 0x8c, 0xd8, 0x01, 0xd8,
		0x8e, 0xd8, 0xe9, 0x72}
)

func findCustomDOSHeaders(file *os.File, b []byte) []parse.Layout {

	pos := int64(0)
	headerSizeInParagraphs, _ := parse.ReadUint16le(file, pos+8)
	ip, _ := parse.ReadUint16le(file, pos+20)
	cs, _ := parse.ReadUint16le(file, pos+22)

	pos = 0x1c
	tok := string(b[pos+2 : pos+9])
	if tok == "PKLITE Co" {

		return []parse.Layout{{
			Offset: pos,
			Length: 2 + 52, // XXX
			Info:   "PKLITE header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 1, Info: "minor version", Type: parse.Uint8},
				{Offset: pos + 1, Length: 1, Info: "bit mapped", Type: parse.Uint8, Masks: []parse.Mask{
					{Low: 0, Length: 4, Info: "major version"},
					{Low: 4, Length: 1, Info: "extra compression"},
					{Low: 5, Length: 1, Info: "multi-segment file"},
					{Low: 6, Length: 2, Info: "unknown"},
				}},
				{Offset: pos + 2, Length: 52, Info: "identifier", Type: parse.ASCII},
			}}}
	}

	// point to dos entry point
	lzStart := int64(((headerSizeInParagraphs + cs) * 16) + ip)

	if parse.HasSignatureInHeader(b, lzStart, lzexe090) ||
		parse.HasSignatureInHeader(b, lzStart, lzexe091) {

		// NOTE: some users of lzexe compression changed the bytes at 0x1c to
		// avoid detection, so instead we match the decompression code

		version, _ := parse.ReadToMap(file, parse.Uint16le, pos+2, lzexeVersions)

		return []parse.Layout{
			{
				Offset: pos,
				Length: 20,
				Info:   "LZEXE header",
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: pos, Length: 2, Info: "identifier", Type: parse.ASCII},
					{Offset: pos + 2, Length: 2, Info: "version = " + version, Type: parse.Uint16le},
					{Offset: pos + 4, Length: 4, Info: "real cs:ip", Type: parse.DOSOffsetSegment},
					{Offset: pos + 8, Length: 4, Info: "real ss:sp", Type: parse.DOSOffsetSegment},
					{Offset: pos + 12, Length: 2, Info: "compressed load module size", Type: parse.Uint16le},
					{Offset: pos + 14, Length: 2, Info: "increase load module size", Type: parse.Uint16le},
					{Offset: pos + 16, Length: 2, Info: "uncompressor size", Type: parse.Uint16le}, // XXX ?
					{Offset: pos + 18, Length: 2, Info: "checksum", Type: parse.Uint16le},
				}},
			{
				Offset: lzStart,
				Length: 232,
				Info:   "lzexe uncompressor",
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: lzStart, Length: 232, Info: "code", Type: parse.Bytes},
				}}}
	}

	u32tok := binary.LittleEndian.Uint32(b[pos:])
	if u32tok == 0x018a0001 {

		fmt.Println("info: exe-dos TOPSPEED compressed sample plz!")

		return []parse.Layout{{
			Offset: pos,
			Length: 6,                 // XXX
			Info:   "TOPSPEED header", // topspeed C 3.0 Crunch
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "identifier", Type: parse.Uint32le},
				{Offset: pos + 4, Length: 2, Info: "id 2", Type: parse.Uint16le}, // 0x1565 ...
			}}}
	}

	tlink1 := binary.LittleEndian.Uint16(b[pos:])
	tlinkId := b[pos+2]
	if tlink1 == 0x1 && tlinkId == 0xfb {
		return []parse.Layout{{
			Offset: pos,
			Length: 6,
			Info:   "borland TLINK header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 3, Info: "identifier", Type: parse.Bytes},
				{Offset: pos + 3, Length: 1, Info: "version", Type: parse.MajorMinor8},
				{Offset: pos + 4, Length: 2, Info: "???", Type: parse.ASCII}, // always "jr" ?
			}}}
	}

	// EXEPACK
	// http://www.shikadi.net/moddingwiki/Microsoft_EXEPACK
	exePackOffset := (int64(headerSizeInParagraphs) * 16)
	exePackSize := int64(cs)*16 + int64(ip)

	exepackCheck := parse.ReadBytesFrom(file, exePackOffset+exePackSize-2, 2)
	if exepackCheck[0] == 'R' && exepackCheck[1] == 'B' {

		pos = exePackOffset
		return []parse.Layout{
			{
				Offset: pos,
				Length: exePackSize - 18,
				Info:   "EXEPACK packed exe",
				Type:   parse.Group,
				Childs: []parse.Layout{
					{Offset: pos, Length: exePackSize, Info: "packed exe", Type: parse.Bytes},
				}},
			{
				Offset: pos + exePackSize - 18,
				Length: 18 + 0x105 + 7 + 0x16,
				Info:   "EXEPACK vars",
				Childs: []parse.Layout{
					{Offset: pos + exePackSize - 18, Length: 2, Info: "real IP", Type: parse.Uint16le},
					{Offset: pos + exePackSize - 18 + 2, Length: 2, Info: "real CS", Type: parse.Uint16le},
					{Offset: pos + exePackSize - 18 + 4, Length: 2, Info: "mem start", Type: parse.Uint16le},
					{Offset: pos + exePackSize - 18 + 6, Length: 2, Info: "exepack size", Type: parse.Uint16le},
					{Offset: pos + exePackSize - 18 + 8, Length: 2, Info: "real SP", Type: parse.Uint16le},
					{Offset: pos + exePackSize - 18 + 10, Length: 2, Info: "real SS", Type: parse.Uint16le},
					{Offset: pos + exePackSize - 18 + 12, Length: 2, Info: "dest len", Type: parse.Uint16le},
					{Offset: pos + exePackSize - 18 + 14, Length: 2, Info: "skip len", Type: parse.Uint16le},
					{Offset: pos + exePackSize - 18 + 16, Length: 2, Info: "signature", Type: parse.Uint16le}, // XXX "RB"
					{Offset: pos + exePackSize - 18 + 18, Length: 0x105 + 7, Info: "unpacker code", Type: parse.Bytes},
					{Offset: pos + exePackSize - 18 + 18 + 0x105 + 7, Length: 0x16, Info: "magic", Type: parse.ASCII}, // XXX
				}}}
	}

	return nil
}
