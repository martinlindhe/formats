package bin

// Program Database (Visual Studio debug info)
// https://en.wikipedia.org/wiki/Program_database

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func PDB(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isPDB(c.File) {
		return nil, nil
	}
	return parsePDB(c.File, c.ParsedLayout)
}

func isPDB(file *os.File) bool {

	s, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, 0, 26)

	// XXX just guessing
	if s != "Microsoft C/C++ MSF 7.00"+"\r\n" {
		return false
	}

	return true
}

func parsePDB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: 0,
		Length: 26, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 26, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
