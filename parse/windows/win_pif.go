package windows

// STATUS: 1%
// Extensions: .pif

import (
	"github.com/martinlindhe/formats/parse"
)

var (
	pifHeaderPos = int64(0x171)
)

// PIF parses the Windows Program Information File format
func PIF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isPIF(c.Header) {
		return nil, nil
	}
	return parsePIF(c)
}

func isPIF(b []byte) bool {

	s := string(b[pifHeaderPos : pifHeaderPos+15])
	return s == "MICROSOFT PIFEX"
}

func parsePIF(c *parse.Checker) (*parse.ParsedLayout, error) {

	pos := int64(pifHeaderPos)
	c.ParsedLayout.FileKind = parse.WindowsResource
	c.ParsedLayout.Layout = []parse.Layout{{
		Offset: pos,
		Length: 15, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 15, Info: "magic", Type: parse.Uint32le},
		}}}

	return &c.ParsedLayout, nil
}
