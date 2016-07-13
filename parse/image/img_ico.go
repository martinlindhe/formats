package image

// TODO icon_embedded_png_001.ico has embedded PNG in image data
// TODO decode non-png as "standard BMP image"...
// TODO offer "save to file" for the "image data" (bytes type)

// STATUS: 80%

import (
	"encoding/binary"
	"fmt"

	"github.com/martinlindhe/formats/parse"
)

var (
	iconTypes = map[uint16]string{
		1: "icon",
		2: "cursor",
	}
)

// ICO parses the Windows Icon / Cursor image resource format
func ICO(c *parse.Checker) (*parse.ParsedLayout, error) {

	if !isICO(c.Header) {
		return nil, nil
	}
	return parseICO(c)
}

func isICO(b []byte) bool {

	h := readIconHeader(b)
	if h[0] != 0 {
		return false
	}
	if h[1] != 1 && h[1] != 2 {
		// 1 = icon, 2 = cursor
		return false
	}
	if h[2] > 500 {
		// NOTE: an arbitrary check to get less false matches
		return false
	}
	return true
}

func parseICO(c *parse.Checker) (*parse.ParsedLayout, error) {

	c.ParsedLayout.FileKind = parse.Image
	c.ParsedLayout.MimeType = "image/x-ico"
	pos := int64(0)

	hdr := readIconHeader(c.Header)
	typeName := "?"
	if val, ok := iconTypes[hdr[1]]; ok {
		typeName = val
	}
	fileHeader := parse.Layout{
		Offset: pos,
		Length: 6,
		Info:   "header",
		Type:   parse.Group,
		Childs: []parse.Layout{
			{Offset: 0, Length: 2, Info: "magic", Type: parse.Uint16le},
			{Offset: 2, Length: 2, Info: "type = " + typeName, Type: parse.Uint16le},
			{Offset: 4, Length: 2, Info: "number of resources", Type: parse.Uint16le},
		}}

	pos += fileHeader.Length

	numIcons := hdr[2]
	resourceEntryLength := int64(16)

	c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, fileHeader)

	for i := 0; i < int(numIcons); i++ {
		id := fmt.Sprintf("%d", i+1)

		c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, parse.Layout{
			Offset: pos,
			Length: resourceEntryLength,
			Info:   "resource " + id + " header",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: pos, Length: 1, Info: "width", Type: parse.Uint8},
				{Offset: pos + 1, Length: 1, Info: "height", Type: parse.Uint8},
				{Offset: pos + 2, Length: 1, Info: "max number of colors", Type: parse.Uint8},
				{Offset: pos + 3, Length: 1, Info: "reserved", Type: parse.Uint8},
				{Offset: pos + 4, Length: 2, Info: "planes", Type: parse.Uint16le},
				{Offset: pos + 6, Length: 2, Info: "bit count", Type: parse.Uint16le},
				{Offset: pos + 8, Length: 4, Info: "data size of resource " + id, Type: parse.Uint32le},
				{Offset: pos + 12, Length: 4, Info: "offset to resource " + id, Type: parse.Uint32le},
			}})
		fileHeader.Length += resourceEntryLength
		pos += resourceEntryLength
	}

	for i := 0; i < int(numIcons); i++ {
		id := fmt.Sprintf("%d", i+1)

		dataOffset, err := c.ParsedLayout.ReadUint32leFromInfo(c.File, "offset to resource "+id)
		if err != nil {
			return nil, err
		}
		dataSize, err := c.ParsedLayout.ReadUint32leFromInfo(c.File, "data size of resource "+id)
		if err != nil {
			return nil, err
		}

		c.ParsedLayout.Layout = append(c.ParsedLayout.Layout, parse.Layout{
			Offset: int64(dataOffset),
			Length: int64(dataSize),
			Info:   "resource " + id + " data",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: int64(dataOffset), Length: int64(dataSize), Info: "image data", Type: parse.Bytes},
			}})
	}

	return &c.ParsedLayout, nil
}

func readIconHeader(b []byte) [3]uint16 {

	var h [3]uint16
	h[0] = binary.LittleEndian.Uint16(b)
	h[1] = binary.LittleEndian.Uint16(b[2:])
	h[2] = binary.LittleEndian.Uint16(b[4:])
	return h
}
