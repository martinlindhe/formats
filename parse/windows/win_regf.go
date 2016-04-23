package windows

// ???
// found on Windows 10 Users/vm/AppData/Local/Packages/Microsoft.AAD.BrokerPlugin_cw5n1h2txyewy/Settings/settings.dat
// extensions: .dat

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func REGF(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isREGF(c.Header) {
		return nil, nil
	}
	return parseREGF(c.File, c.ParsedLayout)
}

func isREGF(b []byte) bool {

	if b[0] != 'r' || b[1] != 'e' || b[2] != 'g' || b[3] != 'f' {
		return false
	}
	return true
}

func parseREGF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
