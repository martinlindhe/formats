package bin

// ??? extension .sdg, magic "SGA3", found in
// /usr/lib/libreoffice/share/gallery/environment.sdg on ubuntu

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func SGA(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isSGA(&c.Header) {
		return nil, nil
	}
	return parseSGA(c.File, c.ParsedLayout)
}

func isSGA(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] == 'S' && b[1] == 'G' && b[2] == 'A' && b[3] == '3' {
		return true
	}
	return false
}

func parseSGA(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
