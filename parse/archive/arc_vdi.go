package archive

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// VDI parses the vdi format
func VDI(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isVDI(c.Header) {
		return nil, nil
	}
	return parseVDI(c.File, c.ParsedLayout)
}

func isVDI(b []byte) bool {

	s := string(b[0:40])
	return s == "<<< Oracle VM VirtualBox Disk Image >>>"+"\n"
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
