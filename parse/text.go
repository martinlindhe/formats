package parse

// STATUS: 50%

// XXX detect bom
// XXX try to guess encoding (utf8 / ascii..)

import (
	"fmt"
	"io"
	"strings"
)

func Text(c *ParseChecker) (*ParsedLayout, error) {

	if !isText(c) {
		return nil, fmt.Errorf("no match")
	}

	return parseText(c)
}

func isText(c *ParseChecker) bool {

	b := c.Header

	for pos := int64(0); pos < 10; pos++ {
		if pos >= c.ParsedLayout.FileSize {
			break
		}

		// US-ASCII check
		c := b[pos]
		if c < 32 || c > 126 {
			if c != '\n' && c != '\r' && c != '\t' {
				return false
			}
		}
	}
	return true
}

func parseText(c *ParseChecker) (*ParsedLayout, error) {

	c.ParsedLayout.FormatName = "text"

	pos := int64(0)
	hdr, _, _ := ReadZeroTerminatedASCIIUntil(c.File, pos, 5)
	if strings.ToLower(hdr) == "<?xml" {
		c.ParsedLayout.FormatName = "xml"
	}

	layout := Layout{
		Offset: pos,
		Info:   "text",
		Type:   Group}

	line := 1
	for {

		_, len, err := ReadBytesUntilNewline(c.File, pos)
		if err != nil {
			if err != io.EOF {
				fmt.Println("err!", err)
			}
			break
		}

		layout.Childs = append(layout.Childs, Layout{
			Offset: pos,
			Length: len,
			Info:   "line " + fmt.Sprintf("%d", line),
			Type:   ASCII})

		layout.Length += len
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
