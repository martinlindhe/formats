package windows

// STATUS: 1%
// Extensions: .rsrc
// found on Windows 10 Program Files/WindowsApps/Microsoft.BingNews_4.3.193.0_x86__8wekyb3d8bbwe/_Resources/0.rsrc

import (
	"encoding/binary"
	"os"

	"github.com/martinlindhe/formats/parse"
)

// RSRC parses the rsrc format
func RSRC(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isRSRC(c.Header) {
		return nil, nil
	}
	return parseRSRC(c.File, c.ParsedLayout)
}

func isRSRC(b []byte) bool {

	val := binary.LittleEndian.Uint32(b)
	return val == 0xbeefcace
}

func parseRSRC(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

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
