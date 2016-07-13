package bin

// https://en.wikipedia.org/wiki/Program_database

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// PDB parses the Program Database format (Visual Studio debug info)
func PDB(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isPDB(c.Header) {
		return nil, nil
	}
	return parsePDB(c.File, c.ParsedLayout)
}

func isPDB(b []byte) bool {

	s := string(b[0:26])
	if s != "Microsoft C/C++ MSF 7.00"+"\r\n" {
		return false
	}
	return true
}

func parsePDB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 26, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 26, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
