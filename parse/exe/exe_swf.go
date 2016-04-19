package exe

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func SWF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isSWF(&hdr) {
		return nil, nil
	}
	return parseSWF(file, pl)
}

func isSWF(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] == 'F' || b[0] == 'C' || b[0] == 'Z' {
		if b[1] == 'W' && b[2] == 'S' {
			return true
		}
	}
	return false
}

func parseSWF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Executable
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 14, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 3, Info: "magic", Type: parse.ASCII}, // XXX F = uncompressed, C = zlib compressed, Z = LZMA compressed
			{Offset: pos + 3, Length: 1, Info: "version", Type: parse.Uint8},
			{Offset: pos + 4, Length: 4, Info: "file length", Type: parse.Uint32le},

			// XXX "RECT" type
			// . This field is stored as a RECT structure, meaning that its size may vary according to the number of bits needed to encode the coordinates. The FrameSize RECT always has Xmin and Ymin value of 0; the Xmax and Ymax members define the width and height (see Using bit values).
			{Offset: pos + 8, Length: 2, Info: "frame size", Type: parse.Uint16le},
			{Offset: pos + 10, Length: 2, Info: "frame rate", Type: parse.Uint16le},
			{Offset: pos + 12, Length: 2, Info: "frame count", Type: parse.Uint16le},
		}}}

	return &pl, nil
}
