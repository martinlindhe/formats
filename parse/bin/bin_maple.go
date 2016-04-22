package bin

// Maple help database
// .hdb

// STATUS: 1%

import (
	"encoding/binary"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func MapleDB(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isMapleDB(c.Header) {
		return nil, nil
	}
	return parseMapleDB(c.File, c.ParsedLayout)
}

func isMapleDB(b []byte) bool {

	val := binary.LittleEndian.Uint32(b)
	if val == 0x00000400 {
		return true
	}
	return false
}

func parseMapleDB(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pos := int64(0)
	pl.FileKind = parse.Binary
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
