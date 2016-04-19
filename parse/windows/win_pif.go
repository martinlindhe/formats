package windows

// Windows Program Information File (PIF)

// STATUS: 1%

import (
	"github.com/martinlindhe/formats/parse"
	"os"
)

func PIF(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isPIF(file) {
		return nil, nil
	}
	return parsePIF(file, pl)
}

func isPIF(file *os.File) bool {

	s, _, err := parse.ReadZeroTerminatedASCIIUntil(file, 0x171, 15)
	if err != nil {
		return false
	}
	if s == "MICROSOFT PIFEX" {
		return true
	}
	return false
}

func parsePIF(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0x171)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 15, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 15, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
