package parse

// for unrecognized files

import (
	"encoding/binary"
	"fmt"
)

func RAW(c *ParseChecker) (*ParsedLayout, error) {

	format := "raw"
	if c.ParsedLayout.FileSize == 0 {
		format = "empty"
	}

	pos := int64(0)
	c.ParsedLayout.FileKind = Binary
	c.ParsedLayout.MimeType = "application/octet-stream"
	c.ParsedLayout.Layout = []Layout{{
		Offset: pos,
		Length: 0,
		Info:   "unrecognized data",
		Type:   Group,
		Childs: []Layout{
			{Offset: pos, Length: 0, Info: "data", Type: Bytes},
		}}}

	val := binary.LittleEndian.Uint32(c.Header)
	sig := string(c.Header[0:4])
	c.ParsedLayout.FormatName = fmt.Sprintf(" [%08x, %s]", val, sig) + " " + format

	return &c.ParsedLayout, nil
}
