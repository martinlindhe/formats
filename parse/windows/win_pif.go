package windows

// Windows Program Information File (PIF)

// STATUS: 1%

import (
	"github.com/martinlindhe/formats/parse"
)

var (
	pifHeaderPos = int64(0x171)
)

func PIF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isPIF(c.Header) {
		return nil, nil
	}
	return parsePIF(c)
}

func isPIF(b []byte) bool {

	s := string(b[pifHeaderPos : pifHeaderPos+15])
	return s == "MICROSOFT PIFEX"
}

func parsePIF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

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
