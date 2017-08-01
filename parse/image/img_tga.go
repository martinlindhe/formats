package image

// STATUS: 1%
// https://en.wikipedia.org/wiki/Truevision_TGA

import (
	"encoding/binary"
	"os"

	"github.com/martinlindhe/formats/parse"
)

// TGA parses the Truevision Advanced Raster Graphics Adapter (TARGA) image format
func TGA(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isTGA(c.Header) {
		return nil, nil
	}
	return parseTGA(c.File, c.ParsedLayout)
}

func isTGA(b []byte) bool {
	id := binary.LittleEndian.Uint32(b)
	kind := id & 0xfff7ffff
	switch kind {
	case 0x01010000: // Targa image data - Map
		return true
	case 0x00020000: // Targa image data - RGB
		return true
	case 0x00030000: // Targa image data - Mono
		return true
	default:
		return false
	}
}

func parseTGA(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			// XXX this is le-format only
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32be},
		}}}

	return &pl, nil
}
