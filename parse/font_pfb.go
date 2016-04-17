package parse

// Adobe Printer Font Binary (used in the '90s)
// STATUS 1%

import (
	"os"
)

func PFB(file *os.File) (*ParsedLayout, error) {

	if !isPFB(file) {
		return nil, nil
	}
	return parsePFB(file)
}

func isPFB(file *os.File) bool {

	// XXX just guessing ...
	s, _, _ := readZeroTerminatedASCIIUntil(file, 6, 10)
	if s != "%!FontType" {
		return false
	}
	return true
}

func parsePFB(file *os.File) (*ParsedLayout, error) {

	offset := int64(0)

	res := ParsedLayout{
		FileKind: Font,
		Layout: []Layout{{
			Offset: offset,
			Length: 16, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: offset, Length: 6, Info: "unknown", Type: Bytes},
				{Offset: offset + 6, Length: 10, Info: "magic", Type: ASCII},
			}}}}

	return &res, nil
}
