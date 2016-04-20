package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func VDI(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isVDI(c.File) {
		return nil, nil
	}
	return parseVDI(c.File, c.ParsedLayout)
}

func isVDI(file *os.File) bool {

	s, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, 0, 40)
	if s != "<<< Oracle VM VirtualBox Disk Image >>>"+"\n" {
		return false
	}
	return true
}

func parseVDI(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 40,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 40, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
