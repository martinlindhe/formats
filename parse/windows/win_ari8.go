package windows

// ???
// found on Windows 10 in ProgramData/Microsoft/Windows/AppRepository/Packages/Microsoft.BingFinance_10004.3.193.0_neutral_~_8wekyb3d8bbwe/S-1-5-18.recovery
// extensions: .pckgdep, .recovery

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func ARI8(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isARI8(&c.Header) {
		return nil, nil
	}
	return parseARI8(c.File, c.ParsedLayout)
}

func isARI8(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'A' || b[1] != 'R' || b[2] != 'I' || b[3] != '8' {
		return false
	}
	return true
}

func parseARI8(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
