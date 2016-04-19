package image

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func TIFF(c *parse.ParseChecker)(*parse.ParsedLayout, error) {

	if !isTIFF(&c.Header) {
		return nil, nil
	}
	return parseTIFF(c.File, c.ParsedLayout)
}

func isTIFF(hdr *[0xffff]byte) bool {

	b := *hdr
	// XXX dont know magic numbers just guessing
	if b[0] != 'I' || b[1] != 'I' || b[2] != '*' || b[3] != 0 {
		return false
	}
	return true
}

func parseTIFF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Image
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4,
		Info:   "file header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
