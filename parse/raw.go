package parse

// for unrecognized files

import (
	"fmt"
)

func RAW(c *ParseChecker) (*ParsedLayout, error) {

	format := "raw"
	if c.ParsedLayout.FileSize == 0 {
		format = "empty"
	}

	// TODO: make cmd/formats work without any Layout, to avoid a 0-length selected area
	pos := int64(0)
	c.ParsedLayout.FormatName = format
	c.ParsedLayout.FileKind = Binary
	c.ParsedLayout.Layout = []Layout{{
		Offset: pos,
		Length: 0,
		Info:   "unrecognized data",
		Type:   Group,
		Childs: []Layout{
			{Offset: pos, Length: 0, Info: "data", Type: Bytes},
		}}}

	val, _ := ReadUint32le(c.File, 0)
	sig, _, _ := ReadZeroTerminatedASCIIUntil(c.File, 0, 4)
	c.ParsedLayout.FormatName += fmt.Sprintf(" [%08x, %s]", val, sig)

	return &c.ParsedLayout, nil
}
