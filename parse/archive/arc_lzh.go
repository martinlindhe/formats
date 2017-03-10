package archive

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// STATUS: 1%

// LZH parses the LZH format, created by the LHArc/LHA archiver
func LZH(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isLZH(c.Header) {
		return nil, nil
	}
	return parseLZH(c.File, c.ParsedLayout)
}

func isLZH(b []byte) bool {

	if b[2] != '-' || b[3] != 'l' {
		return false
	}
	if b[4] == 'h' || b[4] == 'z' {
		return true
	}
	return false
}

func parseLZH(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	methodName, _ := parse.ReadToMap(file, parse.Uint8, pos+5, lhaCompressionMethods)

	pl.FileKind = parse.Archive
	pl.MimeType = "application/x-lzh-compressed"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 21, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 1, Info: "header size", Type: parse.Uint8},
			{Offset: pos + 1, Length: 1, Info: "checksum", Type: parse.Uint8},
			{Offset: pos + 2, Length: 3, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 5, Length: 1, Info: "compression = " + methodName, Type: parse.Uint8},
			{Offset: pos + 6, Length: 1, Info: "reserved", Type: parse.Uint8}, // "-"
			{Offset: pos + 7, Length: 4, Info: "compressed size", Type: parse.Uint32le},
			{Offset: pos + 11, Length: 4, Info: "uncompressed size", Type: parse.Uint32le},
			{Offset: pos + 15, Length: 4, Info: "original file date/time", Type: parse.DOSDateTime},
			{Offset: pos + 19, Length: 2, Info: "file attribute", Type: parse.Uint16le}, // XXX flags

			// uncertain:
			//{Offset: pos + 21, Length: 1, Info: "length of filename", Type: parse.Uint8},        // ???
			//{Offset: pos + 22, Length: 2, Info: "filename / path", Type: parse.Bytes},           // ???
			//{Offset: pos + 24, Length: 2, Info: "crc16 of original file", Type: parse.Uint16le}, // ???
			// XXX filename follows
		}}}

	return &pl, nil
}

var (
	lhaCompressionMethods = map[byte]string{
		0x30: "none",
		0x31: "LZW, 4K buffer, Huffman for upper 6 bits of position",
		0x32: "unknown",
		0x33: "unknown",
		0x34: "LZW, Arithmetic Encoding",
		0x35: "LZW, Arithmetic Encoding",
		0x5c: "LHa 2.x archive?",
		0x64: "LHa 2.x archive?",
		0x73: "LHa 2.x archive?",
	}
)
