package archive

// Debian package

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func DEB(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isDEB(c) {
		return nil, nil
	}
	return parseDEB(c.File, c.ParsedLayout)
}

func isDEB(c *parse.ParseChecker) bool {

	s, _, _ := parse.ReadZeroTerminatedASCIIUntil(c.File, 0, 21)
	if s != "!<arch>\n"+"debian-binary" {
		return false
	}
	return true
}

func parseDEB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 21, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 21, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
