package archive

// ???, some compression used on OS/2 Warp 4 setup cd

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func FTCOMP(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isFTCOMP(c.Header) {
		return nil, nil
	}
	return parseFTCOMP(c.File, c.ParsedLayout)
}

func isFTCOMP(b []byte) bool {

	if b[0] != 0xa5 || b[1] != 0x96 || b[2] != 0xfd || b[3] != 0xff {
		return false
	}
	return true
}

func parseFTCOMP(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Archive
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
