package windows

// STATUS: 1%
// Extensions: .etl
// found on Windows 10, Windows/System32/LogFiles/WMI/RtBackup/EtwRTDiagLog.etl

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// RLFS parses the rlfs format
func RLFS(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isRLFS(c.Header) {
		return nil, nil
	}
	return parseRLFS(c.File, c.ParsedLayout)
}

func isRLFS(b []byte) bool {

	if b[0] != 'R' || b[1] != 'l' || b[2] != 'F' || b[3] != 's' {
		return false
	}
	return true
}

func parseRLFS(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
