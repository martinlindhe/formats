package windows

// ??? found on Windows 10, Windows/System32/winevt/Logs/System.evtx
// Extensions: .evtx

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func EVTX(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isEVTX(c.Header) {
		return nil, nil
	}
	return parseEVTX(c.File, c.ParsedLayout)
}

func isEVTX(b []byte) bool {

	s := string(b[0:7])
	return s == "ElfFile" // XXX why elf file???
}

func parseEVTX(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 7, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 7, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
