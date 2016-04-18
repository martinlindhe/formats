package bin

// PDB, Visual Studio debug info
// STATUS: 1%

import (
	"github.com/martinlindhe/formats/parse"
	"os"
)

func PDB(file *os.File) (*parse.ParsedLayout, error) {

	if !isPDB(file) {
		return nil, nil
	}
	return parsePDB(file)
}

func isPDB(file *os.File) bool {

	s, _, _ := parse.ReadZeroTerminatedASCIIUntil(file, 0, 26)

	// XXX just guessing
	if s != "Microsoft C/C++ MSF 7.00"+"\r\n" {
		return false
	}

	return true
}

func parsePDB(file *os.File) (*parse.ParsedLayout, error) {

	pos := int64(0)
	res := parse.ParsedLayout{
		FileKind: parse.Binary,
		Layout: []parse.Layout{{
			Offset: 0,
			Length: 26, // XXX
			Info:   "header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 26, Info: "magic", Type: parse.ASCII},
			}}}}

	return &res, nil
}
