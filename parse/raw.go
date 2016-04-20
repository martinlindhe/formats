package parse

// used to show unrecognized files

import (
	"os"
)

func RAW(file *os.File) (*ParsedLayout, error) {

	// TODO: make cmd/formats work without any Layout, to avoid a 0-length selected area

	return &ParsedLayout{
		FormatName: "raw",
		FileKind:   Binary,
		Layout: []Layout{{
			Offset: 0,
			Length: 0,
			Info:   "unrecognized data",
			Type:   Group,
			Childs: []Layout{
				{Offset: 0, Length: 0, Info: "data", Type: Bytes},
			}}}}, nil
}
