package archive

// STATUS 1%

// https://golang.org/src/compress/bzip2/bzip2.go

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func BZIP2(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isBZIP2(&c.Header) {
		return nil, nil
	}
	return parseBZIP2(c.File, c.ParsedLayout)
}

func isBZIP2(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'B' || b[1] != 'Z' {
		return false
	}
	if b[2] != 'h' {
		// only huffman encoding is used in the format
		return false
	}
	return true
}

func parseBZIP2(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 2, Length: 1, Info: "encoding", Type: parse.Uint8},          // XXX h = huffman
			{Offset: pos + 3, Length: 1, Info: "compression level", Type: parse.ASCII}, // 0=worst, 9=best<
		}}}

	return &pl, nil
}
