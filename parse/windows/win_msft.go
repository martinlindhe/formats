package windows

// STATUS: 1%
// Extensions: .tlb
// found on Windows 10, Windows/WinSxS/amd64_netfx4-mscorlib_tlb_b03f5f7f11d50a3a_4.0.10240.16384_none_cb57103f03cae093/mscorlib.tlb

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// MSFT parses the msft format
func MSFT(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isMSFT(c.Header) {
		return nil, nil
	}
	return parseMSFT(c.File, c.ParsedLayout)
}

func isMSFT(b []byte) bool {

	if b[0] != 'M' || b[1] != 'S' || b[2] != 'F' || b[3] != 'T' {
		return false
	}
	return true
}

func parseMSFT(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Uint32le},
		}}}

	return &pl, nil
}
