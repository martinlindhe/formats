package parse

// used to show unrecognized files

import (
	"os"
)

func RAW(file *os.File) (*ParsedLayout, error) {

	// TODO: make cmd/formats work without any Layout

	return &ParsedLayout{
		FileKind: Binary,
		Layout: []Layout{{
			Offset: 0,
			Length: 0,
			Info:   "unrecognized data",
			Type:   Group,
			Childs: []Layout{
				{Offset: 0, Length: 0, Info: "data", Type: Bytes},
			}}}}, nil
}
