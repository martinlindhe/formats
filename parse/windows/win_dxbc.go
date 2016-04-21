package windows

// Direct3D shader bytecode

// http://timjones.tw/blog/archive/2015/09/02/parsing-direct3d-shader-bytecode

// STATUS: 1%

import (
	"os"

	"github.com/martinlindhe/formats/parse"
)

func DXBC(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isDXBC(&c.Header) {
		return nil, nil
	}
	return parseDXBC(c.File, c.ParsedLayout)
}

func isDXBC(hdr *[0xffff]byte) bool {

	b := *hdr
	if b[0] != 'D' || b[1] != 'X' || b[2] != 'B' || b[3] != 'C' {
		return false
	}
	return true
}

func parseDXBC(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.WindowsResource
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
