package parse

import (
	"fmt"
	"strings"
)

func Text(c *ParseChecker) (*ParsedLayout, error) {

	if !isText(&c.Header) {
		return nil, fmt.Errorf("no match")
	}

	return parseText(c)
}

func isText(hdr *[0xffff]byte) bool {

	b := *hdr

	for pos := 0; pos < 10; pos++ {
		// US-ASCII check
		c := b[pos]
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}

func parseText(c *ParseChecker) (*ParsedLayout, error) {

	// XXX detect bom

	// XXX detect line endings
	// XXX try to guess encoding (utf8 / ascii..)

	// XXX parse line by line as blocks

	c.ParsedLayout.FormatName = "text"

	pos := int64(0)
	hdr, _, _ := ReadZeroTerminatedASCIIUntil(c.File, pos, 5)
	if strings.ToLower(hdr) == "<?xml" {
		c.ParsedLayout.FormatName = "xml"
	}

	layout := Layout{
		Offset: pos,
		Length: 0,
		Info:   "text",
		Type:   Group}

	line := 1
	for {

		_, len, err := ReadBytesUntilNewline(c.File, pos)
		if err != nil {
			fmt.Println("err!", err)
			break
		}

		layout.Childs = append(layout.Childs, Layout{
			Offset: pos,
			Length: len,
			Info:   "line " + fmt.Sprintf("%d", line),
			Type:   Bytes})

		pos += len
		line++
		if pos >= c.ParsedLayout.FileSize {
			break
		}
	}

	c.ParsedLayout.FileKind = Document
	c.ParsedLayout.Layout = []Layout{layout}

	return &c.ParsedLayout, nil
}
