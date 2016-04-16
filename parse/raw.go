package parse

// used to show unrecognized files

import (
	"os"
)

func RAW(file *os.File) (*ParsedLayout, error) {

	len := fileSize(file)

	return &ParsedLayout{
		FileKind: Binary,
		Layout: []Layout{{
			Offset: 0,
			Length: len,
			Info:   "unrecognized data",
			Type:   Group,
			Childs: []Layout{
				{Offset: 0, Length: len, Info: "data", Type: Bytes},
			}}}}, nil
}
