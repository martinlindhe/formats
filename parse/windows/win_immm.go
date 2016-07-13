package windows

// STATUS: 1%
// Extensions: .db
// found on Windows 10 Users/vm/AppData/Local/Microsoft/Windows/Explorer/iconcache_idx.db

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// IMMM parses the immm format
func IMMM(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isIMMM(c.Header) {
		return nil, nil
	}
	return parseIMMM(c.File, c.ParsedLayout)
}

func isIMMM(b []byte) bool {

	// XXX just guessing, saw these 8 bytes on all samples
	if b[0] != 0 || b[1] != 0x30 || b[2] != 0x20 || b[3] != 0x10 ||
		b[4] != 'I' || b[5] != 'M' || b[6] != 'M' || b[7] != 'M' {
		return false
	}
	return true
}

func parseIMMM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
