package image

// Windows Icon / Cursor image resources
// TODO icon_embedded_png_001.ico has embedded PNG in image data
// TODO decode non-png as "standard BMP image"...
// TODO offer "save to file" for the "image data" (bytes type)

// STATUS: 80%

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/martinlindhe/formats/parse"
)

func ICO(file *os.File, hdr [0xffff]byte, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	if !isICO(file) {
		return nil, nil
	}
	return parseICO(file, pl)
}

func isICO(file *os.File) bool {

	b, _ := readIconHeader(file)
	if b[0] != 0 {
		return false
	}

	// 1 = icon, 2 = cursor
	if b[1] != 1 && b[1] != 2 {
		return false
	}

	return true
}

func parseICO(file *os.File, pl parse.ParsedLayout) (*parse.ParsedLayout, error) {

	pl.FileKind = parse.Image
	pos := int64(0)
	typeName := ""

	hdr, _ := readIconHeader(file)
	switch hdr[1] {
	case 1:
		typeName = "icon"
	case 2:
		typeName = "cursor"
	default:
		typeName = "unknown"
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

	pl.Layout = append(pl.Layout, fileHeader)

	for i := 0; i < int(numIcons); i++ {
		id := fmt.Sprintf("%d", i+1)

		pl.Layout = append(pl.Layout, parse.Layout{
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

		dataOffset, err := pl.ReadUint32leFromInfo(file, "offset to resource "+id)
		if err != nil {
			return nil, err
		}
		dataSize, err := pl.ReadUint32leFromInfo(file, "data size of resource "+id)
		if err != nil {
			return nil, err
		}

		pl.Layout = append(pl.Layout, parse.Layout{
			Offset: int64(dataOffset),
			Length: int64(dataSize),
			Info:   "resource " + id + " data",
			Type:   parse.Group,
			Childs: []parse.Layout{
				{Offset: int64(dataOffset), Length: int64(dataSize), Info: "image data", Type: parse.Bytes},
			}})
	}

	return &pl, nil
}

func readIconHeader(file *os.File) ([3]uint16, error) {

	file.Seek(0, os.SEEK_SET)
	var b [3]uint16
	if err := binary.Read(file, binary.LittleEndian, &b); err != nil {
		return b, err
	}
	return b, nil
}
