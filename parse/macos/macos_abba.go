package macos

// ??? found in Mac OS X, /System/Library/PrivateFrameworks/CoreNLP.framework/Versions/A/Resources/de.cache
// Extensions: .cache

// STATUS: 1%

import (
	"encoding/binary"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func ABBA(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isABBA(c.Header) {
		return nil, nil
	}
	return parseABBA(c.File, c.ParsedLayout)
}

func isABBA(b []byte) bool {

	val := binary.LittleEndian.Uint32(b[:])
	return val == 0xabba
}

func parseABBA(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.MacOSResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
