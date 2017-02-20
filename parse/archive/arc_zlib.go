package archive

// Zlib data stream
// http://stackoverflow.com/questions/9050260/what-does-a-zlib-header-look-like
// https://tools.ietf.org/html/rfc1950

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// ZLib parses the zlib format
func ZLib(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isZLib(c.Header) {
		return nil, nil
	}
	return parseZLib(c.File, c.ParsedLayout)
}

func isZLib(b []byte) bool {

	// XXX only matches zlib streams without dictionary.. this dont always work
	if b[0] != 0x78 {
		return false
	}
	if b[1] != 0x01 && b[1] != 0x9c && b[1] != 0xda {
		// compression level
		return false
	}
	return true
}

func parseZLib(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.MimeType = "application/zlib"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 2, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 2, Info: "magic", Type: parse.Uint16le},
		}}}

	return &pl, nil
}
