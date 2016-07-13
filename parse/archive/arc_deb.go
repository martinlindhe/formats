package archive

// Debian package

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// DEB parses the deb format
func DEB(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isDEB(c.Header) {
		return nil, nil
	}
	return parseDEB(c.File, c.ParsedLayout)
}

func isDEB(b []byte) bool {

	s := string(b[0:21])
	return s == "!<arch>\n"+"debian-binary"
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
