package parse

// PDB, Visual Studio debug info

import (
	"os"
)

func PDB(file *os.File) (*ParsedLayout, error) {

	if !isPDB(file) {
		return nil, nil
	}
	return parsePDB(file)
}

func isPDB(file *os.File) bool {

	s, _, _ := readZeroTerminatedASCIIUntil(file, 0, 26)

	// XXX just guessing
	if s != "Microsoft C/C++ MSF 7.00"+"\r\n" {
		return false
	}

	return true
}

func parsePDB(file *os.File) (*ParsedLayout, error) {

	pos := int64(0)

	res := ParsedLayout{
		FileKind: Binary,
		Layout: []Layout{{
			Offset: 0,
			Length: 26, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 26, Info: "magic", Type: ASCII},
			}}}}

	return &res, nil
}
