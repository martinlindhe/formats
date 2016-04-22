package image

// Targa image
// https://en.wikipedia.org/wiki/Truevision_TGA

// STATUS: 1%

import (
	"encoding/binary"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func TGA(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isTGA(c.Header) {
		return nil, nil
	}
	return parseTGA(c.File, c.ParsedLayout)
}

func isTGA(b []byte) bool {

	/* XXX
	   # at 2, byte ImgType must be 1, 2, 3, 9, 10 or 11
	   # at 1, byte CoMapType must be 1 if ImgType is 1 or 9, 0 otherwise
	   # at 3, leshort Index is 0 for povray, ppmtotga and xv outputs
	*/

	val := binary.LittleEndian.Uint32(b)
	chk := val & 0xfff7ffff

	if chk == 0x01010000 {
		// Targa image data - Map
		// >2	byte&8			8		- RLE
		// >12	leshort			>0		%hd x
		// >14	leshort			>0		%hd
		return true
	}
	if chk == 0x00020000 {
		// Targa image data - RGB
		// >2	byte&8			8		- RLE
		// >12	leshort			>0		%hd x
		// >14	leshort			>0		%hd
		return true
	}
	if chk == 0x00030000 {
		// Targa image data - Mono
		// >2	byte&8			8		- RLE
		// >12	leshort			>0		%hd x
		// >14	leshort			>0		%hd
		return true
	}
	return false
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
