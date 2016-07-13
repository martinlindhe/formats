package windows

// ???
// NOTE: not to be confused with .chm (see doc_chm.go)
// found on Windows 10 /Users/m/Downloads/old-windows/win10-os/Users/vm/AppData/Local/Microsoft/Windows/Explorer/thumbcache_sr.db
// extensions: .db

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// CHMM parses the chmm format
func CHMM(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isCHMM(c.Header) {
		return nil, nil
	}
	return parseCHMM(c.File, c.ParsedLayout)
}

func isCHMM(b []byte) bool {

	if b[0] != 'C' || b[1] != 'H' || b[2] != 'M' || b[3] != 'M' {
		return false
	}
	return true
}

func parseCHMM(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
