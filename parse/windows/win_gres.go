package windows

// STATUS: 1%
// Extensions: .db
// found on Windows 10 Windows/WinSxS/amd64_microsoft-windows-mapcontrol_31bf3856ad364e35_10.0.10240.16384_none_1b558da4a5404873/resource.db

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// GRES parses the gres format
func GRES(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isGRES(c.Header) {
		return nil, nil
	}
	return parseGRES(c.File, c.ParsedLayout)
}

func isGRES(b []byte) bool {

	if b[0] != 'G' || b[1] != 'R' || b[2] != 'E' || b[3] != 'S' {
		return false
	}
	return true
}

func parseGRES(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
