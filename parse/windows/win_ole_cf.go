package windows

// Windows system format, used by:
//   MS Word documents (.doc, .pps, .ppt, .xls)
//   Thumbs.DB

// http://www.forensicswiki.org/wiki/Thumbs.db
// http://www.forensicswiki.org/wiki/OLE_Compound_File

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// OLECF parses the OLE Compound File format
func OLECF(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isOLECF(c.Header) {
		return nil, nil
	}
	return parseOLECF(c.File, c.ParsedLayout)
}

func isOLECF(b []byte) bool {

	if b[0] != 0xd0 || b[1] != 0xcf || b[2] != 0x11 || b[3] != 0xe0 ||
		b[4] != 0xa1 || b[5] != 0xb1 || b[6] != 0x1a || b[7] != 0xe1 {
		return false
	}
	return true
}

func parseOLECF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 8, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 8, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
