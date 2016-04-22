package image

// https://en.wikipedia.org/wiki/WebP
// Extensions: .webp

// STATUS: 1%

// XXX uses RIFF container!

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func WebP(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isWebP(c.Header) {
		return nil, nil
	}
	return parseWebP(c.File, c.ParsedLayout)
}

func isWebP(b []byte) bool {

	if b[0] != 'R' || b[1] != 'I' || b[2] != 'F' || b[3] != 'F' {
		return false
	}
	if b[8] != 'W' || b[9] != 'E' || b[10] != 'B' || b[11] != 'P' ||
		b[12] != 'V' || b[13] != 'P' || b[14] != '8' || b[15] != ' ' {
		return false
	}
	return true
}

func parseWebP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Image
	pl.MimeType = "image/webp"
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 16,
		Info:   "file header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
			{Offset: pos + 4, Length: 4, Info: "length", Type: parse.Uint32le},
			{Offset: pos + 8, Length: 8, Info: "magic", Type: parse.ASCII},
			// XXX the rest
		}}}

	return &pl, nil
}
