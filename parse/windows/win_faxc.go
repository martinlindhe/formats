package windows

// STATUS: 1%
// Extensions: .cov
// ??? found on Windows 10, ProgramData/Microsoft/Windows NT/MSFax/Common Coverpages/en-GB/confident.cov

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// FAXC parses the faxc format
func FAXC(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isFAXC(c.Header) {
		return nil, nil
	}
	return parseFAXC(c.File, c.ParsedLayout)
}

func isFAXC(b []byte) bool {

	s := string(b[0:12])
	return s == "FAXCOVER-VER"
}

func parseFAXC(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 12, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 12, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
