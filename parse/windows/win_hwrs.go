package windows

// STATUS: 1%
// found on Windows 10 in Program Files/Common Files/microsoft shared/ink/{hwrcommonlm,hwrlatinlm}.dat
// XXX also this string at 0x10: "Microsoft Common LM Resource File"

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// HWRS parses the hwrs format
func HWRS(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isHWRS(c.Header) {
		return nil, nil
	}
	return parseHWRS(c.File, c.ParsedLayout)
}

func isHWRS(b []byte) bool {

	if b[0] != 'H' || b[1] != 'W' || b[2] != 'R' || b[3] != 'S' {
		return false
	}
	return true
}

func parseHWRS(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
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
