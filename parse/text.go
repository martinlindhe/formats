package parse

// STATUS: 60%

// https://en.wikipedia.org/wiki/Byte_order_mark

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

	if hasRecognizedBOM(c) {
		return true
	}

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

	layout := Layout{
		Offset: pos,
		Info:   "text",
		Type:   Group}

	bom, bomLen := parseBOMMark(c, pos)
	if bomLen > 0 {
		c.ParsedLayout.TextEncoding = bom
		layout.Childs = append(layout.Childs, Layout{
			Offset: pos,
			Length: bomLen,
			Info:   bom.String() + " bom",
			Type:   Bytes})

		layout.Length += bomLen
		pos += bomLen
	}

	data := ReadBytesFrom(c.File, pos, 5)
	if strings.ToLower(string(data)) == "<?xml" {
		c.ParsedLayout.FormatName = "xml"
	}

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

func parseBOMMark(c *ParseChecker, pos int64) (TextEncoding, int64) {

	b := c.Header
	if b[0] == 0xff && b[1] == 0xfe && b[2] == 0 && b[3] == 0 {
		return UTF32le, 2
	}
	if b[0] == 0 && b[1] == 0 && b[2] == 0xfe && b[3] == 0xff {
		return UTF32be, 4
	}
	if b[0] == 0xfe && b[1] == 0xff {
		return UTF16be, 2
	}
	if b[0] == 0xff && b[1] == 0xfe {
		return UTF16le, 2
	}
	if b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		return UTF8, 3
	}
	return None, 0
}

func hasRecognizedBOM(c *ParseChecker) bool {

	_, len := parseBOMMark(c, 0)
	if len > 0 {
		return true
	}
	return false
}
