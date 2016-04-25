package av

// ???, from OS/2 Warp 4 setup cd OS2IMAGE/FI/JAVAOS2/DEMO/
// Extensions: .au

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func AU(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isAU(c.Header) {
		return nil, nil
	}
	return parseAU(c.File, c.ParsedLayout)
}

func isAU(b []byte) bool {

	if b[0] != '.' || b[1] != 's' || b[2] != 'n' || b[3] != 'd' {
		return false
	}
	return true
}

func parseAU(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.AudioVideo
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
