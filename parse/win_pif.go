package parse

// Windows Program Information File (PIF)
// STATUS: 1%

import (
	"os"
)

func PIF(file *os.File) (*ParsedLayout, error) {

	if !isPIF(file) {
		return nil, nil
	}
	return parsePIF(file)
}

func isPIF(file *os.File) bool {

	s, err := knownLengthASCII(file, 0x171, 15)
	if err != nil {
		return false
	}
	if s == "MICROSOFT PIFEX" {
		return true
	}
	return false
}

func parsePIF(file *os.File) (*ParsedLayout, error) {

	pos := int64(0x171)
	res := ParsedLayout{
		FileKind: WindowsResource,
		Layout: []Layout{{
			Offset: pos,
			Length: 15, // XXX
			Info:   "header",
			Type:   Group,
			Childs: []Layout{
				{Offset: pos, Length: 15, Info: "magic", Type: Uint32le},
			}}}}

	return &res, nil
}
