package windows

// ??? found on Windows 10, Windows/WinSxS/FileMaps/program_files_x86_common_files_microsoft_shared_dao_9d0cb78256d5a29d.cdf-ms

// Extensions: .cdf-ms

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func PCMH(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isPCMH(c.Header) {
		return nil, nil
	}
	return parsePCMH(c.File, c.ParsedLayout)
}

func isPCMH(b []byte) bool {

	if b[0] != 'P' || b[1] != 'c' || b[2] != 'm' || b[3] != 'H' {
		return false
	}
	return true
}

func parsePCMH(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
