package exe

import (
	"github.com/martinlindhe/formats/parse"
	"os"
)

// TODO EXEPACK: http://www.shikadi.net/moddingwiki/Microsoft_EXEPACK

func findCustomDOSHeaders(file *os.File) *parse.Layout {

	pos := int64(0x1c)

	tok, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, pos+2, 9)
	if tok == "PKLITE Co" {

		return &parse.Layout{
			Offset: pos,
			Length: 2 + 52, // XXX
			Info:   "PKLITE header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 1, Info: "minor version", Type: parse.Uint8},
				{Offset: pos + 1, Length: 1, Info: "bit mapped", Type: parse.Uint8},
				{Offset: pos + 2, Length: 52, Info: "identifier", Type: parse.ASCII},
				// XXX bit map:
				// 0-3 - major version
				// 4 - Extra compression
				// 5 - Multi-segment file
			}}
	}

	tok, _, _ = parse.ReadZeroTerminatedASCIIUntil(file, 0x1c, 4)
	if tok == "LZ09" || tok == "LZ91" {

		return &parse.Layout{
			Offset: pos,
			Length: 6, // XXX
			Info:   "LZEXE header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 2, Info: "identifier", Type: parse.ASCII},
				{Offset: pos + 2, Length: 2, Info: "version", Type: parse.ASCII},
				{Offset: pos + 4, Length: 2, Info: "exe version", Type: parse.MajorMinor16le},
			}}

		// XXX version
		// "09" = v 0.9
		// "91" = v 0.91
	}

	u32tok, _ := parse.ReadUint32le(file, 0x1c)
	if u32tok == 0x018a0001 {

		panic("TOPSPEED")
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
		return &parse.Layout{
			Offset: pos,
			Length: 6,
			Info:   "borland TLINK header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 3, Info: "identifier", Type: parse.Bytes},
				{Offset: pos + 3, Length: 1, Info: "version", Type: parse.MajorMinor8},
				{Offset: pos + 4, Length: 2, Info: "???", Type: parse.ASCII}, // always "jr" ?
			}}
	}

	return nil
}
