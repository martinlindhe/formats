package bin

// ??? found in MacOS /System/Library/Perl/Extras/5.18/darwin-thread-multi-2level/XML/Parser/Encodings/euc-kr.enc

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

// PerlENC parses the perl encodings format
func PerlENC(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isPerlENC(c.Header) {
		return nil, nil
	}
	return parsePerlENC(c.File, c.ParsedLayout)
}

func isPerlENC(b []byte) bool {

	if b[0] != 0xfe || b[1] != 0xeb || b[2] != 0xfa || b[3] != 0xce {
		return false
	}
	return true
}

func parsePerlENC(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
	pl.Layout = []parse.Layout{{
		Offset: pos,
		Length: 4, // XXX
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: pos, Length: 4, Info: "magic", Type: parse.ASCII},
		}}}

	return &pl, nil
}
