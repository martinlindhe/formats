package windows

// ???
// found on Windows 10 Program Files/WindowsApps/Microsoft.BingNews_4.3.193.0_x86__8wekyb3d8bbwe/_Resources/0.rsrc
// extensions: .rsrc

// STATUS: 1%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func RSRC(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isRSRC(c.File) {
		return nil, nil
	}
	return parseRSRC(c.File, c.ParsedLayout)
}

func isRSRC(file *os.File) bool {

	val, _ := parse.ReadUint32le(file, 0)
	fmt.Printf("rsrc xxx = %04x\n", val)
	if val == 0xbeefcace {
		return true
	}
	return false
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
