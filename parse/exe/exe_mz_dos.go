package exe

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func findCustomDOSHeaders(file *os.File) []parse.Layout {

	pos := int64(0x1c)

	tok, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, pos+2, 9)
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

	tok, _, _ = parse.ReadZeroTerminatedASCIIUntil(file, pos, 4)
	if tok == "LZ09" || tok == "LZ91" {

		return []parse.Layout{{
			Offset: pos,
			Length: 6, // XXX
			Info:   "LZEXE header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 2, Info: "identifier", Type: parse.ASCII},
				{Offset: pos + 2, Length: 2, Info: "version", Type: parse.ASCII},
				{Offset: pos + 4, Length: 2, Info: "exe version", Type: parse.MajorMinor16le},
			}}}

		// XXX version
		// "09" = v 0.9
		// "91" = v 0.91
	}

	u32tok, _ := parse.ReadUint32le(file, pos)
	if u32tok == 0x018a0001 {

		fmt.Println("info: TOPSPEED sample plz!")

		return []parse.Layout{{
			Offset: pos,
			Length: 4, // XXX
			Info:   "TOPSPEED header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 4, Info: "identifier", Type: parse.Uint32le},
			}}}
		/*
			return &Layout{
				Offset: offset,
				Length: 6,                 // XXX
				Info:   "TOPSPEED header", // topspeed C 3.0 Crunch
				Type:   Group,
				Childs: []Layout{
					Layout{Offset: offset, Length: 4, Info: "identifier", Type: parse.Uint32le},
					Layout{Offset: offset + 4, Length: 2, Info: "id 2", Type: parse.Uint16le}, // 0x1565 ...
				}}
		*/
	}

	tlink1, _ := parse.ReadUint16le(file, pos)
	tlinkId, _ := parse.ReadUint8(file, pos+2)
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
	headerSizeInParagraphs, _ := parse.ReadUint16le(file, 8)
	cs, _ := parse.ReadUint16le(file, 22)
	ip, _ := parse.ReadUint16le(file, 20)
	exePackOffset := (int64(headerSizeInParagraphs) * 16) // XXX hack
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
