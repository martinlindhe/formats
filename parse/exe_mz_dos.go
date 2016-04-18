package parse

import (
	"os"
)

func findCustomDOSHeaders(file *os.File) *Layout {

	pos := int64(0x1c)

	tok, _ := knownLengthASCII(file, pos+2, 9)
	if tok == "PKLITE Co" {

		return &Layout{
			Offset: pos,
			Length: 2 + 52, // XXX
			Info:   "PKLITE header",
			Type:   Group,
			Childs: []Layout{
				Layout{Offset: pos, Length: 1, Info: "minor version", Type: Uint8},
				Layout{Offset: pos + 1, Length: 1, Info: "bit mapped", Type: Uint8},
				Layout{Offset: pos + 2, Length: 52, Info: "identifier", Type: ASCII},
				// XXX bit map:
				// 0-3 - major version
				// 4 - Extra compression
				// 5 - Multi-segment file
			}}
	}

	tok, _ = knownLengthASCII(file, 0x1c, 4)
	if tok == "LZ09" || tok == "LZ91" {

		return &Layout{
			Offset: pos,
			Length: 6, // XXX
			Info:   "LZEXE header",
			Type:   Group,
			Childs: []Layout{
				Layout{Offset: pos, Length: 4, Info: "identifier", Type: ASCII},
				Layout{Offset: pos + 4, Length: 2, Info: "exe version", Type: MajorMinor16le},
			}}

		// XXX
		// "LZ09" = v 0.9
		// "LZ91" = v 0.91
	}

	u32tok, _ := readUint32le(file, 0x1c)
	if u32tok == 0x018A0001 {

		panic("TOPSPEED")
		/*
			return &Layout{
				Offset: offset,
				Length: 6,                 // XXX
				Info:   "TOPSPEED header", // topspeed C 3.0 Crunch
				Type:   Group,
				Childs: []Layout{
					Layout{Offset: offset, Length: 4, Info: "identifier", Type: Uint32le},
					Layout{Offset: offset + 4, Length: 2, Info: "id 2", Type: Uint16le}, // 0x1565 ...
				}}
		*/
	}

	// TODO EXEPACK: http://www.shikadi.net/moddingwiki/Microsoft_EXEPACK

	tlink1, _ := readUint16le(file, pos)
	tlinkId, _ := readUint8(file, pos+2)
	if tlink1 == 0x1 && tlinkId == 0xfb {
		return &Layout{
			Offset: pos,
			Length: 6,
			Info:   "borland TLINK header",
			Type:   Group,
			Childs: []Layout{
				Layout{Offset: pos, Length: 3, Info: "identifier", Type: Bytes},
				Layout{Offset: pos + 3, Length: 1, Info: "version", Type: MajorMinor8},
				Layout{Offset: pos + 4, Length: 2, Info: "???", Type: ASCII}, // always "jr" ?
			}}
	}

	return nil
}
