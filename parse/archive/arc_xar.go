package archive

// Xar Archive
// Extensions: .xar .pkg
// https://en.wikipedia.org/wiki/Xar_%28archiver%29

// STATUS: 3%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func XAR(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isXAR(c.Header) {
		return nil, nil
	}
	return parseXAR(c.File, c.ParsedLayout)
}

func isXAR(b []byte) bool {

	if b[0] != 'x' || b[1] != 'a' || b[2] != 'r' || b[3] != '!' {
		return false
	}
	return true
}

func parseXAR(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)

	hdrLen, _ := parse.ReadUint16be(file, pos+4)
	tocLen, _ := parse.ReadUint64be(file, pos+8)

	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 28, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 4, Length: 2, Info: "header size", Type: parse.Uint16be},
			{Offset: pos + 6, Length: 2, Info: "version", Type: parse.MinorMajor16le}, // XXX be.. the version types are a mess
			{Offset: pos + 8, Length: 8, Info: "toc length compressed", Type: parse.Uint64be},
			{Offset: pos + 16, Length: 8, Info: "toc length uncompressed", Type: parse.Uint64be},
			{Offset: pos + 24, Length: 4, Info: "checksum", Type: parse.Uint32be},
		}}}

	if hdrLen != 0x1c {
		fmt.Println("warning: xar header len expected ", 0x1c, "found", hdrLen)
	}

	pos += int64(hdrLen)

	if tocLen > 0 {
		pl.Layout = append(pl.Layout, parse.Layout{
			Offset: pos,
			Length: int64(tocLen),
			Info:   "toc",
			Type:   parse.Group,

			Childs: []parse.Layout{
				{Offset: pos, Length: int64(tocLen), Info: "toc data", Type: parse.Bytes},
			}})
	}

	return &pl, nil
}
