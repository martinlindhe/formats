package bin

// Maple help database
// .hdb

// STATUS: 1%

import (
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func MapleDB(c *parse.ParseChecker) (*parse.ParsedLayout, error) {

	if !isMapleDB(c.File) {
		return nil, nil
	}
	return parseMapleDB(c.File, c.ParsedLayout)
}

func isMapleDB(file *os.File) bool {

	val, _ := parse.ReadUint32le(file, 0)
	fmt.Printf("val %08x", val)
	if val == 0x00040000 {
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
			{Offset: pos, Length: 4, Info: "magic", Type: parse.Bytes},
		}}}

	return &pl, nil
}
