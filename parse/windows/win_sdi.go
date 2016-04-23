package windows

// ???
// found on Windows 10 Windows/SysWOW64/boot.sdi
// extensions: .sdi

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func SDI(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isSDI(c.Header) {
		return nil, nil
	}
	return parseSDI(c.File, c.ParsedLayout)
}

func isSDI(b []byte) bool {

	// XXX just guessing
	if b[0] != '$' || b[1] != 'S' || b[2] != 'D' || b[3] != 'I' ||
		b[4] != '0' || b[5] != '0' || b[6] != '0' || b[7] != '1' {
		return false
	}
	return true
}

func parseSDI(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 8, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 8, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
